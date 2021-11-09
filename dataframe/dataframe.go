package dataframe

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
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
		"parse_csv": starlark.NewBuiltin("parse_csv", parseCsv),
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrameBuiltin),
		"Index":     starlark.NewBuiltin("Index", newIndex),
		"Series":    starlark.NewBuiltin("Series", newSeries),
	},
}

// DataFrame is the primary data structure of this package, it represents
// a column-oriented table of data, and provides spreadsheet and sql like
// functionality.
type DataFrame struct {
	frozen     bool
	columns    *Index
	index      *Index
	body       []Series
	stringConf stringConfig
}

// compile-time interface assertions
var (
	_ starlark.Value       = (*DataFrame)(nil)
	_ starlark.Mapping     = (*DataFrame)(nil)
	_ starlark.HasAttrs    = (*DataFrame)(nil)
	_ starlark.HasSetField = (*DataFrame)(nil)
	_ starlark.HasSetKey   = (*DataFrame)(nil)
	_ starlark.HasBinary   = (*DataFrame)(nil)
)

func readCsv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, fmt.Errorf("dataframe.read_csv is disabled, use dataframe.parse_csv(string) instead to parse csv text that was already downloaded using the http package. In the future, dataframe.read_csv(url) will be restored")
}

func parseCsv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var content starlark.Value

	if err := starlark.UnpackArgs("parse_csv", args, kwargs,
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
		if columns == nil {
			columns = inData.columns
		}
		if index == nil {
			index = inData.index
		}

	case *starlark.Dict:
		var keys []string
		body, keys, err = constructBodyFromStarlarkDict(inData)
		if err != nil {
			return nil, err
		}

		if columns != nil {
			body = constructBodyFromReindexedColumns(body, keys, columns)
		} else {
			columns = NewIndex(keys, "")
		}

	case *starlark.List:
		var colNames []string
		body, colNames, err = constructBodyFromStarlarkList(inData)
		if err != nil {
			return nil, err
		}
		if colNames != nil {
			columns = NewIndex(colNames, "")
		}

	default:
		return nil, fmt.Errorf("Not implemented, constructing DataFrame using %s", reflect.TypeOf(data))
	}

	// Check that the index and columns, if present, match the body size
	numCols, numRows := sizeOfBody(body)
	if index.Len() > 0 && index.Len() != numRows {
		// TODO(dustmop): Add test
		return nil, fmt.Errorf("size of index does not match body size")
	}
	if columns.Len() > 0 && columns.Len() != numCols {
		// TODO(dustmop): Add test
		return nil, fmt.Errorf("number of columns %d, does not match body size %d", columns.Len(), numCols)
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
func constructBodyFromStarlarkList(data *starlark.List) ([]Series, []string, error) {
	numRows := data.Len()
	numCols := -1
	var builder *tableBuilder
	// Iterate the input data row-size
	for y := 0; y < data.Len(); y++ {
		elem := data.Index(y)
		nr := toNamedRowOrNil(elem)
		if nr != nil {
			if builder == nil {
				builder = newTableBuilder(0, numRows)
			}
			builder.pushNamedRow(nr)
			continue
		}
		row := toInterfaceSliceOrNil(elem)
		if row == nil {
			return nil, nil, fmt.Errorf("invalid value for body: %v", elem)
		}
		// Validate that the size of each row is the same
		// TODO(dustmop): Add test
		if numCols == -1 {
			numCols = len(row)
		} else if numCols != len(row) {
			return nil, nil, fmt.Errorf("rows need to be the same length")
		}
		if builder == nil {
			builder = newTableBuilder(numCols, numRows)
		}
		builder.pushRow(row)
	}
	newBody, err := builder.body()
	if err != nil {
		return nil, nil, err
	}
	return newBody, builder.colNames(), nil
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

// dictionaries are re-indexed by column names to create the body
func constructBodyFromReindexedColumns(orig []Series, names []string, columns *Index) []Series {
	_, numRows := sizeOfBody(orig)
	newBody := make([]Series, len(columns.texts))
	for i, col := range columns.texts {
		pos := findKeyPos(col, names)
		if pos == -1 {
			newBody[i] = newTypedSliceBuilderNaNFilled(numRows).toSeries(nil, "")
		} else {
			newBody[i] = orig[pos]
		}
	}
	return newBody
}

// get the size of the body, width and height
func sizeOfBody(body []Series) (int, int) {
	if len(body) == 0 {
		return 0, 0
	}
	return len(body), body[0].Len()
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
	return df.body[0].Len()
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
	// Column names can be accessed as attributes
	v, found, err := df.Get(starlark.String(name))
	if found {
		return v, err
	}
	// Find non-method attribute
	attrImpl, found := dataframeAttributes[name]
	if found {
		return attrImpl(df)
	}
	// Find method
	return builtinAttr(df, name, dataframeMethods)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	// get the non-method attributes
	attrNames := make([]string, 0, len(dataframeAttributes))
	for name := range dataframeAttributes {
		attrNames = append(attrNames, name)
	}
	sort.Strings(attrNames)
	// append the methods
	return append(attrNames, builtinAttrNames(dataframeMethods)...)
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

	// If dataframe has no columns yet, create an empty index
	if df.columns == nil {
		df.columns = NewIndex([]string{}, "")
	}

	// Figure out if a column already exists with the given name
	columnIndex := findKeyPos(name, df.columns.texts)

	// Either prepend the new column, or keep the names the same
	var newNames []string
	if columnIndex == -1 {
		newNames = append(df.columns.texts, name)
	} else {
		newNames = df.columns.texts
	}

	// Assignment of a scalar (int, bool, float, string) to the column
	if scalar, ok := toScalarMaybe(val); ok {
		var newBody []Series
		newCol := newSeriesFromRepeatScalar(scalar, max(1, df.NumRows()))
		if columnIndex == -1 {
			// New columns are added to the right side of the dataframe
			newBody = append(df.body, *newCol)
		} else {
			newBody = df.body
			newBody[columnIndex] = *newCol
		}
		df.columns = NewIndex(newNames, "")
		df.body = newBody
		return nil
	}

	// Convert list to a series
	if list, ok := val.(*starlark.List); ok {
		var err error
		val, err = newSeriesFromList(*list)
		if err != nil {
			return err
		}
	}

	// Assignment of a Series to the column
	series, ok := val.(*Series)
	if !ok {
		return fmt.Errorf("SetKey: val must be int, string, list, or Series")
	}
	if df.NumRows() > 0 && (df.NumRows() != series.Len()) {
		return fmt.Errorf("SetKey: val len must match number of rows")
	}

	var newBody []Series
	if columnIndex == -1 {
		newBody = append(df.body, *series)
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
	if key, ok := toStrMaybe(keyVal); ok {
		val, err := df.accessDataFrameByString(key)
		if err != nil {
			return val, false, err
		}
		return val, true, nil
	}

	if _, ok := keyVal.(starlark.Bool); ok {
		return nil, false, fmt.Errorf("cannot call DataFrame.Get with bool. If you are trying `df[df[column] == val], instead use `df[df[column].equals(val)]`")
	}

	if ser, ok := keyVal.(*Series); ok {
		val, err := df.accessDataFrameBySeries(ser)
		if err != nil {
			return val, false, err
		}
		return val, true, nil
	}

	if list, ok := keyVal.(*starlark.List); ok {
		val, err := df.accessDataFrameByList(list)
		if err != nil {
			return val, false, err
		}
		return val, true, nil
	}

	return nil, false, fmt.Errorf("DataFrame.Get given %v", keyVal)
}

func (df *DataFrame) accessDataFrameByString(key string) (starlark.Value, error) {
	if df.columns == nil {
		return starlark.None, fmt.Errorf("DataFrame.Get: key not found %q", key)
	}

	// Find the column being retrieved, fail if not found
	keyPos := findKeyPos(key, df.columns.texts)
	if keyPos == -1 {
		return starlark.None, fmt.Errorf("DataFrame.Get: key not found %q", key)
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
	}, nil
}

// accessing a dataframe using a series of bools (example: df[Series([True, False])])
// will return a new dataframe that only contains the rows which correspond to
// True booleans
func (df *DataFrame) accessDataFrameBySeries(series *Series) (starlark.Value, error) {
	if series.Len() != df.NumRows() {
		return starlark.None, fmt.Errorf("Item wrong length %d instead of %d", series.Len(), df.NumRows())
	}
	builder := newTableBuilder(df.NumCols(), df.NumRows())
	indexVals := make([]string, 0, df.NumRows())
	line := 0
	for rowIter := newRowIter(df); !rowIter.Done(); rowIter.Next() {
		if line >= series.Len() {
			break
		}
		elem := series.Index(line)
		b, ok := elem.(starlark.Bool)
		if !ok {
			return starlark.None, fmt.Errorf("DataFrame.Get(Series) must be a Series of bools, got %d: %v of %T", line, elem, elem)
		}
		if b {
			builder.pushRow(rowIter.GetRow().data)
			indexVals = append(indexVals, strconv.Itoa(line))
		}
		line++
	}
	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: df.columns,
		index:   NewIndex(indexVals, ""),
		body:    body,
	}, nil
}

// accessing a dataframe using a list of bools (example: df[[True, False, True]])
// will return a new dataframe that only contains the rows which correspond to
// True booleans
func (df *DataFrame) accessDataFrameByList(list *starlark.List) (starlark.Value, error) {
	if list.Len() != df.NumRows() {
		return starlark.None, fmt.Errorf("Item wrong length %d instead of %d", list.Len(), df.NumRows())
	}
	bs := make([]bool, list.Len())
	for i := 0; i < list.Len(); i++ {
		val := list.Index(i)
		b, ok := val.(starlark.Bool)
		if !ok {
			return starlark.None, fmt.Errorf("DataFrame.Get(list) must be a list of bools, got %d: %v of %T", i, val, val)
		}
		bs[i] = bool(b)
	}
	return df.accessDataFrameBySeries(newSeriesFromBools(bs, nil, ""))
}

// At2d returns the cell as position 'i,j' as a go native type
func (df *DataFrame) At2d(i, j int) (interface{}, error) {
	if j >= len(df.body) {
		return nil, fmt.Errorf("index (%d,%d) out of range: %d >= %d (num cols)", i, j, j, len(df.body))
	}
	series := df.body[j]
	if i >= series.Len() {
		return nil, fmt.Errorf("index (%d,%d) out of range: %d >= %d (num rows)", i, j, i, series.Len())
	}
	cell := series.At(i)
	return cell, nil
}

// SetAt2d assigns a go native type to the cell at position 'i,j'
func (df *DataFrame) SetAt2d(i, j int, any interface{}) error {
	if j >= len(df.body) {
		return fmt.Errorf("index (%d,%d) out of range: %d >= %d (num cols)", i, j, j, len(df.body))
	}
	series := df.body[j]
	if i >= series.Len() {
		return fmt.Errorf("index (%d,%d) out of range: %d >= %d (num rows)", i, j, i, series.Len())
	}
	return series.SetAt(i, any)
}

// Binary performs binary operations (like addition) on the DataFrame
func (df *DataFrame) Binary(op syntax.Token, y starlark.Value, side starlark.Side) (starlark.Value, error) {
	// Currently only handle addition, where this DataFrame is the left-hand-side
	if op != syntax.PLUS {
		return nil, nil
	}
	if side {
		return nil, fmt.Errorf("TODO(dustmop): implement DataFrame as rhs of binary +")
	}

	// The right-hand-side is either a DataFrame, or can be used to construct one
	other, ok := y.(*DataFrame)
	if !ok {
		var err error
		other, err = NewDataFrame(y, nil, nil)
		if err != nil {
			return starlark.None, err
		}
	}

	return addTwoDataframes(df, other, df.columns)
}

func addTwoDataframes(left, right *DataFrame, columns *Index) (starlark.Value, error) {
	// Currently must have matching number of columns
	if left.NumCols() != right.NumCols() {
		return nil, fmt.Errorf("TODO(dustmop): handle binary + for different number of columns")
	}

	numCols := left.NumCols()
	numRows := left.NumRows() + right.NumRows()

	// Build the result by concating the rows
	builder := newTableBuilder(numCols, numRows)
	for rowIter := newRowIter(left); !rowIter.Done(); rowIter.Next() {
		builder.pushRow(rowIter.GetRow().data)
	}
	for rowIter := newRowIter(right); !rowIter.Done(); rowIter.Next() {
		builder.pushRow(rowIter.GetRow().data)
	}

	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	return &DataFrame{
		columns: columns,
		body:    body,
	}, nil
}

// ColumnNamesTypes returns the column names and types if they exist
func (df *DataFrame) ColumnNamesTypes() ([]string, []string) {
	if df.columns == nil {
		return nil, nil
	}

	dtypes := make([]string, len(df.body))
	for i, series := range df.body {
		dtypes[i] = series.dtype
	}

	return df.columns.texts, dtypes
}

// at returns an atIndexer which can retrieve or set individual cells
func dataframeAttrAt(self *DataFrame) (starlark.Value, error) {
	return NewAtIndexer(self), nil
}

// columns returns the columns of the dataframe as an index
func dataframeAttrColumns(self *DataFrame) (starlark.Value, error) {
	if self.columns == nil {
		return NewRangeIndex(self.NumCols()), nil
	}
	return self.columns, nil
}

// index returns the index of the dataframe, or a rangeIndex if none exists
func dataframeAttrIndex(self *DataFrame) (starlark.Value, error) {
	if self.index == nil {
		return NewRangeIndex(self.NumRows()), nil
	}
	return self.index, nil
}

// shape returns a tuple with the rows and columns in the dataframe
func dataframeAttrShape(self *DataFrame) (starlark.Value, error) {
	rows := starlark.MakeInt(self.NumRows())
	cols := starlark.MakeInt(self.NumCols())
	return starlark.Tuple{rows, cols}, nil
}

// apply method iterates the rows of a DataFrame, calls the given function for
// each row, creating a single Series of the results
func dataframeApply(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		funcVal, axisVal starlark.Value
		self             = b.Receiver().(*DataFrame)
	)

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

	builder := newTypedSliceBuilder(self.NumRows())
	for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
		r := rowIter.GetRow()
		arguments := r.toTuple()
		res, err := starlark.Call(thread, funcObj, arguments, nil)
		if err != nil {
			return nil, err
		}
		// TODO(dustmop): This won't handle complex types.
		obj, ok := toScalarMaybe(res)
		if !ok {
			return nil, fmt.Errorf("could not convert: %v", res)
		}
		builder.push(obj)
	}
	if err := builder.error(); err != nil {
		return nil, err
	}
	s := builder.toSeries(nil, "")
	return &s, nil
}

