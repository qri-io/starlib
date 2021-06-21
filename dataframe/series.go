package dataframe

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

// Series represents a sequence of values, either ints, floats, or objects. This is
// the underlying data structure that is used to create DataFrames. A single column
// of a DataFrame is a Series.
type Series struct {
	frozen bool
	// which determines which of the slice of values holds meaningful data
	which     int
	valInts   []int
	valFloats []float64
	valObjs   []string
	index     *Index
	// dtype is the user-provided and printable data type that the series contains.
	// This will usually match `which`, but not necessarily
	// TODO: Do more research to determine how python pandas treats this value, and
	// when if ever it differs from the true type of data
	dtype string
	name  string
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*Series)(nil)
	_ starlark.Mapping  = (*Series)(nil)
	_ starlark.HasAttrs = (*Series)(nil)
)

var seriesMethods = map[string]*starlark.Builtin{
	"get": starlark.NewBuiltin("get", seriesGet),
}

// Freeze prevents the series from being mutated
func (s *Series) Freeze() {
	s.frozen = true
}

// Hash cannot be used with Series
func (s *Series) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", s.Type())
}

// String returns the Series as a string in a readable, tabular form
func (s *Series) String() string {
	return s.stringify()
}

// Truth converts the series into a bool
func (s *Series) Truth() starlark.Bool {
	// NOTE: In python, calling bool(Series) raises this exception: "ValueError: The truth
	// value of a Series is ambiguous. Use a.empty, a.bool(), a.item(), a.any() or a.all()."
	// Since starlark does not have exceptions, just always return true.
	return true
}

// Type returns the type as a string
func (s *Series) Type() string {
	return fmt.Sprintf("%s.Series", Name)
}

// Attr gets a value for a string attribute
func (s *Series) Attr(name string) (starlark.Value, error) {
	return builtinAttr(s, name, seriesMethods)
}

// AttrNames lists available attributes
func (s *Series) AttrNames() []string {
	return builtinAttrNames(seriesMethods)
}

// Get retrieves a single cell from the Series
func (s *Series) Get(keyVal starlark.Value) (value starlark.Value, found bool, err error) {
	if name, ok := toStrMaybe(keyVal); ok {
		pos := findKeyPos(name, s.index.texts)
		if pos == -1 {
			return starlark.None, false, fmt.Errorf("not found: %q", name)
		}
		val, err := convertToStarlark(s.values()[pos])
		if err != nil {
			return starlark.None, false, err
		}
		return val, true, nil
	}
	if index, ok := toIntMaybe(keyVal); ok {
		val, err := convertToStarlark(s.values()[index])
		if err != nil {
			return starlark.None, false, err
		}
		return val, true, nil
	}
	return starlark.None, false, fmt.Errorf("not found: %q", keyVal)
}

func (s *Series) stringify() string {
	// Calculate how wide the index column needs to be
	indexWidth := 0
	if s.index.len() == 0 {
		indexWidth = len(fmt.Sprintf("%d", s.len()-1))
	} else {
		for _, elem := range s.index.texts {
			w := len(elem)
			if w > indexWidth {
				indexWidth = w
			}
		}
	}

	// Calculate how wide the data column needs to be
	colWidth := 0
	for _, elem := range s.stringValues() {
		elem = coerceToDatatype(elem, s.dtype)
		w := len(elem)
		if w > colWidth {
			colWidth = w
		}
	}

	// Final line shows the (optional) name and dtype
	epilogue := fmt.Sprintf("dtype: %s", s.dtype)
	if s.dtype == "" {
		epilogue = fmt.Sprintf("dtype: int64")
	}
	if s.name != "" {
		epilogue = fmt.Sprintf("Name: %s, %s", s.name, epilogue)
	}

	// Determine how to format each line, based upon the column width
	padding := "    "
	var tmpl string
	if s.index.len() == 0 {
		// Result looks like '%-2d    %6s'
		tmpl = fmt.Sprintf("%%-%dd%s%%%ds", indexWidth, padding, colWidth)
	} else {
		// Result looks like '%-4s    %6s'
		tmpl = fmt.Sprintf("%%-%ds%s%%%ds", indexWidth, padding, colWidth)
	}

	// Space for the lines of rendered output, the body, plus optional index.name and types
	render := make([]string, 0, s.len()+2)

	// If the index has a name, it appears on the first line
	if s.index != nil && s.index.name != "" {
		render = append(render, s.index.name)
	}

	// Render each value in the series
	for i, elem := range s.stringValues() {
		elem = coerceToDatatype(elem, s.dtype)
		line := ""
		if s.index.len() == 0 {
			line = fmt.Sprintf(tmpl, i, elem)
		} else {
			line = fmt.Sprintf(tmpl, s.index.texts[i], elem)
		}
		render = append(render, line)
	}

	// Combine the lines together
	render = append(render, epilogue)
	return strings.Join(render, "\n")
}

// values returns a slice of some go native type
func (s *Series) values() []interface{} {
	if s.which == typeInt {
		result := make([]interface{}, len(s.valInts))
		for i, elem := range s.valInts {
			result[i] = elem
		}
		return result
	} else if s.which == typeFloat {
		result := make([]interface{}, len(s.valFloats))
		for i, elem := range s.valFloats {
			result[i] = elem
		}
		return result
	}
	result := make([]interface{}, len(s.valObjs))
	for i, elem := range s.valObjs {
		result[i] = elem
	}
	return result
}

