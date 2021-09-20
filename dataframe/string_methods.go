package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// StringContainer represents something that contains a set of strings over
// which various operations may be performed
type StringContainer interface {
	Len() int
	StrAt(int) string
	CloneWithStrings([]string) starlark.Value
}

// stringMethods provides access to string methods on string collection objects
type stringMethods struct {
	subject StringContainer
}

// compile-time interface assertions
var (
	_ starlark.Value     = (*stringMethods)(nil)
	_ starlark.HasAttrs  = (*stringMethods)(nil)
	_ starlark.Indexable = (*stringMethods)(nil)
)

var stringMethodsMethods = map[string]*starlark.Builtin{
	"contains":   starlark.NewBuiltin("contains", stringMethodsContains),
	"endswith":   starlark.NewBuiltin("endswith", stringMethodsEndsWith),
	"lower":      starlark.NewBuiltin("lower", stringMethodsLower),
	"replace":    starlark.NewBuiltin("replace", stringMethodsReplace),
	"startswith": starlark.NewBuiltin("startswith", stringMethodsStartsWith),
	"strip":      starlark.NewBuiltin("strip", stringMethodsStrip),
}

// Freeze has no effect on the immutable stringMethods
func (sm *stringMethods) Freeze() {
	// pass
}

// Hash cannot be used with stringMethods
func (sm *stringMethods) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", sm.Type())
}

// String returns a string representation of the stringMethods
func (sm *stringMethods) String() string {
	return fmt.Sprintf("<%s>", sm.Type())
}

// Truth converts the stringMethods into a bool
func (sm *stringMethods) Truth() starlark.Bool {
	return true
}

// Type returns the type as a string
func (sm *stringMethods) Type() string {
	return fmt.Sprintf("%s.StringMethods", Name)
}

// Attr gets a value for a string attribute
func (sm *stringMethods) Attr(name string) (starlark.Value, error) {
	return builtinAttr(sm, name, stringMethodsMethods)
}

// AttrNames lists available attributes
func (sm *stringMethods) AttrNames() []string {
	return builtinAttrNames(stringMethodsMethods)
}

// Len returns the number of strings being operated on
func (sm *stringMethods) Len() int {
	return sm.subject.Len()
}

// Index returns a series where each element is the ith character in the subject
func (sm *stringMethods) Index(i int) starlark.Value {
	num := sm.subject.Len()
	result := make([]interface{}, 0, num)
	for k := 0; k < num; k++ {
		str := sm.subject.StrAt(k)
		r := str[i]
		str = string([]byte{r})
		result = append(result, str)
	}
	return newSeriesFromObjects(result, nil, "")
}

func stringMethodsLower(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	self := b.Receiver().(*stringMethods)

	num := self.subject.Len()
	result := make([]string, 0, num)
	for i := 0; i < num; i++ {
		lowerCol := strings.ToLower(self.subject.StrAt(i))
		result = append(result, lowerCol)
	}
	return self.subject.CloneWithStrings(result), nil
}

func stringMethodsReplace(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var oldStr, newStr starlark.String
	if err := starlark.UnpackArgs("replace", args, kwargs,
		"old", &oldStr,
		"new", &newStr,
	); err != nil {
		return nil, err
	}

	self := b.Receiver().(*stringMethods)
	oldText, _ := toStrMaybe(oldStr)
	newText, _ := toStrMaybe(newStr)

	num := self.subject.Len()
	result := make([]string, 0, num)
	for i := 0; i < num; i++ {
		replaced := strings.Replace(self.subject.StrAt(i), oldText, newText, -1)
		result = append(result, replaced)
	}
	return self.subject.CloneWithStrings(result), nil
}

func stringMethodsStrip(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	self := b.Receiver().(*stringMethods)

	num := self.subject.Len()
	result := make([]string, 0, num)
	for i := 0; i < num; i++ {
		trimmed := strings.Trim(self.subject.StrAt(i), " \t")
		result = append(result, trimmed)
	}
	return self.subject.CloneWithStrings(result), nil
}

func stringMethodsStartsWith(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var needleStr starlark.String
	if err := starlark.UnpackArgs("startswith", args, kwargs, "needle", &needleStr); err != nil {
		return nil, err
	}

	self := b.Receiver().(*stringMethods)
	needle, _ := toStrMaybe(needleStr)

	num := self.subject.Len()
	result := make([]bool, 0, num)
	for i := 0; i < num; i++ {
		b := strings.HasPrefix(self.subject.StrAt(i), needle)
		result = append(result, b)
	}
	return newSeriesFromBools(result, nil, ""), nil
}

func stringMethodsEndsWith(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var needleStr starlark.String
	if err := starlark.UnpackArgs("endswith", args, kwargs, "needle", &needleStr); err != nil {
		return nil, err
	}

	self := b.Receiver().(*stringMethods)
	needle, _ := toStrMaybe(needleStr)

	num := self.subject.Len()
	result := make([]bool, 0, num)
	for i := 0; i < num; i++ {
		b := strings.HasSuffix(self.subject.StrAt(i), needle)
		result = append(result, b)
	}
	return newSeriesFromBools(result, nil, ""), nil
}

func stringMethodsContains(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var needleStr starlark.String
	if err := starlark.UnpackArgs("contains", args, kwargs, "needle", &needleStr); err != nil {
		return nil, err
	}

	self := b.Receiver().(*stringMethods)
	needle, _ := toStrMaybe(needleStr)

	num := self.subject.Len()
	result := make([]bool, 0, num)
	for i := 0; i < num; i++ {
		b := strings.Contains(self.subject.StrAt(i), needle)
		result = append(result, b)
	}
	return newSeriesFromBools(result, nil, ""), nil
}
