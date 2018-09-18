package html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
	"github.com/qri-io/starlib/util"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('html.sky', 'html')
const ModuleName = "html.sky"

// LoadModule loads the html module
func LoadModule() (skylark.StringDict, error) {
	return skylark.StringDict{
		"html": skylark.NewBuiltin("html", NewDocument),
	}, nil
}

// NewDocument creates a skylark selection from input text
func NewDocument(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var body skylark.String
	if err := skylark.UnpackArgs("html", args, kwargs, "body", &body); err != nil {
		return nil, err
	}

	str, err := util.AsString(body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	return NewSelectionStruct(doc.Selection), err
}

// Selection is a wrapper for a goquery selection mapping to skylark values
type Selection struct {
	sel *goquery.Selection
}

// NewSelectionStruct creates a skylark struct from a goquery selection
func NewSelectionStruct(s *goquery.Selection) *skylarkstruct.Struct {
	sel := &Selection{sel: s}
	return sel.Struct()
}

// Struct returns a skylark struct of methods
func (s *Selection) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylarkstruct.Default, skylark.StringDict{
		"attr":              skylark.NewBuiltin("attr", s.Attr),
		"children":          skylark.NewBuiltin("children", s.Children),
		"children_filtered": skylark.NewBuiltin("children_filtered", s.ChildrenFiltered),
		"contents":          skylark.NewBuiltin("contents", s.Contents),
		"find":              skylark.NewBuiltin("find", s.Find),
		"filter":            skylark.NewBuiltin("filter", s.Filter),
		"get":               skylark.NewBuiltin("get", s.Get),
		"has":               skylark.NewBuiltin("has", s.Has),
		"parent":            skylark.NewBuiltin("parent", s.Parent),
		"parents_until":     skylark.NewBuiltin("parents_until", s.ParentsUntil),
		"siblings":          skylark.NewBuiltin("siblings", s.Siblings),
		"text":              skylark.NewBuiltin("text", s.Text),
		"first":             skylark.NewBuiltin("first", s.First),
		"last":              skylark.NewBuiltin("last", s.Last),
		"len":               skylark.NewBuiltin("len", s.Len),
		"eq":                skylark.NewBuiltin("eq", s.Eq),
	})
}

// Attr gets the specified attribute's value for the first element in the Selection. To get the value for each element individually, use a looping construct such as Each or Map method
func (s *Selection) Attr(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("attr", args, kwargs)
	if err != nil {
		return nil, err
	}
	val, exists := s.sel.Attr(sstr)
	if !exists {
		return skylark.None, nil
	}
	return skylark.String(val), nil
}

// Children gets the child elements of each element in the Selection. It returns a new Selection object containing these elements
func (s *Selection) Children(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return NewSelectionStruct(s.sel.Children()), nil
}

// ChildrenFiltered gets the child elements of each element in the Selection, filtered by the specified selector. It returns a new Selection object containing these elements
func (s *Selection) ChildrenFiltered(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("children_filtered", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.ChildrenFiltered(sstr)
	return NewSelectionStruct(sel), nil
}

// Contents gets the children of each element in the Selection, including text and comment nodes. It returns a new Selection object containing these elements
func (s *Selection) Contents(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return NewSelectionStruct(s.sel.Contents()), nil
}

// // Each
// func (s *Selection) Each(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// }

// Get retrieves the underlying node at the specified index. Get without parameter is not implemented, since the node array is available on the Selection object
func (s *Selection) Get(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Int
	if err := skylark.UnpackPositionalArgs("get", args, kwargs, 1, &x); err != nil {
		var t skylark.Tuple
		for _, node := range s.sel.Nodes {
			t = append(t, skylark.String(node.Data))
		}
		return t, nil
	}

	i, _ := x.Int64()
	if int(i) > len(s.sel.Nodes)-1 {
		return skylark.None, nil
	}
	sel := s.sel.Get(int(i))
	return skylark.String(sel.Data), nil
}

// Find gets the descendants of each element in the current set of matched elements, filtered by a selector. It returns a new Selection object containing these matched element
func (s *Selection) Find(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("find", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Find(sstr)
	return NewSelectionStruct(sel), nil
}

// Filter reduces the set of matched elements to those that match the selector string. It returns a new Selection object for this subset of matching elements
func (s *Selection) Filter(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("filter", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Filter(sstr)
	return NewSelectionStruct(sel), nil
}

// Has reduces the set of matched elements to those that have a descendant that matches the selector. It returns a new Selection object with the matching elements
func (s *Selection) Has(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("has", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.Has(sstr)
	return NewSelectionStruct(sel), nil
}

// Parent gets the parent of each element in the Selection. It returns a new Selection object containing the matched elements
func (s *Selection) Parent(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return NewSelectionStruct(s.sel.Parent()), nil
}

// ParentsUntil gets the ancestors of each element in the Selection, up to but not including the element matched by the selector. It returns a new Selection object containing the matched elements
func (s *Selection) ParentsUntil(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sstr, err := s.selectorArg("parents_until", args, kwargs)
	if err != nil {
		return nil, err
	}
	sel := s.sel.ParentsUntil(sstr)
	return NewSelectionStruct(sel), nil
}

// selectorArg is a convenience method for functions that only accept a string selector
func (s *Selection) selectorArg(method string, args skylark.Tuple, kwargs []skylark.Tuple) (string, error) {
	var selector skylark.String
	if err := skylark.UnpackPositionalArgs(method, args, kwargs, 1, &selector); err != nil {
		return "", err
	}
	return util.AsString(selector)
}

// Siblings gets the siblings of each element in the Selection. It returns a new Selection object containing the matched elements
func (s *Selection) Siblings(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sel := s.sel.Siblings()
	return NewSelectionStruct(sel), nil
}

// Text gets the combined text contents of each element in the set of matched elements, including their descendants
func (s *Selection) Text(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return skylark.String(s.sel.Text()), nil
}

// First gets the first element of the selection
func (s *Selection) First(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return NewSelectionStruct(s.sel.First()), nil
}

// Last gets the last element of the selection
func (s *Selection) Last(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return NewSelectionStruct(s.sel.Last()), nil
}

// Eq gets the element i of the selection
func (s *Selection) Eq(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Int
	if err := skylark.UnpackPositionalArgs("eq", args, kwargs, 1, &x); err != nil {
		return nil, err
	}

	i, _ := x.Int64()
	return NewSelectionStruct(s.sel.Eq(int(i))), nil
}

// Len returns the length of the nodes in the selection
func (s *Selection) Len(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return skylark.MakeInt(len(s.sel.Nodes)), nil
}
