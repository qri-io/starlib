package dataframe

import (
	"fmt"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

// Index represents a sequence used for indexing and aligning data.
// Used for storing axis labels in a Series or DataFrame.
type Index struct {
	frozen bool
	impl   indexImpl
	name   string
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*Index)(nil)
	_ starlark.HasAttrs = (*Index)(nil)
	_ starlark.Sequence = (*Index)(nil)
)

// NewObjIndex returns a new Index with the values and name
func NewObjIndex(objs []interface{}, name string) *Index {
	return &Index{impl: newObjIndexImpl(objs), name: name}
}

// NewTextIndex returns a new Index with the strings and name
func NewTextIndex(texts []string, name string) *Index {
	return &Index{impl: newObjIndexImpl(convertStringsToObjects(texts)), name: name}
}

// NewRangeIndex returns a new Index with given range of ints
func NewRangeIndex(size int, name string) *Index {
	return &Index{impl: newRangeIndexImpl(size), name: name}
}

// NewInt64Index returns a new Index of integer values, with a name
func NewInt64Index(nums []int, name string) *Index {
	return &Index{impl: newInt64IndexImpl(nums), name: name}
}

// construct a new index, of ints if possible, otherwise objects
func newIndexFrom(vals []interface{}, name string) *Index {
	tryNums := make([]int, len(vals))
	for k, v := range vals {
		if n, ok := v.(int); ok {
			tryNums[k] = n
		} else {
			tryNums = nil
			break
		}
	}
	if tryNums != nil {
		return NewInt64Index(tryNums, name)
	}
	return NewObjIndex(vals, name)
}

// CloneWithStrings returns a clone of the index but with replaced string values
func (i *Index) CloneWithStrings(texts []string) starlark.Value {
	return NewTextIndex(texts, i.name)
}

// Freeze prevents the index from being mutated
func (i *Index) Freeze() {
	i.frozen = true
}

// Hash cannot be used with Index
func (i *Index) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", i.Type())
}

// String returns the index as a string
func (i *Index) String() string {
	if i == nil || i.impl == nil {
		return "<nil>"
	}
	typ := i.impl.Type()
	cols := i.impl.ColumnsString()
	if i.name == "" {
		return fmt.Sprintf("%s(%s)", typ, cols)
	}
	return fmt.Sprintf("%s(%s, name='%s')", typ, cols, i.name)
}

// Truth converts the index into a bool
func (i *Index) Truth() starlark.Bool {
	// NOTE: In python, calling bool(Index) raises this exception: "ValueError: The truth
	// value of a Index is ambiguous. Use a.empty, a.bool(), a.item(), a.any() or a.all()."
	// Since starlark does not have exceptions, just always return true.
	return true
}

// Type returns the type as a string
func (i *Index) Type() string {
	return fmt.Sprintf("%s.Index", Name)
}

// Attr gets a value for a string attribute
func (i *Index) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(i.name), nil
	case "str":
		return &stringMethods{subject: i}, nil
	}
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings
func (i *Index) AttrNames() []string {
	return []string{"name", "str"}
}

// Iterate returns an iterator for the index
func (i *Index) Iterate() starlark.Iterator {
	return &indexIterator{idx: i, count: 0}
}

// Len returns the length of the index
func (i *Index) Len() int {
	if i == nil || i.impl == nil {
		return 0
	}
	return i.impl.Len()
}

// StrAt returns the string at index k
func (i *Index) StrAt(k int) string {
	if i == nil {
		return strconv.Itoa(k)
	}
	return i.impl.StrAt(k)
}

// At returns the data at index k
func (i *Index) At(k int) interface{} {
	if i == nil {
		return k
	}
	return i.impl.At(k)
}

// Columns returns the columns as string values
func (i *Index) Columns() []string {
	result := make([]string, i.impl.Len())
	for k := 0; k < i.impl.Len(); k++ {
		result[k] = i.impl.StrAt(k)
	}
	return result
}

func newIndex(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		dataVal, nameVal starlark.Value
	)
	if err := starlark.UnpackArgs("Index", args, kwargs,
		"data?", &dataVal,
		"name?", &nameVal,
	); err != nil {
		return nil, err
	}
	data := toInterfaceSliceOrNil(dataVal)
	name := toStrOrEmpty(nameVal)
	return newIndexFrom(data, name), nil
}

type indexIterator struct {
	count int
	idx   *Index
}

// Done does cleanup work when iteration finishes, not needed
func (it *indexIterator) Done() {}

// Next assigns the next item and returns whether one was found
func (it *indexIterator) Next(p *starlark.Value) bool {
	if it.count < it.idx.Len() {
		*p = starlark.String(it.idx.StrAt(it.count))
		it.count++
		return true
	}
	return false
}

// interface for implementations of the index
type indexImpl interface {
	Type() string
	ColumnsString() string
	Len() int
	StrAt(int) string
	At(int) interface{}
}

// objects for an index implementation
type objIndexImpl struct {
	objs []interface{}
}

func newObjIndexImpl(objs []interface{}) *objIndexImpl {
	return &objIndexImpl{objs: objs}
}

func (ti *objIndexImpl) Type() string {
	return "Index"
}

func (ti *objIndexImpl) ColumnsString() string {
	result := make([]string, ti.Len())
	for i, col := range ti.objs {
		if str, ok := col.(string); ok {
			// TODO(dustmop): Use proper Starlark string literal quoting, to handle
			// column names that have quotes in them.
			result[i] = fmt.Sprintf("'%s'", str)
			continue
		}
		result[i] = fmt.Sprintf("%v", col)
	}
	return fmt.Sprintf("[%s], dtype='object'", strings.Join(result, ", "))
}

func (ti *objIndexImpl) Len() int {
	return len(ti.objs)
}

func (ti *objIndexImpl) StrAt(k int) string {
	return fmt.Sprintf("%v", ti.objs[k])
}

func (ti *objIndexImpl) At(k int) interface{} {
	return ti.objs[k]
}

// range from 0 up to some limit for an index implementation
type rangeIndexImpl struct {
	size int
}

func newRangeIndexImpl(size int) *rangeIndexImpl {
	return &rangeIndexImpl{size: size}
}

func (ri *rangeIndexImpl) Type() string {
	return "RangeIndex"
}

func (ri *rangeIndexImpl) ColumnsString() string {
	return fmt.Sprintf("start=0, stop=%d, step=1", ri.size)
}

func (ri *rangeIndexImpl) Len() int {
	return ri.size
}

func (ri *rangeIndexImpl) StrAt(k int) string {
	return fmt.Sprintf("%v", k)
}

func (ri *rangeIndexImpl) At(k int) interface{} {
	return k
}

// list of integer values for an index implementation
type int64IndexImpl struct {
	nums []int
}

func newInt64IndexImpl(nums []int) *int64IndexImpl {
	return &int64IndexImpl{nums: nums}
}

func (ii *int64IndexImpl) Type() string {
	return "Int64Index"
}

func (ii *int64IndexImpl) ColumnsString() string {
	result := make([]string, len(ii.nums))
	for i, n := range ii.nums {
		result[i] = fmt.Sprintf("%v", n)
	}
	return fmt.Sprintf("[%s], dtype='int64'", strings.Join(result, ", "))
}

func (ii *int64IndexImpl) Len() int {
	return len(ii.nums)
}

func (ii *int64IndexImpl) StrAt(k int) string {
	return fmt.Sprintf("%v", ii.nums[k])
}

func (ii *int64IndexImpl) At(k int) interface{} {
	return ii.nums[k]
}
