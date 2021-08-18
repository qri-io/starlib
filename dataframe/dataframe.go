package dataframe

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
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
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrameBuiltin),
		"Index":     starlark.NewBuiltin("Index", newIndex),
		"Series":    starlark.NewBuiltin("Series", newSeries),
	},
}

// DataFrame is the primary data structure of this package, it represents
// a column-oriented table of data, and provides spreadsheet and sql like
// functionality.
type DataFrame struct {
	frozen  bool
	columns *Index
	index   *Index
	body    []Series
}

// compile-time interface assertions
var (
	_ starlark.Value       = (*DataFrame)(nil)
	_ starlark.Mapping     = (*DataFrame)(nil)
	_ starlark.HasAttrs    = (*DataFrame)(nil)
	_ starlark.HasSetField = (*DataFrame)(nil)
	_ starlark.HasSetKey   = (*DataFrame)(nil)
)

var dataframeMethods = map[string]*starlark.Builtin{
	"apply":           starlark.NewBuiltin("apply", dataframeApply),
	"drop_duplicates": starlark.NewBuiltin("drop_duplicates", dataframeDropDuplicates),
	"groupby":         starlark.NewBuiltin("groupby", dataframeGroupBy),
	"head":            starlark.NewBuiltin("head", dataframeHead),
	"merge":           starlark.NewBuiltin("merge", dataframeMerge),
	"reset_index":     starlark.NewBuiltin("reset_index", dataframeResetIndex),
	"append":          starlark.NewBuiltin("append", dataframeAppend),
	"set_csv":         starlark.NewBuiltin("set_csv", dataframeSetCSV),
}

func readCsv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var content starlark.Value

	if err := starlark.UnpackArgs("read_csv", args, kwargs,
		"content", &content,
	); err != nil {
		return nil, err
	}

	text, ok := toStrMaybe(content)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}
	body, header, err := constructBodyHeaderFromCSV(text, true)
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: NewIndex(header, ""),
		body:    body,
	}, nil
}

// newDataFrameBuiltin constructs a dataframe, meant to be called from starlark
func newDataFrameBuiltin(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

	columns := toStrSliceOrNil(columnsVal)
	index, _ := toIndexMaybe(indexVal)

	return NewDataFrame(dataVal, columns, index)
}

// NewDataFrame constructs a DataFrame from data, and optionally column names and an index
// data can be any datatype supported by this DataFrame implementation, either
// go native types or starlark types
func NewDataFrame(data interface{}, columnNames []string, index *Index) (*DataFrame, error) {
	var (
		body    []Series
		columns *Index
		err     error
	)

	if columnNames != nil {
		columns = NewIndex(columnNames, "")
	}

	switch inData := data.(type) {
	case nil:
		body = []Series{}

	case []interface{}:
		body, err = constructBodyFromNativeSlice(inData)
		if err != nil {
			return nil, err
		}

	case [][]interface{}:
		body, err = constructBodyFromRows(inData)
		if err != nil {
			return nil, err
		}

	case Series:
		body = []Series{inData}

	case *Series:
		if inData == nil {
			return nil, fmt.Errorf("cannot convert nil series pointer into a dataframe body")
		}
		body = []Series{*inData}

	case *DataFrame:
		body = inData.body

	case *starlark.Dict:
		var keys []string
		body, keys, err = constructBodyFromStarlarkDict(inData)
		if err != nil {
			return nil, err
		}
		columns = NewIndex(keys, "")
		// TODO(dustmop): `index` will re-index

	case *starlark.List:
		body, err = constructBodyFromStarlarkList(inData)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Not implemented, constructing DataFrame using %s", reflect.TypeOf(data))
	}

	// Check that the index and columns, if present, match the body size
	numCols, numRows := sizeOfBody(body)
	if index.len() > 0 && index.len() != numRows {
		// TODO(dustmop): Add test
		return nil, fmt.Errorf("size of index does not match body size")
	}
	if columns.len() > 0 && columns.len() != numCols {
		// TODO(dustmop): Add test
		return nil, fmt.Errorf("number of columns does not match body size")
	}

	return &DataFrame{
		columns: columns,
		index:   index,
		body:    body,
	}, nil
}

