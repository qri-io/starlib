package dataframe

import (
	"encoding/csv"
	"fmt"
	"io"
	//"strconv"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

const (
	// Name of the module
	Name = "dataframe"
	// ModuleName is the filename of this module for the loader
	ModuleName = "dataframe.star"
)

// Module exposes the dataframe module
var Module = &starlarkstruct.Module{
	Name: Name,
	Members: starlark.StringDict{
		"read_csv":  starlark.NewBuiltin("read_csv", readCsv),
		"Index":     starlark.NewBuiltin("Index", newIndex),
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrame),
		"Series":    starlark.NewBuiltin("Series", newSeries),
	},
}

func unfinishedError(v starlark.Value, msg string) error {
	return fmt.Errorf("%s %s unfinished implementation: %s", Name, v.Type(), msg)
}

type DataFrame struct {
	frozen bool
	// TODO: This shoud be an Index
	columnNames []string
	body        []Series
	// TODO: This shoud be an Index
	index []string
}

// compile-time interface assertions
var (
	_ starlark.Value   = (*DataFrame)(nil)
	_ starlark.Mapping = (*DataFrame)(nil)
)

func readCsv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var content starlark.Value

	if err := starlark.UnpackArgs("readCsv", args, kwargs,
		"content", &content,
	); err != nil {
		return nil, err
	}

	text, ok := content.(starlark.String)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}

	reader := csv.NewReader(strings.NewReader(string(text)))

	// Assume header row
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}
	header := record

	// Body rows
	rowLength := -1
	var csvData [][]string
	for lineNum := 0; ; lineNum++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if rowLength == -1 {
			rowLength = len(record)
		} else if rowLength != len(record) {
			return nil, fmt.Errorf("rows must be same length, line %d is %d instead of %d", lineNum, len(record), rowLength)
		}
		csvData = append(csvData, record)
	}

	newBody := transposeToSeriesList(csvData, rowLength)
	return &DataFrame{
		columnNames: header,
		body:        newBody,
	}, nil
}

func newDataFrame(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		dataVal, indexVal, columnsVal, dtypeVal starlark.Value
		kopyVal                                 starlark.Bool
	)
	if err := starlark.UnpackArgs("DataFrame", args, kwargs,
		"data?", &dataVal,
		"index?", &indexVal,
		"columns?", &columnsVal,
		"dtype?", &dtypeVal,
		"copy?", &kopyVal,
	); err != nil {
		return nil, err
	}

	columns := toStrListOrNil(columnsVal)
	index := toStrListOrNil(indexVal)

	if dataDict, ok := dataVal.(*starlark.Dict); ok {
		// data is dict
		newBody := make([]Series, 0)
		keyList := make([]string, 0, dataDict.Len())

		inKeys := dataDict.Keys()
		for i := 0; i < len(inKeys); i++ {
			inKey := inKeys[i]
			val, _, _ := dataDict.Get(inKey)
			keyList = append(keyList, toStr(inKey))
			valList := toStrListOrNil(val)
			// TODO: Handle `list == nil`
			series := Series{which: typeObj, valObjs: valList}
			newBody = append(newBody, series)
		}

		// TODO: `index` will re-index the columns
		return &DataFrame{
			columnNames: keyList,
			body:        newBody,
		}, nil
	}

	if list, ok := dataVal.(*starlark.List); ok {
		// data is list
		collectRows := make([][]string, 0, list.Len())
		for i := 0; i < list.Len(); i++ {
			row := toStrListOrNil(list.Index(i))
			// TODO: Handle `row == nil`
			// TODO: Assert rows are the same length
			collectRows = append(collectRows, row)
		}
		newBody := transposeToSeriesList(collectRows, len(collectRows[0]))
		return &DataFrame{
			columnNames: columns,
			index:       index,
			body:        newBody,
		}, nil
	}

	return nil, fmt.Errorf("Not implemented, constructing DataFrame using %s", dataVal.Type())
}

func (df *DataFrame) numRows() int {
	if len(df.body) == 0 {
		return 0
	}
	return df.body[0].len()
}

