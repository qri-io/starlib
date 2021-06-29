package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
)

// GroupByResult is the result of using groupBy on a DataFrame
type GroupByResult struct {
	columns  *Index
	label    string
	grouping map[string][]*rowTuple
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*GroupByResult)(nil)
	_ starlark.Mapping  = (*GroupByResult)(nil)
	_ starlark.HasAttrs = (*GroupByResult)(nil)
)

// Freeze has no effect on the immutable GroupByResult
func (gbr *GroupByResult) Freeze() {
	// pass
}

// Hash cannot be used with GroupByResult
func (gbr *GroupByResult) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", gbr.Type())
}

// String returns a string representation of the GroupByResult
func (gbr *GroupByResult) String() string {
	return fmt.Sprintf("<%s>", gbr.Type())
}

// Truth converts the GroupByResult into a bool
func (gbr *GroupByResult) Truth() starlark.Bool {
	return true
}

// Type returns the type as a string
func (gbr *GroupByResult) Type() string {
	return fmt.Sprintf("%s.GroupByResult", Name)
}

// Attr gets a value for an attribute
func (gbr *GroupByResult) Attr(name string) (starlark.Value, error) {
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available attributes
func (gbr *GroupByResult) AttrNames() []string {
	return []string{}
}

// Get returns a series like object by indexing into the result of a groupBy call
func (gbr *GroupByResult) Get(key starlark.Value) (value starlark.Value, found bool, err error) {
	name, ok := toStrMaybe(key)
	if !ok {
		return nil, false, fmt.Errorf("key must be a string")
	}

	keyPos := findKeyPos(name, gbr.columns.texts)
	if keyPos == -1 {
		return nil, false, fmt.Errorf("not found")
	}

	result := map[string][]string{}
	for group, frame := range gbr.grouping {
		newRow := []string{}
		for _, row := range frame {
			val := row.strAt(keyPos)
			newRow = append(newRow, val)
		}
		result[group] = newRow
	}

	return &SeriesGroupByResult{lhsLabel: gbr.label, rhsLabel: name, grouping: result}, true, nil
}