// head method returns a copy of the DataFrame but only with the first n rows
func dataframeHead(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var nVal starlark.Value
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("head", args, kwargs,
		"n?", &nVal,
	); err != nil {
		return nil, err
	}

	numRows, ok := toIntMaybe(nVal)
	if !ok {
		// n defaults to 5 if not given
		numRows = 5
	}

	builder := newTableBuilder(self.NumCols(), 0)
	for rowIter := newRowIter(self); !rowIter.Done() && numRows > 0; rowIter.Next() {
		r := rowIter.GetRow()
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
	self := b.Receiver().(*DataFrame)

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
	groupBy, ok := toStrMaybe(first)
	if !ok {
		return nil, fmt.Errorf("by[0] should be a string")
	}

	result := map[string][]*rowTuple{}
	keyPos := findKeyPos(groupBy, self.columns.texts)
	if keyPos == -1 {
		return starlark.None, nil
	}

	for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
		r := rowIter.GetRow()
		groupValue := rowIter.GetStr(keyPos)
		result[groupValue] = append(result[groupValue], r)
	}

	return &GroupByResult{label: groupBy, columns: self.columns, grouping: result}, nil
}

// drop method returns a copy of a DataFrame with rows or columns dropped
func dataframeDrop(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		labelsVal  starlark.Value
		axisVal    starlark.Value
		indexVal   starlark.Value
		columnsVal starlark.Value
		self       = b.Receiver().(*DataFrame)
	)

	if err := starlark.UnpackArgs("drop", args, kwargs,
		"labels?", &labelsVal,
		"axis?", &axisVal,
		"index?", &indexVal,
		"columns?", &columnsVal,
	); err != nil {
		return nil, err
	}

	labels := toStrSliceOrNil(labelsVal)
	axis, ok := toIntMaybe(axisVal)
	if !ok {
		axis = -1
	}
	columns := toStrSliceOrNil(columnsVal)
	index := toIntSliceOrNil(indexVal)

	// Validate axis value, must be 1, or -1 is when it is not given
	if axis == 1 {
		columns = labels
	} else if axis != -1 {
		return nil, fmt.Errorf("axis must equal 1 for dropping columns")
	}

	// Validate index and column parameters
	if index == nil && columns == nil {
		return nil, fmt.Errorf("drop requires either an index or columns")
	}
	if index != nil && columns != nil {
		return nil, fmt.Errorf("drop with both index and column is not supported")
	}
	if index != nil && len(index) != 1 {
		return nil, fmt.Errorf("dropping from index only supports dropping 1 row")
	}
	if columns != nil && len(columns) != 1 {
		return nil, fmt.Errorf("dropping from columns only supports dropping 1 column")
	}

	if columns != nil {
		// Drop columns
		colIndex := findKeyPos(columns[0], self.columns.texts)
		newColumns := removeElemFromStringList(self.columns.texts, colIndex)

		builder := newTableBuilder(self.NumCols()-1, self.NumRows())
		for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
			newRow := removeElemFromInterfaceList(rowIter.GetRow().data, colIndex)
			builder.pushRow(newRow)
		}
		// Finish building the body, return any errors
		body, err := builder.body()
		if err != nil {
			return nil, err
		}
		// Return copy of the dataframe
		return &DataFrame{
			columns: NewIndex(newColumns, ""),
			body:    body,
		}, nil
	}

	// Drop rows using index
	builder := newTableBuilder(self.NumCols(), self.NumRows()-1)
	line := 0
	for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
		if line == index[0] {
			line++
			continue
		}
		builder.pushRow(rowIter.GetRow().data)
		line++
	}
	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}
	// TODO(dustomp): Support indexes with names, remove from them
	// Return copy of the dataframe
	return &DataFrame{
		columns: self.columns,
		body:    body,
	}, nil
}