// stringValues returns a slice of the stringified values, fit for printing
func (s *Series) stringValues() []string {
	if s.which == typeInt {
		result := make([]string, len(s.valInts))
		for i, elem := range s.valInts {
			result[i] = strconv.Itoa(elem)
		}
		return result
	} else if s.which == typeFloat {
		result := make([]string, len(s.valFloats))
		for i, elem := range s.valFloats {
			result[i] = stringifyFloat(elem)
		}
		return result
	}
	return s.valObjs
}

// len returns the number of values
func (s *Series) len() int {
	if s.which == typeInt {
		return len(s.valInts)
	} else if s.which == typeFloat {
		return len(s.valFloats)
	}
	return len(s.valObjs)
}

// strAt returns the cell at position 'i', as a string fit for printing
func (s *Series) strAt(i int) string {
	if s.which == typeInt {
		return strconv.Itoa(s.valInts[i])
	} else if s.which == typeFloat {
		return stringifyFloat(s.valFloats[i])
	}
	return s.valObjs[i]
}

func builtinAttr(recv starlark.Value, name string, methods map[string]*starlark.Builtin) (starlark.Value, error) {
	b := methods[name]
	if b == nil {
		return nil, nil // no such method
	}
	return b.BindReceiver(recv), nil
}

func builtinAttrNames(methods map[string]*starlark.Builtin) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func seriesGet(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &key); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)
	ret, _, err := self.Get(key)
	return ret, err
}

func newSeries(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		dataVal  starlark.Value
		indexVal starlark.Value
		dtypeVal starlark.Value
		nameVal  starlark.Value
	)
	if err := starlark.UnpackArgs("Series", args, kwargs,
		"data", &dataVal,
		"index?", &indexVal,
		"dtype?", &dtypeVal,
		"name?", &nameVal,
	); err != nil {
		return nil, err
	}

	name := toStrOrEmpty(nameVal)
	dtype := toStrOrEmpty(dtypeVal)
	index, _ := toIndexMaybe(indexVal)

	// Series built from a scalar value
	if scalarNum, ok := toIntMaybe(dataVal); ok {
		if dtype == "float64" {
			return newSeriesFromFloats([]float64{float64(scalarNum)}, index, name), nil
		} else if dtype == "object" {
			return newSeriesFromStrings([]string{strconv.Itoa(scalarNum)}, index, name), nil
		}
		return newSeriesFromInts([]int{scalarNum}, index, name), nil
	}
	if scalarFloat, ok := toFloatMaybe(dataVal); ok {
		return newSeriesFromFloats([]float64{scalarFloat}, index, name), nil
	}
	if scalarStr, ok := toStrMaybe(dataVal); ok {
		return newSeriesFromStrings([]string{scalarStr}, index, name), nil
	}

	switch inData := dataVal.(type) {
	case *starlark.List:
		builder := newTypedArrayBuilder(inData.Len())
		builder.setType(dtype)

		for k := 0; k < inData.Len(); k++ {
			elemVal := inData.Index(k)
			if scalar, ok := toScalarMaybe(elemVal); ok {
				builder.push(scalar)
				continue
			}
			// TODO: return an error for this invalid element, add a test
		}
		if err := builder.error(); err != nil {
			return starlark.None, err
		}
		series := builder.toSeries(index, name)
		return &series, nil
	case *starlark.Dict:
		builder := newTypedArrayBuilder(inData.Len())
		builder.setType(dtype)

		keys := inData.Keys()
		for i := 0; i < len(keys); i++ {
			keyVal := keys[i]
			key, ok := keyVal.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("dict key must be string")
			}
			val, _, _ := inData.Get(keyVal)
			if scalar, ok := toScalarMaybe(val); ok {
				builder.pushKeyVal(string(key), scalar)
				continue
			}
			// TODO: return an error for this invalid element, add a test
		}
		if err := builder.error(); err != nil {
			return starlark.None, err
		}
		// TODO: If index is provided, reindex the series.
		index := NewIndex(builder.keys(), "")
		series := builder.toSeries(index, name)
		return &series, nil
	}

	return starlark.None, fmt.Errorf("`data` type unrecognized: %q of %s", dataVal.String(), dataVal.Type())
}

func newSeriesFromRepeatScalar(val interface{}, size int) *Series {
	switch x := val.(type) {
	case int:
		vals := make([]int, size)
		for i := 0; i < size; i++ {
			vals[i] = x
		}
		return newSeriesFromInts(vals, nil, "")
	case float64:
		vals := make([]float64, size)
		for i := 0; i < size; i++ {
			vals[i] = x
		}
		return newSeriesFromFloats(vals, nil, "")
	case string:
		vals := make([]string, size)
		for i := 0; i < size; i++ {
			vals[i] = x
		}
		return newSeriesFromStrings(vals, nil, "")
	default:
		return nil
	}
}

func newSeriesFromInts(vals []int, index *Index, name string) *Series {
	return &Series{
		dtype:   "int64",
		which:   typeInt,
		valInts: vals,
		index:   index,
		name:    name,
	}
}

func newSeriesFromFloats(vals []float64, index *Index, name string) *Series {
	return &Series{
		dtype:     "float64",
		which:     typeFloat,
		valFloats: vals,
		index:     index,
		name:      name,
	}
}

func newSeriesFromStrings(vals []string, index *Index, name string) *Series {
	return &Series{
		dtype:   "object",
		which:   typeObj,
		valObjs: vals,
		index:   index,
		name:    name,
	}
}

func findKeyPos(needle string, subject []string) int {
	for i, elem := range subject {
		if elem == needle {
			return i
		}
	}
	return -1
}
