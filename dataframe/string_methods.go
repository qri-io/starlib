package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

type stringMethods struct {
	subject *Index
}

// Freeze ...
func (sm *stringMethods) Freeze() {
	// pass
}

func (sm *stringMethods) Hash() (uint32, error) {
	// TODO
	return 0, nil
}

func (sm *stringMethods) String() string {
	return "<class 'StringMethods'>"
}

// Truth ...
func (sm *stringMethods) Truth() starlark.Bool {
	return true
}

func (sm *stringMethods) Type() string {
	return "dataframe.StringMethods"
}

type starlarkFunc func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

// Attr
func (sm *stringMethods) Attr(name string) (starlark.Value, error) {
	if name == "strip" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return sm.stringStrip(b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "lower" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return sm.stringLower(b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	if name == "replace" {
		impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return sm.stringReplace(b.Name(), args, kwargs)
		}
		return starlark.NewBuiltin(name, impl), nil
	}
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (sm *stringMethods) AttrNames() []string {
	return []string{"strip", "lower", "replace"}
}

//
func (sm *stringMethods) stringStrip(fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	columns := sm.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.Trim(text, " \t")
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}

//
func (sm *stringMethods) stringLower(fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	columns := sm.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.ToLower(text)
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}

//
func (sm *stringMethods) stringReplace(fnname string, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		oldVal, newVal starlark.Value
	)
	if err := starlark.UnpackArgs("replace", args, kwargs,
		"old", &oldVal,
		"new", &newVal,
	); err != nil {
		return nil, err
	}

	oldStr, ok := oldVal.(starlark.String)
	if !ok {
		return starlark.None, fmt.Errorf("invalid conversion 'old' to string")
	}
	newStr, ok := newVal.(starlark.String)
	if !ok {
		return starlark.None, fmt.Errorf("invalid conversion 'new' to string")
	}
	oldText := string(oldStr)
	newText := string(newStr)

	columns := sm.subject.texts
	result := make([]string, 0, len(columns))
	for _, text := range columns {
		lowerCol := strings.Replace(text, oldText, newText, -1)
		result = append(result, lowerCol)
	}
	return &Index{texts: result}, nil
}
