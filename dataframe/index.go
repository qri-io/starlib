package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

type Index struct {
	texts []string
	name  string
}

func NewIndex(texts []string, name string) *Index {
	return &Index{texts: texts, name: name}
}

// Freeze ...
func (idx *Index) Freeze() {
	// pass
}

func (idx *Index) Hash() (uint32, error) {
	// TODO
	return 0, nil
}

func (idx *Index) String() string {
	result := make([]string, 0, len(idx.texts))
	for _, col := range idx.texts {
		text := fmt.Sprintf("'%s'", col)
		result = append(result, text)
	}
	cols := strings.Join(result, ", ")
	if idx.name == "" {
		return fmt.Sprintf("Index([%s], dtype='object')", cols)
	}
	return fmt.Sprintf("Index([%s], dtype='object', name='%s')", cols, idx.name)
}

// Truth ...
func (idx *Index) Truth() starlark.Bool {
	return true
}

func (idx *Index) Type() string {
	return "dataframe.Index"
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (idx *Index) Attr(name string) (starlark.Value, error) {
	if name == "str" {
		return &stringMethods{subject: idx}, nil
	}
	if name == "name" {
		return starlark.String(idx.name), nil
	}
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (idx *Index) AttrNames() []string {
	return []string{"str", "name"}
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
	data := toStrListOrNil(dataVal)
	name := toStrOrEmpty(nameVal)
	return NewIndex(data, name), nil
}
