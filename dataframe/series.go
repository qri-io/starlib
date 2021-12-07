package dataframe

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
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
	valObjs   []interface{}
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
	_ starlark.Value     = (*Series)(nil)
	_ starlark.Mapping   = (*Series)(nil)
	_ starlark.HasAttrs  = (*Series)(nil)
	_ starlark.Indexable = (*Series)(nil)
	_ starlark.Sequence  = (*Series)(nil)
	_ starlark.HasUnary  = (*Series)(nil)
)

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
	if name == "dtype" {
		return starlark.String(s.dtype), nil
	} else if name == "index" {
		return s.index, nil
	} else if name == "str" {
		return &stringMethods{subject: s}, nil
	}
	// Find non-method attribute
	attrImpl, found := seriesAttributes[name]
	if found {
		return attrImpl(s)
	}
	return builtinAttr(s, name, seriesMethods)
}

// AttrNames lists available attributes
func (s *Series) AttrNames() []string {
	// TODO: Use seriesAttributes
	attributeNames := []string{"dtype", "index", "str"}
	return append(attributeNames, builtinAttrNames(seriesMethods)...)
}

// Get retrieves a single cell from the Series
func (s *Series) Get(keyVal starlark.Value) (value starlark.Value, found bool, err error) {
	if name, ok := toStrMaybe(keyVal); ok {
		pos := findKeyPos(name, s.index.texts)
		if pos == -1 {
			return starlark.None, false, fmt.Errorf("Series.Get: not found: %q", name)
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
	// TODO(dustmop): Also support series.get(list)
	if keyList, ok := keyVal.(*Series); ok {
		if keyList.dtype != "bool" {
			return starlark.None, false, fmt.Errorf("Series.Get[series] only supported for dtype bool")
		}
		vals := s.stringValues()
		newIdx := make([]string, 0, len(vals))
		newVals := make([]interface{}, 0, len(vals))
		for i, key := range keyList.values() {
			// NOTE: The dtype is checked above, to validate it is "bool"
			if key == 0 {
				continue
			}
			newIdx = append(newIdx, fmt.Sprintf("%d", i))
			newVals = append(newVals, vals[i])
		}
		return newSeriesFromObjects(newVals, NewIndex(newIdx, ""), s.name), true, nil
	}
	return starlark.None, false, fmt.Errorf("Series.Get: not found: %q", keyVal)
}

func seriesEquals(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &key); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)
	return self.selectElemsWhereEqString(string(key))
}

func (s *Series) selectElemsWhereEqString(cmp string) (*Series, error) {
	builder := newTypedSliceBuilder(s.Len())
	builder.setType("bool")

	for k := 0; k < s.Len(); k++ {
		elemVal := s.Index(k)
		if elemVal == nil || elemVal == starlark.None {
			builder.pushNil()
			continue
		}
		if elemStr, ok := toStrMaybe(elemVal); ok {
			if elemStr == cmp {
				builder.push(true)
			} else {
				builder.push(false)
			}
			continue
		} else {
			return nil, fmt.Errorf("expected Series to contain strings, got %v", elemVal)
		}
	}
	if err := builder.error(); err != nil {
		return nil, err
	}
	ans := builder.toSeries(nil, s.name)
	return &ans, nil
}

func seriesNotEquals(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	val, err := seriesEquals(thread, b, args, kwargs)
	if err != nil {
		return val, err
	}

	series := val.(*Series)
	res := make([]bool, series.Len())
	for k := 0; k < series.Len(); k++ {
		elem := series.Index(k)
		b := elem.(starlark.Bool)
		res[k] = !bool(b)
	}

	return newSeriesFromBools(res, series.index, series.name), nil
}

