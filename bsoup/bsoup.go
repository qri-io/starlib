package bsoup

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/dustmop/soup"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"golang.org/x/net/html"
)

// ModuleName defines the name for loading this module, using `load('bsoup.star', 'bsoup')
const ModuleName = "bsoup.star"

var (
	once        sync.Once
	bsoupModule starlark.StringDict
)

// LoadModule loads the bsoup module. Concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		bsoupModule = starlark.StringDict{
			"bsoup": &starlarkstruct.Module{
				Name: "bsoup",
				Members: starlark.StringDict{
					"parseHtml": starlark.NewBuiltin("parseHtml", ParseHTML),
				},
			},
		}
	})
	return bsoupModule, nil
}

// ParseHTML parses html from a string, and returns it as a SoupNode
func ParseHTML(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var htmlText starlark.String
	if err := starlark.UnpackArgs("parseHtml", args, kwargs, "htmlText", &htmlText); err != nil {
		return nil, err
	}

	content, err := AsString(htmlText)
	if err != nil {
		return nil, err
	}

	root := soup.HTMLParse(string(content))
	return NewSoupNode(&root), nil
}

// AsString converts a starlark Value into a string, with outer quotes trimmed
func AsString(x starlark.Value) (string, error) {
	if x == nil || x == starlark.None {
		return "nil", nil
	}
	return strconv.Unquote(x.String())
}

// SoupNode extends soup's Root struct with starlark support
type SoupNode soup.Root

// String converts a SoupNode to a string by rendering each node
func (n *SoupNode) String() string {
	buf := new(bytes.Buffer)
	err := html.Render(buf, n.Pointer)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// Type returns the type of SoupNode as string
func (n *SoupNode) Type() string {
	return "SoupNode"
}

// Freeze freezes a SoupNode struct, which is already immutable
func (n *SoupNode) Freeze() {}

// Hash calculates a hash of a SoupNode
func (n *SoupNode) Hash() (uint32, error) {
	return hashString(fmt.Sprintf("%v", *n)), nil
}

// Truth returns whether a SoupNode is non-nil
func (n *SoupNode) Truth() starlark.Bool {
	return n != nil
}

// Attr returns an attribute of a SoupNode
func (n *SoupNode) Attr(name string) (starlark.Value, error) {
	method := bsoupMethods[name]
	if method == nil {
		return nil, nil // no such method
	}

	impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return method(b.Name(), b.Receiver(), args, kwargs)
	}
	return starlark.NewBuiltin(name, impl).BindReceiver(n), nil
}

