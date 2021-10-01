package dataframe

import (
	"fmt"
	"sort"
	"strconv"

	"go.starlark.net/starlark"
)

// SeriesGroupByResult is the result of indexing into a groupBy result
type SeriesGroupByResult struct {
	lhsLabel string
	rhsLabel string
	grouping map[string][]string
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*SeriesGroupByResult)(nil)
	_ starlark.HasAttrs = (*SeriesGroupByResult)(nil)
)

var seriesGroupByResultMethods = map[string]*starlark.Builtin{
	"count": starlark.NewBuiltin("count", seriesGroupByResultCount),
	"sum":   starlark.NewBuiltin("sum", seriesGroupByResultSum),
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
		for _, elem := range series {
			num, err := strconv.Atoi(elem)
			if err == nil {
				sum += num
			}
		}

		indexTexts = append(indexTexts, groupName)
		vals = append(vals, sum)
	}

	// TODO(dustmop): sgbr.lhsLabel as the index.name, add a test
	index := NewIndex(indexTexts, "")
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
		count := len(series)
		indexTexts = append(indexTexts, groupName)
		vals = append(vals, count)
	}

	// TODO(dustmop): sgbr.lhsLabel as the index.name, add a test
	index := NewIndex(indexTexts, "")
	return newSeriesFromInts(vals, index, self.rhsLabel), nil
}

func getSortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
