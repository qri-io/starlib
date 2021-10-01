package regexp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module exposes the regexp module
var Module = &starlarkstruct.Module{
	Name: "regexp",
	Members: starlark.StringDict{
		"compile": starlark.NewBuiltin("compile", compile),
	},
}

// The regular expression used to indentify all the backreferences in a replacement pattern.
var backreferenceRe = regexp.MustCompile(`((\\\\)*)\\(\d)`)

func compile(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern string
	)

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &pattern); err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Regexp{re: re}, nil
}

func toList(strs []string) *starlark.List {
	elems := make([]starlark.Value, len(strs))
	for i, s := range strs {
		elems[i] = starlark.String(s)
	}
	return starlark.NewList(elems)
}

// A Regexp represents a compiled RE2 regular expression.
type Regexp struct {
	re *regexp.Regexp
}

// String implements the Stringer interface.
func (r *Regexp) String() string { return r.re.String() }

// Type returns a short string describing the value's type.
func (r *Regexp) Type() string { return "regexp" }

// Freeze renders r immutable. Required by starlark.Value interface.
// The interface regex presents to the Starlark runtime renders it immutable,
// making this a no-op.
func (r *Regexp) Freeze() {}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (r *Regexp) Hash() (uint32, error) { return starlark.String(r.re.String()).Hash() }

// Truth always returns true for a Regexp.
func (r *Regexp) Truth() starlark.Bool { return true }

// Attr gets a value for a string attribute, implementing dot expression support
// in Starklark. required by starlark.HasAttrs interface.
func (r *Regexp) Attr(name string) (starlark.Value, error) {
	return builtinAttr(r, name, regexMethods)
}

// AttrNames lists available dot expression strings for time. Required by
// starlark.HasAttrs interface.
func (r *Regexp) AttrNames() []string { return builtinAttrNames(regexMethods) }

var regexMethods = map[string]*starlark.Builtin{
	"find":            starlark.NewBuiltin("find", find),
	"find_all":        starlark.NewBuiltin("find_all", findAll),
	"find_submatches": starlark.NewBuiltin("find_submatches", findSubmatches),
	"matches":         starlark.NewBuiltin("matches", matches),
	"replace_all":     starlark.NewBuiltin("replace_all", replaceAll),
	"split":           starlark.NewBuiltin("split", split),
}

func builtinAttr(recv starlark.Value, name string, methods map[string]*starlark.Builtin) (starlark.Value, error) {
	b := methods[name]
	if b == nil {
		return nil, nil // no such method
	}
	return b.BindReceiver(recv), nil
}

func builtinAttrNames(methods map[string]*starlark.Builtin) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func matches(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src string

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	return starlark.Bool(re.MatchString(src)), nil
}

func find(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src string

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	return starlark.String(re.FindString(src)), nil
}

func findAll(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		src string
		max int = -1
	)

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &src, &max); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	return toList(re.FindAllString(src, max)), nil
}

func findSubmatches(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src string

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	return toList(re.FindStringSubmatch(src)), nil
}

func replaceAll(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		src  string
		repl starlark.Value
	)

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &src, &repl); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	switch x := repl.(type) {
	case starlark.Callable:
		var fnErr error
		fn := func(matched string) string {
			res, err := starlark.Call(thread, repl, starlark.Tuple{starlark.String(matched)}, nil)
			if err != nil {
				// Save the error to be able to return it to the caller
				fnErr = err
				return ""
			}
			resp, ok := res.(starlark.String)
			if !ok {
				// Save the error to be able to return it to the caller
				fnErr = fmt.Errorf("%s returned %s, want string", x.Name(), res.Type())
				return ""
			}
			return string(resp)
		}
		result := re.ReplaceAllStringFunc(src, fn)
		if fnErr != nil {
			return nil, fnErr
		}
		return starlark.String(result), nil
	case starlark.String:
		return starlark.String(re.ReplaceAllString(src, convertReplacementPattern(string(x)))), nil
	}
	return nil, fmt.Errorf("got %s, want a string or callable", repl.Type())
}

func split(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		src string
		max int = -1
	)

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &src, &max); err != nil {
		return nil, err
	}

	re := b.Receiver().(*Regexp).re
	return toList(re.Split(src, max)), nil
}

// Converts the replacement pattern to make it compatible with what the regexp package expects.
func convertReplacementPattern(repl string) string {
	// Escapes the special dollar characters if any
	repl = strings.ReplaceAll(repl, "$", "$$")
	startIdx := 0
	var sb strings.Builder
	// Replaces the backreferences represented by backslash-escaped digits (\1 to \9) with their
	// counter parts (${1} to ${9}).
	for _, match := range backreferenceRe.FindAllStringSubmatchIndex(repl, -1) {
		startIdxMatch := match[0]
		endIdxMatch := match[1]
		if startIdxMatch > 0 && repl[startIdxMatch-1] == '\\' {
			// The preceding character is a slash so we don't have a real match let's skip it
			// Unescapes slashes to get the expected result
			sb.WriteString(strings.ReplaceAll(repl[startIdx:endIdxMatch], `\\`, `\`))
			startIdx = endIdxMatch
			continue

		}
		endIdxGroup1 := match[3]
		if startIdx < endIdxGroup1 {
			// Adds the slashes of the first group and unscapes them
			sb.WriteString(strings.ReplaceAll(repl[startIdx:endIdxGroup1], `\\`, `\`))
		}
		// Builds the backreference from the third group
		sb.WriteString("${")
		sb.WriteString(repl[match[6]:match[7]])
		sb.WriteString("}")
		startIdx = endIdxMatch
	}
	if startIdx < len(repl) {
		// Adds the remaining characters if any
		sb.WriteString(repl[startIdx:])
	}
	return sb.String()
}