// drop_duplicates method returns a copy of a DataFrame without duplicates
func dataframeDropDuplicates(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		subset starlark.Value
		self   = b.Receiver().(*DataFrame)
	)

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
	for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
		matchOn := rowIter.Marshal(subsetPos)
		if seen[matchOn] {
			continue
		}
		seen[matchOn] = true
		builder.pushRow(rowIter.GetRow().data)
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
		self                        = b.Receiver().(*DataFrame)
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
		for rowIter := newRowIter(self); !rowIter.Done(); rowIter.Next() {
			key := rowIter.GetStr(leftKey)
			n, has := seen[key]
			if has {
				idxs[n] = append(idxs[n], rowIter.Index())
			} else {
				n = len(idxs)
				idxs = append(idxs, []int{rowIter.Index()})
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

func dataframeSortValues(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		byList, ascendingVal starlark.Value
		self                 = b.Receiver().(*DataFrame)
	)

	if err := starlark.UnpackArgs("sort_values", args, kwargs,
		"by", &byList,
		"ascending?", &ascendingVal,
	); err != nil {
		return nil, err
	}

	// Get values of the column to sort by
	byStrs := toStrSliceOrNil(byList)
	if byStrs == nil {
		return nil, fmt.Errorf("invalid `by` value")
	}
	sortPos := findKeyPos(byStrs[0], self.columns.texts)
	values := self.body[sortPos].stringValues()

	// Make an order list, indexes that refer to the sorted order
	order := make([]int, self.NumRows())
	for i := 0; i < self.NumRows(); i++ {
		order[i] = i
	}
	// `ascending` parameter defaults to true if not explicitly set
	if ascending, ok := ascendingVal.(starlark.Bool); ok && !bool(ascending) {
		// descending order
		sort.Slice(order, func(i, j int) bool {
			return values[order[i]] > values[order[j]]
		})
	} else {
		// ascending order
		sort.Slice(order, func(i, j int) bool {
			return values[order[i]] < values[order[j]]
		})
	}

	// Create the index from the sorted values
	orderStr := make([]string, self.NumRows())
	for i := 0; i < self.NumRows(); i++ {
		if self.index == nil {
			orderStr[i] = strconv.Itoa(order[i])
		} else {
			orderStr[i] = self.index.texts[order[i]]
		}
	}

	// Build the new body using this new order
	builder := newTableBuilder(self.NumCols(), self.NumRows())
	for rowIter := newRowIterWithOrder(self, order); !rowIter.Done(); rowIter.Next() {
		builder.pushRow(rowIter.GetRow().data)
	}
	// Finish building the body, return any errors
	body, err := builder.body()
	if err != nil {
		return nil, err
	}

	return &DataFrame{
		columns: self.columns,
		index:   NewIndex(orderStr, ""),
		body:    body,
	}, nil
}

// reset_index method turns the DataFrame index into a new column
func dataframeResetIndex(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	self := b.Receiver().(*DataFrame)

	if err := starlark.UnpackArgs("reset_index", args, kwargs); err != nil {
		return nil, err
	}

	if self.index == nil {
		return self, nil
	}

	newColumns := append([]string{"index"}, self.columns.texts...)
	newBody := make([]Series, 0, self.NumCols())

	objs := convertStringsToObjects(self.index.texts)
	newBody = append(newBody, Series{which: typeObj, valObjs: objs})
	newBody = append(newBody, self.body...)

	return &DataFrame{
		columns: NewIndex(newColumns, ""),
		body:    newBody,
	}, nil
}

// append adds a new row to the body
func dataframeAppend(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		otherVal starlark.Value
		self     = b.Receiver().(*DataFrame)
	)

	if err := starlark.UnpackArgs("append", args, kwargs,
		"other", &otherVal,
	); err != nil {
		return nil, err
	}

	// TODO(dustmop): append another DataFrame should work

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

func removeElemFromInterfaceList(ls []interface{}, i int) []interface{} {
	if i == -1 {
		return ls
	}
	a := make([]interface{}, len(ls))
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