// construct a body from a slice, either cooercing it into rows, or treating it as a column
func constructBodyFromNativeSlice(ls []interface{}) ([]Series, error) {
	if rows := toTwoDimensionalRows(ls); rows != nil {
		return constructBodyFromRows(rows)
	}
	// One dimensional list, treat it as a column
	builder := newTypedSliceBuilder(len(ls))
	for i := 0; i < len(ls); i++ {
		builder.push(ls[i])
	}
	if err := builder.error(); err != nil {
		return nil, err
	}
	s := builder.toSeries(nil, "")
	return []Series{s}, nil
}

// try to convert a slice into 2-d rows, returning nil if not possible
func toTwoDimensionalRows(ls []interface{}) [][]interface{} {
	build := [][]interface{}{}
	for i := 0; i < len(ls); i++ {
		elem := ls[i]
		items, ok := elem.([]interface{})
		if !ok {
			return nil
		}
		build = append(build, items)
	}
	return build
}

// construct a body from 2-d rows
func constructBodyFromRows(rows [][]interface{}) ([]Series, error) {
	rowLength := -1
	var builder *tableBuilder
	for lineNum, record := range rows {
		if rowLength == -1 {
			rowLength = len(record)
		} else if rowLength != len(record) {
			return nil, fmt.Errorf("rows must be same length, line %d is %d instead of %d", lineNum, len(record), rowLength)
		}
		if builder == nil {
			builder = newTableBuilder(rowLength, 0)
		}
		builder.pushRow(record)
	}
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return body, nil
}

// construct a body from a starlark.Dict
func constructBodyFromStarlarkDict(data *starlark.Dict) ([]Series, []string, error) {
	newBody := make([]Series, 0)
	keyList := make([]string, 0, data.Len())
	inKeys := data.Keys()
	numRows := -1
	for i := 0; i < len(inKeys); i++ {
		// Collect each key, use them as the default index
		inKey := inKeys[i]
		keyList = append(keyList, toStr(inKey))
		// Get each value, which should be a list of values
		val, _, _ := data.Get(inKey)
		items := toInterfaceSliceOrNil(val)
		if items == nil {
			return nil, nil, fmt.Errorf("invalid values for column")
		}
		// Validate that the size of each column is the same
		// TODO(dustmop): Add test
		if numRows == -1 {
			numRows = len(items)
		} else if numRows != len(items) {
			return nil, nil, fmt.Errorf("columns need to be the same length")
		}
		// The list of values should be of the same type
		builder := newTypedSliceBuilder(len(items))
		for _, it := range items {
			builder.push(it)
		}
		if err := builder.error(); err != nil {
			return nil, nil, err
		}
		newBody = append(newBody, builder.toSeries(nil, ""))
	}
	return newBody, keyList, nil
}

// construct a body from a starlark.List
func constructBodyFromStarlarkList(data *starlark.List) ([]Series, error) {
	numRows := data.Len()
	numCols := -1
	var builder *tableBuilder
	// Iterate the input data row-size
	for y := 0; y < data.Len(); y++ {
		row := toInterfaceSliceOrNil(data.Index(y))
		if row == nil {
			obj, _ := json.Marshal(data)
			return nil, fmt.Errorf("invalid values for row: %s", string(obj))
		}
		// Validate that the size of each row is the same
		// TODO(dustmop): Add test
		if numCols == -1 {
			numCols = len(row)
		} else if numCols != len(row) {
			return nil, fmt.Errorf("rows need to be the same length")
		}
		if builder == nil {
			builder = newTableBuilder(numCols, numRows)
		}
		builder.pushRow(row)
	}
	newBody, err := builder.body()
	if err != nil {
		return nil, err
	}
	return newBody, nil
}

func constructBodyHeaderFromCSV(text string, hasHeader bool) ([]Series, []string, error) {
	reader := csv.NewReader(ReplaceReader(strings.NewReader(text)))

	header := []string{}
	if hasHeader {
		record, err := reader.Read()
		if err != nil {
			return nil, nil, err
		}
		header = record
	}

	// Body rows
	rowLength := -1
	var builder *tableBuilder
	for lineNum := 0; ; lineNum++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if rowLength == -1 {
			rowLength = len(record)
		} else if rowLength != len(record) {
			return nil, nil, fmt.Errorf("rows must be same length, line %d is %d instead of %d", lineNum, len(record), rowLength)
		}
		if builder == nil {
			builder = newTableBuilder(rowLength, 0)
		}
		builder.pushTextRow(record)
	}
	body, err := builder.body()
	if err != nil {
		return nil, nil, err
	}
	return body, header, nil
}

