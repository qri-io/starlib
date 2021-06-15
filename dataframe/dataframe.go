package dataframe

import (
	"fmt"
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
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrame),
		"Index":     starlark.NewBuiltin("Index", newIndex),
		"Series":    starlark.NewBuiltin("Series", newSeries),
	},
}

func unfinishedError(v starlark.Value, msg string) error {
	return fmt.Errorf("%s %s unfinished implementation: %s", Name, v.Type(), msg)
}

type DataFrame struct {
	frozen  bool
	columns *Index
	index   *Index
	body    []Series
}

// compile-time interface assertions
var (
	_ starlark.Value   = (*DataFrame)(nil)
	_ starlark.Mapping = (*DataFrame)(nil)
	_ starlark.HasAttrs = (*DataFrame)(nil)
	_ starlark.HasSetField = (*DataFrame)(nil)
	_ starlark.HasSetKey = (*DataFrame)(nil)
)

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

	// TODO: Assert that all columns have the same size (height == numRows)
	// TODO: Assert that len(columns) == 0 || len(columns) == len(body)
	// TODO: Assert that len(index) == 0 || len(index) == numRows()

	columns := toStrListOrNil(columnsVal)
	index, _ := toIndexMaybe(indexVal)

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
			// TODO: Generalize this. If the list is integers, collect as an []int
			// TODO: Put another way, ensure that when a DataFrame is created, regardless of
			// what is used to construct it, type check each column. Int columns should make
			// int serieses. Same for booleans, floats, timestamps.
			if valInts, ok := maybeIntList(valList); ok {
				series := Series{which: typeInt, valInts: valInts}
				newBody = append(newBody, series)
				continue
			}
			// TODO: Handle `valList == nil`, don't crash
			series := Series{which: typeObj, valObjs: valList}
			newBody = append(newBody, series)
		}

		// TODO: `index` will re-index the columns
		return &DataFrame{
			columns: NewIndex(keyList, ""),
			body:    newBody,
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
			columns: NewIndex(columns, ""),
			index:   index,
			body:    newBody,
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
	return df.stringify()
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
		return df.columns, nil
	}
	return nil, nil
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	return []string{"columns"}
}

func (df *DataFrame) SetField(name string, val starlark.Value) error {
	if df.frozen {
		return fmt.Errorf("cannot set, dataframe is frozen")
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

func (df *DataFrame) SetKey(nameVal, val starlark.Value) error {
	if df.frozen {
		return fmt.Errorf("cannot set, dataframe is frozen")
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
		df.columns = NewIndex(newNames, "")
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
		df.columns = NewIndex(newNames, "")
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
	df.columns = NewIndex(newNames, "")
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
	keyPos := findKeyPos(key, df.columns.texts)
	if keyPos == -1 {
		return starlark.None, false, fmt.Errorf("not found")
	}

	got := df.body[keyPos]
	// TODO: index should be the left-hand-side index, need a test
	index := NewIndex(nil, "")

	dtype := got.dtype
	if dtype == "" {
		dtype = dtypeFromWhich(got.which)
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
		bodyHeight := len(df.body)
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
	widths := make([]int, len(df.body))
	for i, name := range df.columns.texts {
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

	// Render columns
	header := make([]string, 0, len(df.columns.texts))
	if len(df.columns.texts) > 0 {
		// Render the column names
		for i, name := range df.columns.texts {
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
	for i := 0; i < df.numRows(); i++ {
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

func maybeIntList(vals []string) ([]int, bool) {
	newVals := make([]int, 0, len(vals))
	for _, elem := range vals {
		num, err := strconv.Atoi(elem)
		if err != nil {
			return nil, false
		}
		newVals = append(newVals, num)
	}
	return newVals, true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
