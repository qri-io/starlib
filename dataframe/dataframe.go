package dataframe

import (
	"fmt"
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
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrame),
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

	columns := toStrSliceOrNil(columnsVal)
	index, _ := toIndexMaybe(indexVal)

	switch inData := dataVal.(type) {
	case *starlark.Dict:
		newBody := make([]Series, 0)
		keyList := make([]string, 0, inData.Len())
		inKeys := inData.Keys()
		numCols := len(inKeys)
		numRows := -1
		for i := 0; i < len(inKeys); i++ {
			// Collect each key, use them as the default index
			inKey := inKeys[i]
			keyList = append(keyList, toStr(inKey))
			// Get each value, which should be a list of values
			val, _, _ := inData.Get(inKey)
			items := toInterfaceSliceOrNil(val)
			if items == nil {
				return starlark.None, fmt.Errorf("invalid values for column")
			}
			// Validate that the size of each column is the same
			// TODO(dustmop): Add test
			if numRows == -1 {
				numRows = len(items)
			} else if numRows != len(items) {
				return starlark.None, fmt.Errorf("columns need to be the same length")
			}
			// The list of values should be of the same type
			builder := newTypedArrayBuilder(len(items))
			for _, it := range items {
				builder.push(it)
			}
			if err := builder.error(); err != nil {
				return starlark.None, err
			}

			newBody = append(newBody, builder.toSeries(nil, ""))
		}

		if index.len() > 0 && index.len() != numRows {
			// TODO(dustmop): Add test
			return starlark.None, fmt.Errorf("size of index does not match body size")
		}
		if len(columns) > 0 && len(columns) != numCols {
			// TODO(dustmop): Add test
			return starlark.None, fmt.Errorf("number of columns does not match body size")
		}

		// TODO(dustmop): `index` will re-index
		return &DataFrame{
			columns: NewIndex(keyList, ""),
			body:    newBody,
		}, nil
	case *starlark.List:
		numRows := inData.Len()
		numCols := -1
		var builders []*typedSliceBuilder
		// Iterate the input data row-size
		for y := 0; y < inData.Len(); y++ {
			row := toInterfaceSliceOrNil(inData.Index(y))
			if row == nil {
				return starlark.None, fmt.Errorf("invalid values for row")
			}
			// Validate that the size of each row is the same
			// TODO(dustmop): Add test
			if numCols == -1 {
				numCols = len(row)
			} else if numCols != len(row) {
				return starlark.None, fmt.Errorf("rows need to be the same length")
			}
			for x := 0; x < numCols; x++ {
				// Allocate builders once we know how many and how large they are
				if builders == nil {
					builders = make([]*typedSliceBuilder, numCols)
					for i := 0; i < numCols; i++ {
						builders[i] = newTypedArrayBuilder(numRows)
					}
				}
				// Accumlate each cell into the appropriate column builder
				builders[x].push(row[x])
			}
		}
		// Get a series for each column of the body
		newBody := make([]Series, numCols)
		for x := 0; x < numCols; x++ {
			newBody[x] = builders[x].toSeries(nil, "")
		}

		if index.len() > 0 && index.len() != numRows {
			// TODO(dustmop): Add test
			return starlark.None, fmt.Errorf("size of index does not match body size")
		}
		if len(columns) > 0 && len(columns) != numCols {
			// TODO(dustmop): Add test
			return starlark.None, fmt.Errorf("number of columns does not match body size")
		}

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
	return nil, nil
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	return []string{"columns"}
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
		newCol := newSeriesFromRepeatScalar(scalar, max(1, df.numRows()))
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
	// TODO(dustmop): index should be the left-hand-side index, need a test
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

func dtypeFromWhich(which int) string {
	if which == typeInt {
		return "int64"
	} else if which == typeFloat {
		return "float64"
	}
	return "object"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