// get the size of the body, width and height
func sizeOfBody(body []Series) (int, int) {
	if len(body) == 0 {
		return 0, 0
	}
	return len(body), body[0].len()
}

// NumCols returns the number of columns
func (df *DataFrame) NumCols() int {
	return len(df.body)
}

// NumRows returns the number of rows
func (df *DataFrame) NumRows() int {
	if len(df.body) == 0 {
		return 0
	}
	return df.body[0].len()
}

// Row returns the ith row as a slice of go native types
func (df *DataFrame) Row(i int) []interface{} {
	if i >= df.NumRows() {
		return nil
	}
	row := make([]interface{}, df.NumCols())
	for k := 0; k < df.NumCols(); k++ {
		series := df.body[k]
		row[k] = series.At(i)
	}
	return row
}

// String returns the DataFrame as a string in a readable, tabular form
func (df *DataFrame) String() string {
	return df.stringify()
}

// Type returns a short string describing the value's type.
func (DataFrame) Type() string { return fmt.Sprintf("%s.DataFrame", Name) }

// Freeze renders DataFrame immutable.
func (df *DataFrame) Freeze() { df.frozen = true }

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (df *DataFrame) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", df.Type())
}

// Truth reports whether the DataFrame is non-zero.
func (df *DataFrame) Truth() starlark.Bool {
	// NOTE: In python, calling bool(DataFrame) raises this exception: "ValueError: The truth
	// value of a DataFrame is ambiguous. Use a.empty, a.bool(), a.item(), a.any() or a.all()."
	// Since starlark does not have exceptions, just always return true.
	return true
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (df *DataFrame) Attr(name string) (starlark.Value, error) {
	switch name {
	case "columns":
		return df.columns, nil
	}
	return builtinAttr(df, name, dataframeMethods)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	methodNames := builtinAttrNames(seriesMethods)
	return append([]string{"columns"}, methodNames...)
}

// SetField assigns to a field of the DataFrame
func (df *DataFrame) SetField(name string, val starlark.Value) error {
	if df.frozen {
		return fmt.Errorf("cannot set, DataFrame is frozen")
	}

	if name == "columns" {
		idx, ok := val.(*Index)
		if !ok {
			return fmt.Errorf("cannot assign to 'columns', wrong type")
		}
		df.columns = idx
		return nil
	}
	return starlark.NoSuchAttrError(name)
}

// SetKey assigns a value to a DataFrame at the given key
func (df *DataFrame) SetKey(nameVal, val starlark.Value) error {
	if df.frozen {
		return fmt.Errorf("cannot set, DataFrame is frozen")
	}

	name, ok := toStrMaybe(nameVal)
	if !ok {
		return fmt.Errorf("SetKey: name must be string")
	}

	// Figure out if a column already exists with the given name
	columnIndex := findKeyPos(name, df.columns.texts)

	// Either prepend the new column, or keep the names the same
	newNames := make([]string, 0, len(df.columns.texts)+1)
	if columnIndex == -1 {
		newNames = append([]string{name}, df.columns.texts...)
	} else {
		newNames = df.columns.texts
	}

	// Assignment of a scalar (int, bool, float, string) to the column
	if scalar, ok := toScalarMaybe(val); ok {
		var newBody []Series
		newCol := newSeriesFromRepeatScalar(scalar, max(1, df.NumRows()))
		if columnIndex == -1 {
			// New columns are added to the left side of the dataframe
			newBody = append([]Series{*newCol}, df.body...)
		} else {
			newBody = df.body
			newBody[columnIndex] = *newCol
		}
		df.columns = NewIndex(newNames, "")
		df.body = newBody
		return nil
	}

	// Assignment of a Series to the column
	series, ok := val.(*Series)
	if !ok {
		return fmt.Errorf("SetKey: val must be int, string, or Series")
	}
	if df.NumRows() > 0 && (df.NumRows() != series.len()) {
		return fmt.Errorf("SetKey: val len must match number of rows")
	}

	var newBody []Series
	if columnIndex == -1 {
		newBody = append([]Series{*series}, df.body...)
	} else {
		newBody = df.body
		newBody[columnIndex] = *series
	}
	df.columns = NewIndex(newNames, "")
	df.body = newBody
	return nil
}

// Get returns a column of the DataFrame as a Series
func (df *DataFrame) Get(keyVal starlark.Value) (value starlark.Value, found bool, err error) {
	key, ok := toStrMaybe(keyVal)
	if !ok {
		return starlark.None, false, fmt.Errorf("Get key must be string")
	}

	// Find the column being retrieved, fail if not found
	keyPos := findKeyPos(key, df.columns.texts)
	if keyPos == -1 {
		return starlark.None, false, fmt.Errorf("DataFrame.Get: key not found %q", key)
	}

	got := df.body[keyPos]
	// TODO(dustmop): index should be the left-hand-side index, need a test
	index := NewIndex(nil, "")

	dtype := got.dtype
	if dtype == "" {
		switch got.which {
		case typeInt:
			dtype = "int64"
		case typeFloat:
			dtype = "float64"
		default:
			dtype = "object"
		}
	}

	return &Series{
		name:      key,
		dtype:     dtype,
		which:     got.which,
		valInts:   got.valInts,
		valFloats: got.valFloats,
		valObjs:   got.valObjs,
		index:     index,
	}, true, nil
}

func (df *DataFrame) stringify() string {
	// Get width of the left-hand label
	labelWidth := 0
	if df.index == nil {
		bodyHeight := df.NumRows()
		k := len(fmt.Sprintf("%d", bodyHeight))
		if k > labelWidth {
			labelWidth = k
		}
	} else {
		for _, str := range df.index.texts {
			k := len(str)
			if k > labelWidth {
				labelWidth = k
			}
		}
	}

	// Create array of max widths, starting at 0
	widths := make([]int, df.NumCols())
	colTexts := []string{}
	if df.columns != nil {
		colTexts = df.columns.texts
	}
	for i, name := range colTexts {
		w := len(fmt.Sprintf("%s", name))
		if w > widths[i] {
			widths[i] = w
		}
	}
	for i := 0; i < df.NumRows(); i++ {
		for j, col := range df.body {
			elem := col.strAt(i)
			w := len(elem)
			if w > widths[j] {
				widths[j] = w
			}
		}
	}

	// Render columns
	header := make([]string, 0, len(colTexts))
	if len(colTexts) > 0 {
		// Render the column names
		for i, name := range colTexts {
			tmpl := fmt.Sprintf("%%%ds", widths[i])
			header = append(header, fmt.Sprintf(tmpl, name))
		}
	} else {
		// Render the column indicies
		for i := range df.body {
			tmpl := fmt.Sprintf("%%%dd", widths[i])
			header = append(header, fmt.Sprintf(tmpl, i))
		}
	}
	padding := strings.Repeat(" ", labelWidth)
	answer := fmt.Sprintf("%s    %s\n", padding, strings.Join(header, "  "))

	// Render each row
	for i := 0; i < df.NumRows(); i++ {
		render := []string{""}
		// Render the index number or label to start the line
		if df.index == nil {
			tmpl := fmt.Sprintf("%%%dd  ", labelWidth)
			render[0] = fmt.Sprintf(tmpl, i)
		} else {
			tmpl := fmt.Sprintf("%%%ds  ", labelWidth)
			render[0] = fmt.Sprintf(tmpl, df.index.texts[i])
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

// apply method iterates the rows of a DataFrame, calls the given function for
// each row, creating a single Series of the results
func dataframeApply(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		funcVal, axisVal starlark.Value
	)
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("apply", args, kwargs,
		"function", &funcVal,
		"axis?", &axisVal,
		// TODO(dustmop): Add other arguments that pandas.DataFrame.apply has.
	); err != nil {
		return nil, err
	}

	axis, ok := toIntMaybe(axisVal)
	if !ok || axis != 1 {
		return nil, fmt.Errorf("axis must equal 1 for row-wise application")
	}

	funcObj, ok := funcVal.(*starlark.Function)
	if !ok {
		return nil, fmt.Errorf("first argument must be a function")
	}

	var result []string
	for rows := newRowIter(self); !rows.Done(); rows.Next() {
		r := rows.GetRow()
		arguments := r.toTuple()
		res, err := starlark.Call(thread, funcObj, arguments, nil)
		if err != nil {
			return nil, err
		}

		// TODO(dustmop): Accept other return value types.
		text, ok := res.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("fn.apply should have returned String")
		}

		result = append(result, string(text))
	}

	return &Series{dtype: "object", which: typeObj, valObjs: result}, nil
}

// head method returns a copy of the DataFrame but only with the first n rows
func dataframeHead(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var nVal starlark.Value

	if err := starlark.UnpackArgs("head", args, kwargs,
		"n?", &nVal,
	); err != nil {
		return nil, err
	}
	self := b.Receiver().(*DataFrame)

	numRows, ok := toIntMaybe(nVal)
	if !ok {
		// n defaults to 5 if not given
		numRows = 5
	}

	builder := newTableBuilder(self.NumCols(), 0)
	for rows := newRowIter(self); !rows.Done() && numRows > 0; rows.Next() {
		r := rows.GetRow()
		builder.pushRow(r.data)
		numRows--
	}

	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: self.columns,
		body:    body,
	}, nil
}

// groupby method returns a grouped set of rows collected by some given column value
func dataframeGroupBy(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var by starlark.Value

	if err := starlark.UnpackArgs("groupby", args, kwargs,
		"by", &by,
	); err != nil {
		return nil, err
	}
	self := b.Receiver().(*DataFrame)

	byList, ok := by.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("by should be a list of strings")
	}

	// TODO(dustmop): Support multiple values for the `by` value
	first := byList.Index(0)
	groupBy, ok := toStrMaybe(first)
	if !ok {
		return nil, fmt.Errorf("by[0] should be a string")
	}

	result := map[string][]*rowTuple{}
	keyPos := findKeyPos(groupBy, self.columns.texts)
	if keyPos == -1 {
		return starlark.None, nil
	}

	for rows := newRowIter(self); !rows.Done(); rows.Next() {
		r := rows.GetRow()
		groupValue := rows.GetStr(keyPos)
		result[groupValue] = append(result[groupValue], r)
	}

	return &GroupByResult{label: groupBy, columns: self.columns, grouping: result}, nil
}

