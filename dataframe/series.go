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
	frozen    bool
	which     int
	valInts   []int
	valFloats []float64
	valObjs   []string
	index     []string
	dtype     string
	name      string
}

// A Series contains values of one of these three types. The which field uses these
// constants to determine which slice holds the actual values of the Series.
const (
	typeInt   = 1
	typeFloat = 2
	typeObj   = 3
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
	return 0, fmt.Errorf("'Series' objects are mutable, thus they cannot be hashed")
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
	return "dataframe.Series"
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
		pos := findKeyPos(name, s.index)
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
	if len(s.index) == 0 {
		indexWidth = len(fmt.Sprintf("%d", s.len()-1))
	} else {
		for _, elem := range s.index {
			w := len(elem)
			if w > indexWidth {
				indexWidth = w
			}
		}
	}

	// Calculate how wide the data column needs to be
	colWidth := 0
	for _, elem := range s.stringValues() {
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
	if len(s.index) == 0 {
		// Result looks like '%-2d    %6s'
		tmpl = fmt.Sprintf("%%-%dd%s%%%ds", indexWidth, padding, colWidth)
	} else {
		// Result looks like '%-4s    %6s'
		tmpl = fmt.Sprintf("%%-%ds%s%%%ds", indexWidth, padding, colWidth)
	}

	// Render each value in the series
	render := make([]string, 0, s.len()+1)
	for i, elem := range s.stringValues() {
		line := ""
		if len(s.index) == 0 {
			line = fmt.Sprintf(tmpl, i, elem)
		} else {
			line = fmt.Sprintf(tmpl, s.index[i], elem)
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
			result[i] = fmt.Sprintf("%1.1f", elem)
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
		return fmt.Sprintf("%1.1f", s.valFloats[i])
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
	index := toStrListOrNil(indexVal)

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

	which := 0
	if dtype == "int64" {
		which = typeInt
	} else if dtype == "float64" {
		which = typeFloat
	}

	dataList, ok := dataVal.(*starlark.List)
	if ok {
		// Series built from a list
		valInts := make([]int, 0, dataList.Len())
		valFloats := make([]float64, 0, dataList.Len())
		valObjs := make([]string, 0, dataList.Len())

		for k := 0; k < dataList.Len(); k++ {
			elemVal := dataList.Index(k)
			// If an int, convert to float if that's what this Series is typed as
			if num, ok := toIntMaybe(elemVal); ok {
				if which == 0 || which == typeInt {
					which = typeInt
					valInts = append(valInts, num)
					continue
				} else if which == typeFloat {
					valFloats = append(valFloats, float64(num))
					continue
				}
				which = typeObj
			}
			// If a float, convert an existing Series of ints to floats
			if f, ok := toFloatMaybe(elemVal); ok {
				if which == 0 || which == typeFloat {
					which = typeFloat
					valFloats = append(valFloats, f)
					continue
				} else if which == typeInt {
					which = typeFloat
					valFloats = append(convertIntsToFloats(valInts), f)
					continue
				}
				which = typeObj
			}
			// Otherwise, or if the type is object, convert everything to strings
			if which == typeInt {
				valObjs = convertIntsToStrings(valInts)
			} else if which == typeFloat {
				valObjs = convertFloatsToStrings(valFloats)
			}
			which = typeObj
			valObjs = append(valObjs, toStr(elemVal))
		}

		// If no dtype was provided, derive it from the values in the Series
		if dtype == "" {
			dtype = dtypeFromWhich(which)
		}
		return &Series{
			dtype:     dtype,
			which:     which,
			valInts:   valInts,
			valFloats: valFloats,
			valObjs:   valObjs,
			index:     index,
			name:      name,
		}, nil
	}

	dataDict, ok := dataVal.(*starlark.Dict)
	if ok {
		// Series built from a dict
		valObjs := make([]string, 0, dataDict.Len())
		index = make([]string, 0, dataDict.Len())

		keys := dataDict.Keys()
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			val, _, _ := dataDict.Get(key)

			// TODO: Building from a dict coerces everything to string, should retain
			// types as is done for lists
			index = append(index, toStr(key))
			valObjs = append(valObjs, toStr(val))
		}
		return &Series{
			dtype:   dtype,
			which:   typeObj,
			valObjs: valObjs,
			index:   index,
			name:    name,
		}, nil
	}

	return starlark.None, fmt.Errorf("`data` type unrecognized: %q of %s", dataVal.String(), dataVal.Type())
}

func newSeriesFromInts(vals []int, index []string, name string) *Series {
	return &Series{
		dtype:   "int64",
		which:   typeInt,
		valInts: vals,
		index:   index,
		name:    name,
	}
}

func newSeriesFromFloats(vals []float64, index []string, name string) *Series {
	return &Series{
		dtype:     "float64",
		which:     typeFloat,
		valFloats: vals,
		index:     index,
		name:      name,
	}
}

func newSeriesFromStrings(vals, index []string, name string) *Series {
	return &Series{
		dtype:   "object",
		which:   typeObj,
		valObjs: vals,
		index:   index,
		name:    name,
	}
}

func dtypeFromWhich(which int) string {
	if which == typeInt {
		return "int64"
	} else if which == typeFloat {
		return "float64"
	}
	return "object"
}

func findKeyPos(needle string, subject []string) int {
	for i, elem := range subject {
		if elem == needle {
			return i
		}
	}
	return -1
}
