package dataframe

import (
	"fmt"
	"strconv"
	"sort"

	"go.starlark.net/starlark"
)

type SeriesGroupByResult struct {
	gbLabel  string
	rhsLabel string
	grouping map[string][]string
}

// Freeze ...
func (sgbr *SeriesGroupByResult) Freeze() {
	// pass
}

func (sgbr *SeriesGroupByResult) Hash() (uint32, error) {
	// TODO
	return 0, nil
}

func (sgbr *SeriesGroupByResult) String() string {
	return fmt.Sprintf("<class 'SeriesGroupByResult'>")
}

// Truth ...
func (sgbr *SeriesGroupByResult) Truth() starlark.Bool {
	return true
}

func (sgbr *SeriesGroupByResult) Type() string {
	return "dataframe.SeriesGroupByResult"
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (sgbr *SeriesGroupByResult) Attr(name string) (starlark.Value, error) {
	if name == "sum" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return sgbr.seriesGroupByResultSum(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	return nil, nil
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (sgbr *SeriesGroupByResult) AttrNames() []string {
	return []string{"sum"}
}

func (sgbr *SeriesGroupByResult) seriesGroupByResultSum(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("sum", args, kwargs); err != nil {
		return nil, err
	}

	//result := []Series{}

	//newColumns := []string{sgbr.gbLabel, sgbr.rhsLabel}

	//tmpDf := &DataFrame{}
	//makeRows := newRowCollectOfSize(tmpDf, len(newColumns))

	index := []string{}
	vals := []int{}

	sortedKeys := getSortedKeys(sgbr.grouping)
	for _, groupName := range sortedKeys { // groupName, series := range sgbr.grouping {
		series := sgbr.grouping[groupName]

		//fmt.Printf("group = %s\n", groupName)

		sum := 0
		for _, elem := range series {
			//fmt.Printf("elem = %v\n", elem)
			num, err := strconv.Atoi(elem)
			if err == nil {
				sum += num
			}
		}
		//fmt.Printf("sum = %d\n\n", sum)
		//var vals []string
		//vals = []string{groupName, strconv.Itoa(sum)}
		//series := newSeriesFromStrings(vals, nil, "")
		//result = append(result, *series)

/*
		items := make([]interface{}, 2)
		items[0] = groupName
		items[1] = strconv.Itoa(sum)
		t := &rowTuple{0, tmpDf, items}
		makeRows.Push(t)
*/

		index = append(index, groupName)
		vals = append(vals, sum)

	}

	//if len(newColumns) != len(result) {
	//	fmt.Printf("num column names: %d\n", len(newColumns))
	//	fmt.Printf("columns in body:  %d\n", len(result))
	//	panic("stopping")
	//}

	//return &DataFrame{columnNames: newColumns, body: makeRows.Body()}, nil

	//return &Series{
		//columnNames: newColumns,
		//body: makeRows.Body(),
	//	valObjs: vals,
	//	index: index,
	//}, nil

	// TODO: sgbr.gbLabel as the index.name
	return newSeriesFromInts(vals, index, sgbr.rhsLabel), nil
}

func getSortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