func (s *Series) stringify() string {
	// Calculate how wide the index column needs to be
	indexWidth := 0
	if s.index.Len() == 0 {
		indexWidth = len(fmt.Sprintf("%d", s.Len()-1))
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
		w := len(elem)
		if w > colWidth {
			colWidth = w
		}
	}

	// Final line shows the (optional) name and dtype
	epilogue := fmt.Sprintf("dtype: %s", s.dtype)
	if s.dtype == "" {
		epilogue = "dtype: int64"
	}
	if s.name != "" {
		epilogue = fmt.Sprintf("Name: %s, %s", s.name, epilogue)
	}

	// Determine how to format each line, based upon the column width
	padding := "    "
	var tmpl string
	if s.index.Len() == 0 {
		// Result looks like '%-2d    %6s'
		tmpl = fmt.Sprintf("%%-%dd%s%%%ds", indexWidth, padding, colWidth)
	} else {
		// Result looks like '%-4s    %6s'
		tmpl = fmt.Sprintf("%%-%ds%s%%%ds", indexWidth, padding, colWidth)
	}

	// Space for the lines of rendered output, the body, plus optional index.name and types
	render := make([]string, 0, s.Len()+2)

	// If the index has a name, it appears on the first line
	if s.index != nil && s.index.name != "" {
		render = append(render, s.index.name)
	}

	// Render each value in the series
	for i, elem := range s.stringValues() {
		line := ""
		if s.index.Len() == 0 {
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
		if s.dtype == "bool" {
			for i, elem := range s.valInts {
				if elem == 0 {
					result[i] = "False"
				} else {
					result[i] = "True"
				}
			}
			return result
		}
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
	result := make([]string, len(s.valObjs))
	for i, elem := range s.valObjs {
		if elem == nil {
			result[i] = "None"
		} else if elem == true {
			result[i] = "True"
		} else if elem == false {
			result[i] = "False"
		} else {
			result[i] = fmt.Sprintf("%v", elem)
		}
	}
	return result
}

// Iterate returns an iterator for the series
func (s *Series) Iterate() starlark.Iterator {
	return &seriesIterator{series: s, count: 0}
}

// Len returns the number of values
func (s *Series) Len() int {
	if s.which == typeInt {
		return len(s.valInts)
	} else if s.which == typeFloat {
		return len(s.valFloats)
	}
	return len(s.valObjs)
}

// Index returns the element at index i as a starlark value
func (s *Series) Index(i int) starlark.Value {
	obj, err := convertToStarlark(s.At(i))
	if err != nil {
		return starlark.None
	}
	return obj
}

// StrAt returns the cell at position 'i', as a string fit for printing
func (s *Series) StrAt(i int) string {
	if s.which == typeInt {
		if s.dtype == "bool" {
			if s.valInts[i] == 0 {
				return "False"
			}
			return "True"
		}
		return strconv.Itoa(s.valInts[i])
	} else if s.which == typeFloat {
		return stringifyFloat(s.valFloats[i])
	}
	if s.valObjs[i] == nil {
		return "None"
	} else if s.valObjs[i] == true {
		return "True"
	} else if s.valObjs[i] == false {
		return "False"
	}
	return fmt.Sprintf("%v", s.valObjs[i])
}

// At returns the cell at position 'i' as a go native type
func (s *Series) At(i int) interface{} {
	if s.which == typeInt {
		if s.dtype == "bool" {
			return s.valInts[i] != 0
		}
		return s.valInts[i]
	} else if s.which == typeFloat {
		return s.valFloats[i]
	}
	return s.valObjs[i]
}

// SetAt assigns a go native type to the cell at position 'i'
func (s *Series) SetAt(i int, any interface{}) error {
	switch item := any.(type) {
	case int:
		if s.which == typeInt {
			s.valInts[i] = item
		} else {
			return fmt.Errorf("TODO: implement SetAt(int) conversion")
		}
	case string:
		if s.which == typeObj {
			s.valObjs[i] = item
		} else {
			return fmt.Errorf("TODO: implement SetAt(string) conversion")
		}
	case interface{}:
		if s.which == typeObj {
			s.valObjs[i] = item
		} else {
			return fmt.Errorf("TODO: implement SetAt(interface) conversion")
		}
	}
	return nil
}

// CloneWithStrings returns a clone of the series with contents replaced with the given strings
func (s *Series) CloneWithStrings(txts []string) starlark.Value {
	return &Series{
		dtype:   "object",
		which:   typeObj,
		valObjs: convertStringsToObjects(txts),
		index:   s.index,
		name:    s.name,
	}
}

// Unary implements unary operators, only the tilde (negation) is supported
func (s *Series) Unary(op syntax.Token) (value starlark.Value, err error) {
	if op != syntax.TILDE {
		return starlark.None, fmt.Errorf("only unary ~ is supported")
	}

	result := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		obj := s.At(i)
		if b, ok := obj.(bool); ok {
			if !b {
				result[i] = true
				continue
			}
		}
		result[i] = false
	}

	return newSeriesFromBools(result, s.index, s.name), nil
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

// name returns of the name of the series
func seriesAttrName(self *Series) (starlark.Value, error) {
	return starlark.String(self.name), nil
}

// size returns the number of elements in the series
func seriesAttrSize(self *Series) (starlark.Value, error) {
	return starlark.MakeInt(self.Len()), nil
}

func adaptToSeriesFromDataframe(methodName string) starlarkMethod {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		self := b.Receiver().(*Series)
		outconf, _ := thread.Local("OutputConfig").(*OutputConfig)

		// Convert the series to a DataFrame
		df, err := NewDataFrame(self, nil, self.index, outconf)
		if err != nil {
			return starlark.None, err
		}

		// Lookup the adapted method
		method, err := df.Attr(methodName)
		if err != nil {
			return starlark.None, err
		}

		// Call it, and ensure the result is a DataFrame
		res, err := starlark.Call(thread, method, args, kwargs)
		if err != nil {
			return starlark.None, err
		}
		df, ok := res.(*DataFrame)
		if !ok {
			return starlark.None, fmt.Errorf("expected DataFrame result, got %T", res)
		}

		// Retrieve the first column as a series, and return it
		series, _, err := df.Get(starlark.MakeInt(0))
		if err != nil {
			return starlark.None, err
		}
		return series, nil
	}
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

// astype method converts a Series by coercing its values to the given type
func seriesAsType(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var typeVal starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &typeVal); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)

	typeName, _ := toStrMaybe(typeVal)
	if typeName != "int64" {
		return nil, fmt.Errorf("invalid type, only \"int64\" allowed")
	}

	newVals := make([]int, 0, self.Len())
	for _, val := range self.values() {
		text := fmt.Sprintf("%s", val)

		// Special case: convert datetime to nanoseconds
		if self.dtype == "datetime64[ns]" && self.which == typeObj {
			t, err := time.Parse("2006-01-02 15:04:05", text)
			if err != nil {
				return nil, err
			}
			num := t.UnixNano()
			newVals = append(newVals, int(num))
			continue
		}

		// Default case, parse the value as an integer
		num, err := strconv.Atoi(text)
		if err != nil {
			num = -1
		}
		newVals = append(newVals, num)
	}

	return newSeriesFromInts(newVals, self.index, self.name), nil
}

