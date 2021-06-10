package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

type Index struct {
	texts []string
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
	return fmt.Sprintf("Index([%s], dtype='object')", strings.Join(result, ", "))
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
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (idx *Index) AttrNames() []string {
	return []string{"str"}
}