// AttrNames returns all attributes of a SoupNode
func (n *SoupNode) AttrNames() []string {
	names := make([]string, 0, len(bsoupMethods))
	for name := range bsoupMethods {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type builtinMethod func(fnname string, recv starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func hashString(s string) uint32 {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

var bsoupMethods = map[string]builtinMethod{
	"find":         bsoupFind,
	"find_all":     bsoupFindAll,
	"attrs":        bsoupAttrs,
	"contents":     bsoupContents,
	"child":        bsoupChild,
	"parent":       bsoupParent,
	"next_sibling": bsoupNextSibling,
	"prev_sibling": bsoupPrevSibling,
	"get_text":     bsoupGetText,
	// TODO(dlong):
	// .string https://www.crummy.com/software/BeautifulSoup/bs4/doc/#string
	// .strings https://www.crummy.com/software/BeautifulSoup/bs4/doc/#strings-and-stripped-strings
	// .parents https://www.crummy.com/software/BeautifulSoup/bs4/doc/#parents
	// .next_element and .prev_element https://www.crummy.com/software/BeautifulSoup/bs4/doc/#next-element-and-previous-element
	// find(string=...) https://www.crummy.com/software/BeautifulSoup/bs4/doc/#the-string-argument
	// find(limit=...) https://www.crummy.com/software/BeautifulSoup/bs4/doc/#the-limit-argument
	// .find_parents https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find-parents-and-find-parent
	// .select https://www.crummy.com/software/BeautifulSoup/bs4/doc/#css-selectors
	// .prettify https://www.crummy.com/software/BeautifulSoup/bs4/doc/#pretty-printing
}

// NewSoupNode constructs a new SoupNode by cloning each field from the soup.Root
func NewSoupNode(root *soup.Root) starlark.Value {
	// Need to clone, since the input value is not immutable. Removing the clone operation
	// will break things like iteration, wherein a single soup.Root will be mutated by
	// each step of the loop body.
	// Luckily, the soup library exports the field names of this struct, so we can clone.
	clone := &soup.Root{
		Pointer:   root.Pointer,
		NodeValue: root.NodeValue,
		Error:     root.Error,
	}

	return (*SoupNode)(clone)
}

// parseFindArgs converts starlark function arguments into a linear list. BeautifulSoup has very
// flexible ways of passing arguments to `find` and `find_alll`, so we have to do this manually.
func (n *SoupNode) parseFindArgs(args starlark.Tuple, kwargs []starlark.Tuple) ([]string, error) {
	params := []string{}
	haveTagName := false

	// Convert positional arguments to a string list. Each argument may be a string (for a tagName),
	// or a dictionary, which will get flatteneed. Only one tagName is allowed.
	for _, arg := range args {
		if tagName, ok := arg.(starlark.String); ok {
			if haveTagName {
				return nil, fmt.Errorf("only one tagName parameter is allowed, found %s", tagName)
			}
			str, err := AsString(tagName)
			if err != nil {
				return nil, err
			}
			params = append(params, str)
			haveTagName = true
			continue
		}
		if dict, ok := arg.(*starlark.Dict); ok {
			for _, k := range dict.Keys() {
				key, err := AsString(k)
				if err != nil {
					return nil, err
				}
				params = append(params, key)
				v, _, err := dict.Get(k)
				if err != nil {
					return nil, err
				}
				val, err := AsString(v)
				if err != nil {
					return nil, err
				}
				params = append(params, val)
			}
			continue
		}
		return nil, fmt.Errorf("invalid parameter %v", arg)
	}
	// If no tagName was given, prepend a blank string, meaning any tagName.
	if !haveTagName {
		params = append([]string{""}, params...)
	}

	// Convert keyword arguments to a string list, by flattening them.
	// TODO(dlong): Handle meaningful keywords: name, attrs, recurive, string, limit, class_
	for _, kw := range kwargs {
		key, err := AsString(kw[0])
		if err != nil {
			return nil, err
		}
		params = append(params, key)
		val, err := AsString(kw[1])
		if err != nil {
			return nil, err
		}
		params = append(params, val)
	}

	return params, nil
}

// bsoupFind implements soup.find()
func bsoupFind(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	params, err := node.parseFindArgs(args, kwargs)
	if err != nil {
		return starlark.None, err
	}

	elem := (*soup.Root)(node).Find(params...)
	if elem.Pointer == nil {
		return starlark.None, nil
	}

	return NewSoupNode(&elem), nil
}

// bsoupFind implements soup.find_all()
func bsoupFindAll(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	params, err := node.parseFindArgs(args, kwargs)
	if err != nil {
		return starlark.None, err
	}

	elemList := (*soup.Root)(node).FindAll(params...)
	nodeList := starlark.NewList([]starlark.Value{})
	for _, elem := range elemList {
		built := NewSoupNode(&elem)
		nodeList.Append(built)
	}

	return nodeList, nil
}

// bsoupAttrs implements soup.attrs
func bsoupAttrs(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	attrs := (*soup.Root)(node).Attrs()
	result := starlark.NewDict(0)
	for k, v := range attrs {
		result.SetKey(starlark.String(k), starlark.String(v))
	}
	return result, nil
}

// bsoupAttrs implements soup.contents
func bsoupContents(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	children := (*soup.Root)(node).Children()
	nodeList := starlark.NewList([]starlark.Value{})
	for _, elem := range children {
		built := NewSoupNode(&elem)
		nodeList.Append(built)
	}
	return nodeList, nil
}

// bsoupChild is a replacement for things like soup.title. Instead do soup.child('title')
func bsoupChild(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)

	var tagName starlark.String
	if err := starlark.UnpackPositionalArgs("Child", args, kwargs, 1, &tagName); err != nil {
		return starlark.None, err
	}

	name, err := AsString(tagName)
	if err != nil {
		return starlark.None, err
	}

	children := (*soup.Root)(node).Children()
	for _, elem := range children {
		if elem.Pointer.Type == html.ElementNode && elem.Pointer.Data == name {
			return NewSoupNode(&elem), nil
		}
	}
	return starlark.None, nil
}

// bsoupParent implements soup.parent
func bsoupParent(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	parent := node.Pointer.Parent
	result := NewSoupNode(&soup.Root{
		Pointer: parent,
	})
	return result, nil
}

// bsoupNextSibling implements soup.next_sibling
func bsoupNextSibling(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	sibling := (*soup.Root)(node).FindNextSibling()
	return NewSoupNode(&sibling), nil
}

// bsoupPrevSibling implements soup.prev_sibling
func bsoupPrevSibling(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	sibling := (*soup.Root)(node).FindPrevSibling()
	return NewSoupNode(&sibling), nil
}

// bsoupGetText implements soup.get_text
// https://www.crummy.com/software/BeautifulSoup/bs4/doc/#get-text
func bsoupGetText(fnname string, self starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	node := self.(*SoupNode)
	text := (*soup.Root)(node).Text()
	return starlark.String(text), nil
}
