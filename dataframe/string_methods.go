package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// stringMethods provides access to string methods on string collection objects
type stringMethods struct {
	subject *Index
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*stringMethods)(nil)
	_ starlark.HasAttrs = (*stringMethods)(nil)
)

var stringMethodsMethods = map[string]*starlark.Builtin{
	"lower":   starlark.NewBuiltin("lower", stringMethodsLower),
	"replace": starlark.NewBuiltin("replace", stringMethodsReplace),
	"strip":   starlark.NewBuiltin("strip", stringMethodsStrip),
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

func stringMethodsLower(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	self := b.Receiver().(*stringMethods)
	columns := self.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.ToLower(text)
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}

func stringMethodsReplace(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		oldVal, newVal starlark.Value
	)
	if err := starlark.UnpackArgs("replace", args, kwargs,
		"old", &oldVal,
		"new", &newVal,
	); err != nil {
		return nil, err
	}

	self := b.Receiver().(*stringMethods)
	oldText, _ := toStrMaybe(oldVal)
	newText, _ := toStrMaybe(newVal)

	columns := self.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.Replace(text, oldText, newText, -1)
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}

func stringMethodsStrip(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	self := b.Receiver().(*stringMethods)
	columns := self.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.Trim(text, " \t")
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}
