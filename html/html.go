package html

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// AsString unquotes a starlark string value
func AsString(x starlark.Value) (string, error) {
	return strconv.Unquote(x.String())
}

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('html.star', 'html')
const ModuleName = "html.star"

// LoadModule loads the html module
func LoadModule() (starlark.StringDict, error) {
	return starlark.StringDict{
		"html": starlark.NewBuiltin("html", NewDocument),
	}, nil
}

// NewDocument creates a starlark selection from input text
func NewDocument(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var body starlark.String
	if err := starlark.UnpackArgs("html", args, kwargs, "body", &body); err != nil {
		return nil, err
	}

	str, err := AsString(body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	return NewSelectionStruct(doc.Selection), err
}

// Selection is a wrapper for a goquery selection mapping to starlark values
type Selection struct {
	sel *goquery.Selection
}

// NewSelectionStruct creates a starlark struct from a goquery selection
func NewSelectionStruct(s *goquery.Selection) *starlarkstruct.Struct {
	sel := &Selection{sel: s}
	return sel.Struct()
}

// Struct returns a starlark struct of methods
func (s *Selection) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		"attr":              starlark.NewBuiltin("attr", s.Attr),
		"children":          starlark.NewBuiltin("children", s.Children),
		"children_filtered": starlark.NewBuiltin("children_filtered", s.ChildrenFiltered),
		"contents":          starlark.NewBuiltin("contents", s.Contents),
		"find":              starlark.NewBuiltin("find", s.Find),
		"filter":            starlark.NewBuiltin("filter", s.Filter),
		"get":               starlark.NewBuiltin("get", s.Get),
		"has":               starlark.NewBuiltin("has", s.Has),
		"parent":            starlark.NewBuiltin("parent", s.Parent),
		"parents_until":     starlark.NewBuiltin("parents_until", s.ParentsUntil),
		"siblings":          starlark.NewBuiltin("siblings", s.Siblings),
		"text":              starlark.NewBuiltin("text", s.Text),
		"first":             starlark.NewBuiltin("first", s.First),
		"last":              starlark.NewBuiltin("last", s.Last),
		"len":               starlark.NewBuiltin("len", s.Len),
		"eq":                starlark.NewBuiltin("eq", s.Eq),
	})
}

// Attr gets the specified attribute's value for the first element in the Selection. To get the value for each element individually, use a looping construct such as Each or Map method
func (s *Selection) Attr(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("attr", args, kwargs)
	if err != nil {
		return nil, err
	}
	val, exists := s.sel.Attr(sstr)
	if !exists {
		return starlark.None, nil
	}
	return starlark.String(val), nil
}

// Children gets the child elements of each element in the Selection. It returns a new Selection object containing these elements
func (s *Selection) Children(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewSelectionStruct(s.sel.Children()), nil
}

// ChildrenFiltered gets the child elements of each element in the Selection, filtered by the specified selector. It returns a new Selection object containing these elements
func (s *Selection) ChildrenFiltered(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("children_filtered", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.ChildrenFiltered(sstr)
	return NewSelectionStruct(sel), nil
}

// Contents gets the children of each element in the Selection, including text and comment nodes. It returns a new Selection object containing these elements
func (s *Selection) Contents(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewSelectionStruct(s.sel.Contents()), nil
}

// // Each
// func (s *Selection) Each(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// }

// Get retrieves the underlying node at the specified index. Get without parameter is not implemented, since the node array is available on the Selection object
func (s *Selection) Get(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Int
	if err := starlark.UnpackPositionalArgs("get", args, kwargs, 1, &x); err != nil {
		var t starlark.Tuple
		for _, node := range s.sel.Nodes {
			t = append(t, starlark.String(node.Data))
		}
		return t, nil
	}

	i, _ := x.Int64()
	if int(i) > len(s.sel.Nodes)-1 {
		return starlark.None, nil
	}
	sel := s.sel.Get(int(i))
	return starlark.String(sel.Data), nil
}

// Find gets the descendants of each element in the current set of matched elements, filtered by a selector. It returns a new Selection object containing these matched element
func (s *Selection) Find(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("find", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Find(sstr)
	return NewSelectionStruct(sel), nil
}

// Filter reduces the set of matched elements to those that match the selector string. It returns a new Selection object for this subset of matching elements
func (s *Selection) Filter(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("filter", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Filter(sstr)
	return NewSelectionStruct(sel), nil
}

// Has reduces the set of matched elements to those that have a descendant that matches the selector. It returns a new Selection object with the matching elements
func (s *Selection) Has(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("has", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Has(sstr)
	return NewSelectionStruct(sel), nil
}

// Parent gets the parent of each element in the Selection. It returns a new Selection object containing the matched elements
func (s *Selection) Parent(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewSelectionStruct(s.sel.Parent()), nil
}

// ParentsUntil gets the ancestors of each element in the Selection, up to but not including the element matched by the selector. It returns a new Selection object containing the matched elements
func (s *Selection) ParentsUntil(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sstr, err := s.selectorArg("parents_until", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.ParentsUntil(sstr)
	return NewSelectionStruct(sel), nil
}

// selectorArg is a convenience method for functions that only accept a string selector
func (s *Selection) selectorArg(method string, args starlark.Tuple, kwargs []starlark.Tuple) (string, error) {
	var selector starlark.String
	if err := starlark.UnpackPositionalArgs(method, args, kwargs, 1, &selector); err != nil {
		return "", err
	}
	return AsString(selector)
}

// Siblings gets the siblings of each element in the Selection. It returns a new Selection object containing the matched elements
func (s *Selection) Siblings(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sel := s.sel.Siblings()
	return NewSelectionStruct(sel), nil
}

// Text gets the combined text contents of each element in the set of matched elements, including their descendants
func (s *Selection) Text(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.String(s.sel.Text()), nil
}

// First gets the first element of the selection
func (s *Selection) First(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewSelectionStruct(s.sel.First()), nil
}

// Last gets the last element of the selection
func (s *Selection) Last(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewSelectionStruct(s.sel.Last()), nil
}

// Eq gets the element i of the selection
func (s *Selection) Eq(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Int
	if err := starlark.UnpackPositionalArgs("eq", args, kwargs, 1, &x); err != nil {
		return nil, err
	}

	i, _ := x.Int64()
	return NewSelectionStruct(s.sel.Eq(int(i))), nil
}

// Len returns the length of the nodes in the selection
func (s *Selection) Len(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.MakeInt(len(s.sel.Nodes)), nil
}