// String implements the Stringer interface.
func (df *DataFrame) String() string {

	// Get width of the left-hand label
	labelWidth := 0
	if df.index == nil {
		bodyHeight := len(df.body)
		k := len(fmt.Sprintf("%d", bodyHeight))
		if k > labelWidth {
			labelWidth = k
		}
	} else {
		for _, str := range df.index {
			k := len(str)
			if k > labelWidth {
				labelWidth = k
			}
		}
	}

	// Create array of max widths, starting at 0
	widths := make([]int, len(df.columnNames))
	for i, name := range df.columnNames {
		w := len(fmt.Sprintf("%s", name))
		if w > widths[i] {
			widths[i] = w
		}
	}
	for i := 0; i < df.numRows(); i++ {
		for j, col := range df.body {
			elem := col.strAt(i)
			w := len(elem)
			if w > widths[j] {
				widths[j] = w
			}
		}
	}

	// Render the column names
	header := make([]string, 0, len(df.columnNames))
	for i, name := range df.columnNames {
		tmpl := fmt.Sprintf("%%%ds", widths[i])
		header = append(header, fmt.Sprintf(tmpl, name))
	}
	padding := strings.Repeat(" ", labelWidth)
	answer := fmt.Sprintf("%s    %s\n", padding, strings.Join(header, "  "))

	// Render each row
	for i := 0; i < df.numRows(); i++ {
		//for k, row := range df.body {
		render := []string{""}
		// Render the index number or label to start the line
		if df.index == nil {
			tmpl := fmt.Sprintf("%%%dd  ", labelWidth)
			render[0] = fmt.Sprintf(tmpl, i)
		} else {
			tmpl := fmt.Sprintf("%%%ds  ", labelWidth)
			render[0] = fmt.Sprintf(tmpl, df.index[i])
		}
		// Render each element of the row
		for j, col := range df.body {
			elem := col.strAt(i)
			tmpl := fmt.Sprintf("%%%ds", widths[j])
			render = append(render, fmt.Sprintf(tmpl, elem))
		}
		answer += strings.Join(render, "  ") + "\n"
	}

	return answer
}

// Type returns a short string describing the value's type.
func (DataFrame) Type() string { return fmt.Sprintf("%s.DataFrame", Name) }

// Freeze renders DataFrame immutable. required by starlark.Value interface
func (df *DataFrame) Freeze() { df.frozen = true }

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (df *DataFrame) Hash() (uint32, error) {
	// TODO (b5) - finish
	return 0, nil
}