// drop_duplicates method returns a copy of a DataFrame without duplicates
func dataframeDropDuplicates(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		subset starlark.Value
	)
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("drop_duplicates", args, kwargs,
		"subset?", &subset,
	); err != nil {
		return nil, err
	}

	// TODO(dustmop): Support multiple values for the `subset` value
	subsetPos := -1
	if subsetList, ok := subset.(*starlark.List); ok {
		// TODO: Assuming len > 0
		elem := subsetList.Index(0)
		if text, ok := elem.(starlark.String); ok {
			subsetPos = findKeyPos(string(text), self.columns.texts)
		}
	}

	seen := map[string]bool{}
	builder := newTableBuilder(self.NumCols(), 0)
	for rows := newRowIter(self); !rows.Done(); rows.Next() {
		matchOn := rows.Marshal(subsetPos)
		if seen[matchOn] {
			continue
		}
		seen[matchOn] = true
		builder.pushRow(rows.GetRow().data)
	}

	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: self.columns,
		body:    body,
	}, nil
}

// merge method merges the rows of two DataFrames
func dataframeMerge(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		right, leftOn, rightOn, how starlark.Value
		suffixesVal                 starlark.Value
	)
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("merge", args, kwargs,
		"right", &right,
		"left_on?", &leftOn,
		"right_on?", &rightOn,
		"how?", &how,
		"suffixes?", &suffixesVal,
	); err != nil {
		return nil, err
	}

	rightFrame, ok := right.(*DataFrame)
	if !ok {
		return starlark.None, fmt.Errorf("`right` must be a DataFrame")
	}

	leftOnStr := toStr(leftOn)
	rightOnStr := toStr(rightOn)
	leftKey := 0
	rightKey := 0
	if leftOnStr != "" {
		leftKey = findKeyPos(leftOnStr, self.columns.texts)
		rightKey = findKeyPos(rightOnStr, rightFrame.columns.texts)
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
		// Indices are collected using a list of list of ints. Each of the dataframe
		// rows with the same key will have their indices appear adjacent to each other.
		// For example, when running the test dataframe_merge.star, for the first
		// call to `df1.merge`, the `idxs` array will be [[0, 3], [1], [2]], caused by
		// the key "foo" appearing at positions 0 and 3: its rows end up together.
		idxs := make([][]int, 0)
		for rows := newRowIter(self); !rows.Done(); rows.Next() {
			key := rows.GetStr(leftKey)
			n, has := seen[key]
			if has {
				idxs[n] = append(idxs[n], rows.Index())
			} else {
				n = len(idxs)
				idxs = append(idxs, []int{rows.Index()})
				seen[key] = n
			}
		}
		// Collect the rows now based upon the desired order
		for _, numList := range idxs {
			leftOrder = append(leftOrder, numList...)
		}
	} else if howStr == "left" {
		leftOrder = nil
	} else {
		return starlark.None, fmt.Errorf("not implemented: `how` is %q", howStr)
	}

	var suffixes []string
	leftList := toStrSliceOrNil(suffixesVal)

	// TODO(dustmop): Ensure suffixes are the right length, add a test
	if len(leftList) == 2 {
		suffixes = leftList
	} else {
		// Default column indicies are "_x" and "_y"
		suffixes = []string{"_x", "_y"}
	}

	// If column names of the merge key are the same, don't include the second one, ignore it
	ignore := findKeyPos(self.columns.texts[leftKey], rightFrame.columns.texts)

	leftColumns := addSuffixToStringList(self.columns.texts, suffixes[0], leftKey)
	rightColumns := addSuffixToStringList(rightFrame.columns.texts, suffixes[1], rightKey)
	if ignore != -1 {
		rightColumns = removeElemFromStringList(rightColumns, ignore)
	}
	newColumns := append(leftColumns, rightColumns...)

	builder := newTableBuilder(len(newColumns), 0)
	leftIter := newRowIterWithOrder(self, leftOrder)
	for ; !leftIter.Done(); leftIter.Next() {
		for rightIter := newRowIter(rightFrame); !rightIter.Done(); rightIter.Next() {
			newRow := leftIter.MergeWith(rightIter, leftKey, rightKey, ignore)
			if newRow != nil {
				builder.pushRow(newRow.data)
			}
		}
	}

	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: NewIndex(newColumns, ""),
		body:    body,
	}, nil
}