// notnull method returns a Series of booleans that are true for non-null values
func seriesNotNull(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)

	newVals := make([]int, 0, self.Len())
	for _, val := range self.values() {
		if val == nil {
			newVals = append(newVals, 0)
		} else {
			newVals = append(newVals, 1)
		}
	}

	series := newSeriesFromInts(newVals, self.index, self.name)
	series.dtype = "bool"
	return series, nil
}

// to_frame converts a Series to a DataFrame
func seriesToFrame(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)
	outconf, _ := thread.Local("OutputConfig").(*OutputConfig)
	return NewDataFrame(self, nil, self.index, outconf)
}

// unique method returns a list of the unique elements from the series
func seriesUnique(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)

	have := make(map[string]struct{})

	items := make([]starlark.Value, 0)
	for k := 0; k < self.Len(); k++ {
		key := self.StrAt(k)
		if _, found := have[key]; found {
			continue
		}
		have[key] = struct{}{}
		it, err := convertToStarlark(self.At(k))
		if err != nil {
			return starlark.None, err
		}
		items = append(items, it)
	}

	return starlark.NewList(items), nil
}

// reset_index turns the index into a column
func seriesResetIndex(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 0); err != nil {
		return nil, err
	}
	self := b.Receiver().(*Series)
	outconf, _ := thread.Local("OutputConfig").(*OutputConfig)

	df, err := NewDataFrame(self, []string{"id"}, self.index, outconf)
	if err != nil {
		return starlark.None, err
	}

	method, err := df.Attr("reset_index")
	if err != nil {
		return starlark.None, err
	}

	res, err := starlark.Call(thread, method, args, kwargs)
	if err != nil {
		return starlark.None, err
	}

	return res, nil
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
			return newSeriesFromObjects([]interface{}{scalarNum}, index, name), nil
		}
		return newSeriesFromInts([]int{scalarNum}, index, name), nil
	}
	if scalarFloat, ok := toFloatMaybe(dataVal); ok {
		return newSeriesFromFloats([]float64{scalarFloat}, index, name), nil
	}
	if scalarStr, ok := toStrMaybe(dataVal); ok {
		return newSeriesFromObjects([]interface{}{scalarStr}, index, name), nil
	}

	switch inData := dataVal.(type) {
	case *starlark.List:
		builder := newTypedSliceBuilder(inData.Len())
		builder.setType(dtype)

		for k := 0; k < inData.Len(); k++ {
			elemVal := inData.Index(k)
			if elemVal == nil || elemVal == starlark.None {
				builder.pushNil()
				continue
			}
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
		builder := newTypedSliceBuilder(inData.Len())
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

func newSeriesFromList(list starlark.List) (*Series, error) {
	builder := newTypedSliceBuilder(list.Len())
	for k := 0; k < list.Len(); k++ {
		elemVal := list.Index(k)
		elem := toNativeValue(elemVal)
		builder.push(elem)
	}
	if err := builder.error(); err != nil {
		return nil, err
	}
	s := builder.toSeries(nil, "")
	return &s, nil
}

func newSeriesFromRepeatScalar(val interface{}, size int) *Series {
	if val == nil {
		return newSeriesFromObjects(make([]interface{}, size), nil, "")
	}
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
		vals := make([]interface{}, size)
		for i := 0; i < size; i++ {
			vals[i] = x
		}
		return newSeriesFromObjects(vals, nil, "")
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

func newSeriesFromBools(vals []bool, index *Index, name string) *Series {
	newVals := make([]int, len(vals))
	for i, b := range vals {
		if b {
			newVals[i] = 1
		} else {
			newVals[i] = 0
		}
	}
	series := newSeriesFromInts(newVals, index, name)
	series.dtype = "bool"
	return series
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

func newSeriesFromObjects(vals []interface{}, index *Index, name string) *Series {
	return &Series{
		dtype:   "object",
		which:   typeObj,
		valObjs: vals,
		index:   index,
		name:    name,
	}
}

func newSeriesConstructor(vals []interface{}, index *Index, name string) *Series {
	builder := newTypedSliceBuilder(len(vals))
	for _, v := range vals {
		builder.push(v)
	}
	ans := builder.toSeries(index, name)
	return &ans
}

func findKeyPos(needle string, subject []string) int {
	for i, elem := range subject {
		if elem == needle {
			return i
		}
	}
	return -1
}

type seriesIterator struct {
	count  int
	series *Series
}

// Done does cleanup work when iteration finishes, not needed
func (it *seriesIterator) Done() {}

// Next assigns the next item and returns whether one was found
func (it *seriesIterator) Next(p *starlark.Value) bool {
	if it.count < it.series.Len() {
		*p = it.series.Index(it.count)
		it.count++
		return true
	}
	return false
}