// Truth reports whether the DataFrame is non-zero.
func (df *DataFrame) Truth() starlark.Bool {
	// TODO (b5) - finish
	return true
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (df *DataFrame) Attr(name string) (starlark.Value, error) {
	if name == "columns" {
		return &Index{texts: df.columnNames}, nil
	}
	if name == "apply" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameApply(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "head" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameHead(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "groupby" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameGroupBy(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "drop_duplicates" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameDropDuplicates(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "merge" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameMerge(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "reset_index" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return df.dataFrameResetIndex(thread, b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	return nil, nil
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	return []string{"columns", "apply", "head", "groupby", "drop_duplicates", "merge", "reset_index"}
}

var _ starlark.HasAttrs = (*DataFrame)(nil)
var _ starlark.HasSetField = (*DataFrame)(nil)
var _ starlark.HasSetKey = (*DataFrame)(nil)

func (df *DataFrame) SetField(name string, val starlark.Value) error {
	if name == "columns" {
		idx, ok := val.(*Index)
		if !ok {
			return fmt.Errorf("cannot assign to 'columns', wrong type")
		}
		df.columnNames = idx.texts
		return nil
	}
	return starlark.NoSuchAttrError(name)
}

func (df *DataFrame) SetKey(nameVal, val starlark.Value) error {
	name, ok := toStrMaybe(nameVal)
	if !ok {
		return fmt.Errorf("SetKey: name must be string")
	}

	// Figure out if a column already exists with the given name
	columnIndex := findKeyPos(name, df.columnNames)

	// Either prepend the new column, or keep the names the same
	newNames := make([]string, 0, len(df.columnNames)+1)
	if columnIndex == -1 {
		newNames = append([]string{name}, df.columnNames...)
	} else {
		newNames = df.columnNames
	}

	// Assignment of a scalar int
	if num, ok := toIntMaybe(val); ok {
		var newBody []Series
		newCol := newSeriesFromRepeatScalar(num, max(1, df.numRows()))
		if columnIndex == -1 {
			newBody = append([]Series{*newCol}, df.body...)
		} else {
			newBody = df.body
			newBody[columnIndex] = *newCol
		}
		df.columnNames = newNames
		df.body = newBody
		return nil
	}

	// TODO: Float, need test

	// Assignment of a scalar string
	if text, ok := toStrMaybe(val); ok {
		var newBody []Series
		newCol := newSeriesFromRepeatScalar(string(text), max(1, df.numRows()))
		if columnIndex == -1 {
			newBody = append([]Series{*newCol}, df.body...)
		} else {
			newBody = df.body
			newBody[columnIndex] = *newCol
		}
		df.columnNames = newNames
		df.body = newBody
		return nil
	}

	// Assignment of a Series to the column
	series, ok := val.(*Series)
	if !ok {
		return fmt.Errorf("SetKey: val must be int, string, or Series")
	}
	if df.numRows() > 0 && (df.numRows() != series.len()) {
		return fmt.Errorf("SetKey: val len must match number of rows")
	}

	var newBody []Series
	if columnIndex == -1 {
		newBody = append([]Series{*series}, df.body...)
	} else {
		newBody = df.body
		newBody[columnIndex] = *series
	}
	df.columnNames = newNames
	df.body = newBody
	return nil
}

// CompareSameType implements comparison of two DataFrame values. required by
// starlark.Comparable interface.
func (df *DataFrame) CompareSameType(op syntax.Token, v starlark.Value, depth int) (bool, error) {
	return false, unfinishedError(df, "CompareSameType")
}

// Binary implements binary operators, which satisfies the starlark.HasBinary
// interface
func (df *DataFrame) Binary(op syntax.Token, y starlark.Value, side starlark.Side) (starlark.Value, error) {
	return nil, unfinishedError(df, "Binary")
}

// Get returns a column of the DataFrame as a Series
func (df *DataFrame) Get(keyVal starlark.Value) (value starlark.Value, found bool, err error) {
	key, ok := toStrMaybe(keyVal)
	if !ok {
		return starlark.None, false, fmt.Errorf("Get key must be string")
	}

	// Find the column being retrieved, fail if not found
	keyPos := findKeyPos(key, df.columnNames)
	if keyPos == -1 {
		return starlark.None, false, fmt.Errorf("not found")
	}

	got := df.body[keyPos]
	// TODO: index should be the left-hand-side index, need a test
	index := []string{}

	return &Series{
		name:      key,
		which:     got.which,
		valInts:   got.valInts,
		valFloats: got.valFloats,
		valObjs:   got.valObjs,
		index:     index,
	}, true, nil
}

func (df *DataFrame) dataFrameApply(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		function, axis starlark.Value
	)

	if err := starlark.UnpackArgs("apply", args, kwargs,
		"function", &function,
		"axis?", &axis,
	); err != nil {
		return nil, err
	}

	axisNum, err := starlark.AsInt32(axis)
	if err != nil {
		return nil, err
	}
	if axisNum != 1 {
		return nil, fmt.Errorf("axis must equal 1 (row-size), other values Not Implemented")
	}

	funcObj, ok := function.(*starlark.Function)
	if !ok {
		return nil, fmt.Errorf("function must be a function")
	}

	var result []string
	for rows := newRowIter(df); !rows.Done(); rows.Next() {
		r := rows.GetRow()
		arguments := r.toTuple()
		res, err := starlark.Call(thread, funcObj, arguments, nil)
		if err != nil {
			return nil, err
		}

		text, ok := res.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("fn.apply should have returned String")
		}

		result = append(result, string(text))
	}

	return &Series{dtype: "object", which: typeObj, valObjs: result}, nil
}

func (df *DataFrame) dataFrameHead(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("head", args, kwargs); err != nil {
		return nil, err
	}

	// TODO: `n` parameter is number of rows to copy. Default to 5.

	numRows := 5
	newBody := make([]Series, 0, len(df.body))
	for k := 0; k < len(df.body); k++ {
		newBody = append(newBody, df.body[k].takeFirst(numRows))
	}

	return &DataFrame{
		columnNames: df.columnNames,
		body:        newBody,
	}, nil
}

func (df *DataFrame) dataFrameGroupBy(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		by starlark.Value
	)

	if err := starlark.UnpackArgs("groupby", args, kwargs,
		"by", &by,
	); err != nil {
		return nil, err
	}

	byList, ok := by.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("by should be a list of strings")
	}

	// TODO(dustmop): Support multiple values for the `by` value
	first := byList.Index(0)

	groupByStr, ok := first.(starlark.String)
	if !ok {
		return nil, fmt.Errorf("by[0] should be a string")
	}

	groupBy := string(groupByStr)

	result := map[string][]*rowTuple{}

	keyPos := findKeyPos(groupBy, df.columnNames)
	if keyPos == -1 {
		return starlark.None, nil
	}

	for rows := newRowIter(df); !rows.Done(); rows.Next() {
		r := rows.GetRow()
		groupValue := rows.GetStr(keyPos)
		result[groupValue] = append(result[groupValue], r)
	}

	return &GroupByResult{gbLabel: groupBy, columnNames: df.columnNames, grouping: result}, nil
}

func (df *DataFrame) dataFrameDropDuplicates(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		subset starlark.Value
	)

	if err := starlark.UnpackArgs("drop_duplicates", args, kwargs,
		"subset?", &subset,
	); err != nil {
		return nil, err
	}

	subsetPos := -1
	if subsetList, ok := subset.(*starlark.List); ok {
		// TODO: Assuming len 0
		elem := subsetList.Index(0)
		if text, ok := elem.(starlark.String); ok {
			subsetPos = findKeyPos(string(text), df.columnNames)
		}
	}

	seen := map[string]bool{}
	makeRows := newRowCollect(df)
	for rows := newRowIter(df); !rows.Done(); rows.Next() {
		matchOn := rows.Marshal(subsetPos)
		if seen[matchOn] {
			continue
		}
		seen[matchOn] = true
		makeRows.Push(rows.GetRow())
	}

	return &DataFrame{
		columnNames: df.columnNames,
		body:        makeRows.Body(),
	}, nil
}

type rowIndicies struct {
	//	first int
	//	key string
	is []int
}

func (df *DataFrame) dataFrameMerge(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		right, leftOn, rightOn, how starlark.Value
		suffixesVal                 starlark.Value
	)

	if err := starlark.UnpackArgs("merge", args, kwargs,
		"right", &right,
		"left_on?", &leftOn,
		"right_on?", &rightOn,
		"how?", &how,
		"suffixes?", &suffixesVal,
	); err != nil {
		return nil, err
	}

	var newColumns []string

	rightFrame, ok := right.(*DataFrame)
	if !ok {
		return starlark.None, fmt.Errorf("`right` must be a DataFrame")
	}

	leftOnStr := toStr(leftOn)
	rightOnStr := toStr(rightOn)
	leftKey := 0
	rightKey := 0
	if leftOnStr != "" {
		leftKey = findKeyPos(leftOnStr, df.columnNames)
		rightKey = findKeyPos(rightOnStr, rightFrame.columnNames)
		if leftKey == -1 {
			return starlark.None, fmt.Errorf("left key %q not found", leftOnStr)
		}
		if rightKey == -1 {
			return starlark.None, fmt.Errorf("right key %q not found", rightOnStr)
		}
	}

	howStr := toStrOrEmpty(how)

	var leftOrder []int
	if howStr == "" || howStr == "inner" {
		// For an inner merge, the keys appear with identical keys appearing together
		seen := make(map[string]int)
		idxs := make([]rowIndicies, 0)
		//for i, row := range df.body {
		for rows := newRowIter(df); !rows.Done(); rows.Next() {
			key := rows.GetStr(leftKey)
			n, has := seen[key]
			if has {
				idxs[n].is = append(idxs[n].is, rows.Index())
			} else {
				n = len(idxs)
				idxs = append(idxs, rowIndicies{is: []int{rows.Index()}})
				seen[key] = n
			}
		}
		// Collect the rows now based upon the desired order
		for _, numList := range idxs {
			for _, i := range numList.is {
				leftOrder = append(leftOrder, i)
			}
		}
	} else if howStr == "left" {
		leftOrder = nil
	} else {
		return starlark.None, fmt.Errorf("not implemented: `how` is %q", howStr)
	}

	var suffixes []string
	leftList := toStrListOrNil(suffixesVal)

	if len(leftList) == 2 {
		suffixes = leftList
	} else {
		suffixes = []string{"_x", "_y"}
	}

	// If column names of the merge key are the same, don't include the second one, ignore it
	ignore := findKeyPos(df.columnNames[leftKey], rightFrame.columnNames)

	leftColumns := modifyNames(df.columnNames, suffixes[0], leftKey)
	rightColumns := modifyNames(rightFrame.columnNames, suffixes[1], rightKey)
	if ignore != -1 {
		rightColumns = removeFromStringList(rightColumns, ignore)
	}
	newColumns = append(leftColumns, rightColumns...)

	makeRows := newRowCollectOfSize(df, len(newColumns))
	leftIter := newRowIterWithOrder(df, leftOrder)
	for ; !leftIter.Done(); leftIter.Next() {
		for rightIter := newRowIter(rightFrame); !rightIter.Done(); rightIter.Next() {
			newRow := leftIter.MergeWith(rightIter, leftKey, rightKey, ignore)
			if newRow != nil {
				makeRows.Push(newRow)
			}
		}
	}

	return &DataFrame{
		columnNames: newColumns,
		body:        makeRows.Body(),
	}, nil
}

func (df *DataFrame) dataFrameResetIndex(thread *starlark.Thread, fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("reset_index", args, kwargs); err != nil {
		return nil, err
	}

	if df.index == nil {
		return df, nil
	}

	newColumns := append([]string{"index"}, df.columnNames...)
	newBody := make([]Series, 0, len(df.body))

	newBody = append(newBody, Series{which: typeObj, valObjs: df.index})
	for _, col := range df.body {
		newBody = append(newBody, col)
	}

	return &DataFrame{
		columnNames: newColumns,
		body:        newBody,
	}, nil
}

func transposeToSeriesList(rows [][]string, rowLength int) []Series {
	colLength := len(rows)
	newBody := make([]Series, 0, rowLength)
	for i := 0; i < rowLength; i++ {
		// TODO: Detect types for each column, use that here.
		newCol := make([]string, 0, colLength)
		for j := 0; j < colLength; j++ {
			cell := rows[j][i]
			newCol = append(newCol, cell)
		}
		newBody = append(newBody, Series{
			which:   typeObj,
			valObjs: newCol,
		})
	}
	return newBody
}

func mergeMatchRows(left, right []string, leftKey, rightKey, ignore int) ([]string, bool) {
	leftElem := left[leftKey]
	rightElem := right[rightKey]
	if leftElem == rightElem {
		right = removeFromStringList(right, ignore)
		return append(left, right...), true
	}
	return nil, false
}

func modifyNames(names []string, suffix string, keyPos int) []string {
	result := make([]string, len(names))
	for i, elem := range names {
		if i == keyPos {
			result[i] = elem
		} else {
			result[i] = fmt.Sprintf("%s%s", elem, suffix)
		}
	}
	return result
}

func removeFromStringList(ls []string, i int) []string {
	if i == -1 {
		return ls
	}
	a := make([]string, len(ls))
	copy(a, ls)
	copy(a[i:], a[i+1:])
	a[len(a)-1] = ""
	return a[:len(a)-1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