// reset_index method turns the DataFrame index into a new column
func dataframeResetIndex(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("reset_index", args, kwargs); err != nil {
		return nil, err
	}
	self := b.Receiver().(*DataFrame)

	if self.index == nil {
		return self, nil
	}

	newColumns := append([]string{"index"}, self.columns.texts...)
	newBody := make([]Series, 0, self.NumCols())

	newBody = append(newBody, Series{which: typeObj, valObjs: self.index.texts})
	for _, col := range self.body {
		newBody = append(newBody, col)
	}

	return &DataFrame{
		columns: NewIndex(newColumns, ""),
		body:    newBody,
	}, nil
}

// set_csv parses csv data and assigns it to the body
func dataframeSetCSV(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var textVal starlark.String

	if err := starlark.UnpackArgs("text", args, kwargs,
		"text", &textVal,
	); err != nil {
		return nil, err
	}
	self := b.Receiver().(*DataFrame)

	body, _, err := constructBodyHeaderFromCSV(string(textVal), false)
	if err != nil {
		return nil, err
	}

	self.body = body
	return starlark.None, nil
}

// append adds a new row to the body
func dataframeAppend(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var otherVal starlark.Value
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("append", args, kwargs,
		"other", &otherVal,
	); err != nil {
		return nil, err
	}

	dataCols := -1
	var data [][]interface{}
	ls, ok := otherVal.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("append requires a list to append")
	}
	for i := 0; i < ls.Len(); i++ {
		elem := ls.Index(i)
		arr := toInterfaceSliceOrNil(elem)
		if arr == nil {
			return nil, fmt.Errorf("append requires a list of lists")
		}
		if dataCols == -1 {
			dataCols = len(arr)
		}
		data = append(data, arr)
	}
	dataRows := len(data)

	newBody := make([]Series, len(self.body))
	for x := 0; x < dataCols; x++ {
		col := self.body[x]
		builder := newTypedSliceBuilderFromSeries(&col)
		for y := 0; y < dataRows; y++ {
			builder.push(data[y][x])
		}
		if err := builder.error(); err != nil {
			return nil, err
		}
		newBody[x] = builder.toSeries(nil, "")
	}

	return &DataFrame{
		columns: self.columns,
		index:   self.index,
		body:    newBody,
	}, nil
}

func addSuffixToStringList(names []string, suffix string, keyPos int) []string {
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

func removeElemFromStringList(ls []string, i int) []string {
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
