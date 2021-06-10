package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
)

type GroupByResult struct {
	columnNames []string
	gbLabel     string
	grouping    map[string][]*rowTuple
}

// Freeze ...
func (gbr *GroupByResult) Freeze() {
	// pass
}

func (gbr *GroupByResult) Hash() (uint32, error) {
	// TODO
	return 0, nil
}

func (gbr *GroupByResult) String() string {
	return fmt.Sprintf("<class 'GroupByResult'>")
}

// Truth ...
func (gbr *GroupByResult) Truth() starlark.Bool {
	return true
}

func (gbr *GroupByResult) Type() string {
	return "dataframe.GroupByResult"
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (gbr *GroupByResult) Attr(name string) (starlark.Value, error) {
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (gbr *GroupByResult) AttrNames() []string {
	return []string{}
}

func (gbr *GroupByResult) Get(key starlark.Value) (value starlark.Value, found bool, err error) {
	keyVal, ok := key.(starlark.String)
	if !ok {
		return nil, false, fmt.Errorf("index must be a string")
	}

	keyStr := string(keyVal)

	keyPos := findKeyPos(keyStr, gbr.columnNames)
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

	//newColumnNames := []string{gbr.gbLabel, keyStr}
	//return &DataFrame{columnNames: newColumnNames, body: result}, true, nil
	return &SeriesGroupByResult{gbLabel: gbr.gbLabel, rhsLabel: keyStr, grouping: result}, true, nil
}
