package dataframe

import (
	"fmt"
	"sort"

	"go.starlark.net/starlark"
)

// SeriesGroupByResult is the result of indexing into a groupBy result
type SeriesGroupByResult struct {
	lhsLabel string
	rhsLabel string
	grouping map[string]*Series
	// index of the source DataFrame
	dfIndex *Index
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*SeriesGroupByResult)(nil)
	_ starlark.HasAttrs = (*SeriesGroupByResult)(nil)
)

const (
	indexTypeNone         = 0
	indexTypeFromOriginal = 1
	indexTypeBuildNew     = 2
)

var seriesGroupByResultMethods = map[string]*starlark.Builtin{
	"count": starlark.NewBuiltin("count", seriesGroupByResultCount),
	"sum":   starlark.NewBuiltin("sum", seriesGroupByResultSum),
	"apply": starlark.NewBuiltin("apply", seriesGroupByResultApply),
}

// Freeze has no effect on the immutable SeriesGroupByResult
func (sgbr *SeriesGroupByResult) Freeze() {
	// pass
}

// Hash cannot be used with SeriesGroupByResult
func (sgbr *SeriesGroupByResult) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", sgbr.Type())
}

// String returns a string representation of the SeriesGroupByResult
func (sgbr *SeriesGroupByResult) String() string {
	return fmt.Sprintf("<%s>", sgbr.Type())
}

// Truth converts the SeriesGroupByResult into a bool
func (sgbr *SeriesGroupByResult) Truth() starlark.Bool {
	return true
}

// Type returns the type as a string
func (sgbr *SeriesGroupByResult) Type() string {
	return fmt.Sprintf("%s.SeriesGroupByResult", Name)
}

// Attr gets a value for an attribute
func (sgbr *SeriesGroupByResult) Attr(name string) (starlark.Value, error) {
	return builtinAttr(sgbr, name, seriesGroupByResultMethods)
}

// AttrNames lists available attributes
func (sgbr *SeriesGroupByResult) AttrNames() []string {
	return builtinAttrNames(seriesGroupByResultMethods)
}

// sum method returns a Series that is the sum of each grouped result
func seriesGroupByResultSum(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("sum", args, kwargs); err != nil {
		return nil, err
	}
	self := b.Receiver().(*SeriesGroupByResult)

	indexTexts := []string{}
	vals := []int{}

	sortedKeys := getSortedKeys(self.grouping)
	for _, groupName := range sortedKeys {
		series := self.grouping[groupName]

		sum := 0
		for i := 0; i < series.Len(); i++ {
			elem := series.Index(i)
			if num, err := starlark.AsInt32(elem); err == nil {
				sum += num
			}
		}

		indexTexts = append(indexTexts, groupName)
		vals = append(vals, sum)
	}

	index := NewIndex(indexTexts, self.lhsLabel)
	return newSeriesFromInts(vals, index, self.rhsLabel), nil
}

// sum method returns a Series that is the sum of each grouped result
func seriesGroupByResultCount(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("count", args, kwargs); err != nil {
		return nil, err
	}
	self := b.Receiver().(*SeriesGroupByResult)

	indexTexts := []string{}
	vals := []int{}

	sortedKeys := getSortedKeys(self.grouping)
	for _, groupName := range sortedKeys {
		series := self.grouping[groupName]
		count := series.Len()
		indexTexts = append(indexTexts, groupName)
		vals = append(vals, count)
	}

	index := NewIndex(indexTexts, self.lhsLabel)
	return newSeriesFromInts(vals, index, self.rhsLabel), nil
}

// apply method returns a Series that is built by calling the given
// function, and passing each grouped series as an argument to it
func seriesGroupByResultApply(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		funcVal starlark.Value
		self    = b.Receiver().(*SeriesGroupByResult)
	)

	if err := starlark.UnpackArgs("apply", args, kwargs,
		"function", &funcVal,
	); err != nil {
		return nil, err
	}

	funcObj, ok := funcVal.(*starlark.Function)
	if !ok {
		return nil, fmt.Errorf("first argument must be a function")
	}

	sortedKeys := getSortedKeys(self.grouping)
	builder := newTypedSliceBuilder(len(sortedKeys))
	indexNames := make([]string, len(sortedKeys))
	// NOTE: The index is either copied from the original DataFrame if the
	// `apply` function is returning Series's, or it is constructed from the
	// group names if the `apply` function is returning scalars.
	// TODO(dustmop): What happens if `apply` returns a scalar AND a Series?
	indexType := indexTypeNone

	for i, groupName := range sortedKeys {
		series := self.grouping[groupName]
		arguments := starlark.Tuple{series}
		// Call function, passing the series to it
		res, err := starlark.Call(thread, funcObj, arguments, nil)
		if err != nil {
			return nil, err
		}

		// Maybe a series
		if series, ok := res.(*Series); ok {
			if indexType == indexTypeNone {
				indexType = indexTypeFromOriginal
			} else if indexType != indexTypeFromOriginal {
				// TODO(dustmop): Test for this case
				return starlark.None, fmt.Errorf("TODO(dustmop): apply callback can only return scalars or Series's, not both")
			}
			for i := 0; i < series.Len(); i++ {
				builder.push(series.At(i))
			}
			continue
		}

		obj, ok := toScalarMaybe(res)
		if !ok {
			return nil, fmt.Errorf("apply's callback, could not convert from %T : %v", res, res)
		}
		// Accumulate the new series, and build the new index
		builder.push(obj)
		if indexType == indexTypeNone {
			indexType = indexTypeBuildNew
		} else if indexType != indexTypeBuildNew {
			// TODO(dustmop): Test for this case
			return starlark.None, fmt.Errorf("TODO(dustmop): apply callback can only return scalars or Series's, not both")
		}
		indexNames[i] = groupName
	}
	if err := builder.error(); err != nil {
		return nil, err
	}
	index := self.dfIndex
	if indexType == indexTypeBuildNew {
		index = NewIndex(indexNames, self.lhsLabel)
	}
	s := builder.toSeries(index, self.rhsLabel)
	return &s, nil
}

func getSortedKeys(m map[string]*Series) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
