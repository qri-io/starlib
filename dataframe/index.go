package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// Index represents a sequence used for indexing and aligning data.
// Used for storing axis labels in a Series or DataFrame.
type Index struct {
	frozen bool
	texts  []string
	name   string
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*Index)(nil)
	_ starlark.HasAttrs = (*Index)(nil)
)

// NewIndex returns a new Index with the text values and name
func NewIndex(texts []string, name string) *Index {
	return &Index{texts: texts, name: name}
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
	result := make([]string, 0, len(i.texts))
	for _, col := range i.texts {
		// TODO(dustmop): Use proper Starlark string literal quoting, to handle
		// column names that have quotes in them.
		text := fmt.Sprintf("'%s'", col)
		result = append(result, text)
	}
	cols := strings.Join(result, ", ")
	if i.name == "" {
		return fmt.Sprintf("Index([%s], dtype='object')", cols)
	}
	return fmt.Sprintf("Index([%s], dtype='object', name='%s')", cols, i.name)
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

func (i *Index) len() int {
	if i == nil {
		return 0
	}
	return len(i.texts)
}

func newIndex(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		dataVal, nameVal starlark.Value
	)
	if err := starlark.UnpackArgs("DataFrame", args, kwargs,
		"data?", &dataVal,
		"name?", &nameVal,
	); err != nil {
		return nil, err
	}
	data := toStrSliceOrNil(dataVal)
	name := toStrOrEmpty(nameVal)
	return NewIndex(data, name), nil
}
