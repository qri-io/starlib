package re

import (
	"regexp"
	"sync"

	"github.com/qri-io/starlib/util"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('re.star', 're')
const ModuleName = "re.star"

var (
	once     sync.Once
	reModule starlark.StringDict
)

// LoadModule loads the re module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		reModule = starlark.StringDict{
			"re": &starlarkstruct.Module{
				Name: "re",
				Members: starlark.StringDict{
					// TODO
					// "compile" : starlark.NewBuiltin("complile", compile),

					"search": starlark.NewBuiltin("search", search),
					"match":  starlark.NewBuiltin("match", match),
					// "fullmatch": starlark.NewBuiltin("fullmatch", fullmatch),
					"split":   starlark.NewBuiltin("split", split),
					"findall": starlark.NewBuiltin("findall", findall),
					// "finditer":  starlark.NewBuiltin("finditer", finditer),
					"sub": starlark.NewBuiltin("sub", sub),
					// "subn":      starlark.NewBuiltin("subn", subn),
					// "escape":    starlark.NewBuiltin("escape", escape),
				},
			},
		}
	})
	return reModule, nil
}

// search(pattern,string,flags=0)
// Scan through string looking for the first location where the regular expression pattern produces a match,
// and return a corresponding match object. Return None if no position in the string matches the pattern;
// note that this is different from finding a zero-length match at some point in the string.
func search(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern, str starlark.String
		flags        starlark.Int
	)
	if err := starlark.UnpackArgs("search", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
		return starlark.None, err
	}
	r, err := newRegex(pattern)
	if err != nil {
		return starlark.None, err
	}
	matches := r.FindStringIndex(string(str))
	if len(matches) == 0 {
		return starlark.None, nil
	}

	vals := make([]interface{}, len(matches))
	for i, m := range matches {
		vals[i] = m
	}

	return util.Marshal(vals)
}

// match(pattern, string, flags=0)
// If zero or more characters at the beginning of string match the regular expression pattern,
// return a corresponding match object. Return None if the string does not match the pattern;
// note that this is different from a zero-length match.
// Note that even in MULTILINE mode, re.match() will only match at the beginning of the string and not at the beginning of each line.
// If you want to locate a match anywhere in string, use search() instead
func match(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern, str starlark.String
		flags        starlark.Int
	)
	if err := starlark.UnpackArgs("match", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
		return starlark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return starlark.None, err
	}

	vals := starlark.NewList(nil)
	for _, match := range r.FindAllStringSubmatch(string(str), -1) {
		if err = vals.Append(slStrSlice(match)); err != nil {
			return starlark.None, err
		}
	}

	return vals, nil
}

// fullmatch(pattern, string, flags=0)¶
// If the whole string matches the regular expression pattern, return a corresponding match object.
// Return None if the string does not match the pattern; note that this is different from a zero-length match.
// func fullmatch(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// 	var pattern starlark.String
// 	if err := starlark.UnpackArgs("fullmatch", args, kwargs, "pattern", &pattern); err != nil {
// 		return starlark.None, err
// 	}

// 	return starlark.None, nil
// }

// split(pattern, string, maxsplit=0, flags=0)
// Split string by the occurrences of pattern. If capturing parentheses are used in pattern,
// then the text of all groups in the pattern are also returned as part of the resulting list.
// If maxsplit is nonzero, at most maxsplit splits occur, and the remainder of the string
// is returned as the final element of the list.
func split(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern  starlark.String
		str      starlark.String
		maxSplit starlark.Int
		flags    starlark.Int
	)
	if err := starlark.UnpackArgs("split", args, kwargs, "pattern", &pattern, "string", &str, "maxsplit?", &maxSplit, "flags", &flags); err != nil {
		return starlark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return starlark.None, err
	}

	ms, _ := maxSplit.Int64()
	res := r.Split(string(str), int(ms))
	return slStrSlice(res), nil
}

// findall(pattern, string, flags=0)
// Returns all non-overlapping matches of pattern in string, as a list of strings.
// The string is scanned left-to-right, and matches are returned in the order found.
// If one or more groups are present in the pattern, return a list of groups;
// this will be a list of tuples if the pattern has more than one group.
// Empty matches are included in the result.
func findall(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern starlark.String
		str     starlark.String
		flags   starlark.Int
	)
	if err := starlark.UnpackArgs("findall", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
		return starlark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return starlark.None, err
	}
	res := r.FindAllString(string(str), -1)
	return slStrSlice(res), nil
}

// func finditer(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// 	var pattern starlark.String
// 	if err := starlark.UnpackArgs("finditer", args, kwargs, "pattern", &pattern); err != nil {
// 		return starlark.None, err
// 	}

// 	return starlark.None, nil
// }

// sub(pattern, repl, string, count=0, flags=0)
// Return the string obtained by replacing the leftmost non-overlapping occurrences of pattern
// in string by the replacement repl. If the pattern isn’t found, string is returned unchanged.
// repl can be a string or a function; if it is a string, any backslash escapes in it are processed.
// That is, \n is converted to a single newline character, \r is converted to a carriage return, and so forth.
// Unknown escapes such as \& are left alone. Backreferences, such as \6, are replaced with the substring matched by group 6 in the pattern.
func sub(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pattern, repl, str starlark.String
		count, flags       starlark.Int
	)
	if err := starlark.UnpackArgs("sub", args, kwargs, "pattern", &pattern, "repl", &repl, "string", &str, "count?", &count, "flags", &flags); err != nil {
		return starlark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return starlark.None, nil
	}
	res := r.ReplaceAllString(string(str), string(repl))
	return starlark.String(res), nil
}

// subn(pattern, repl, string, count=0, flags=0)
// Perform the same operation as sub(), but return a tuple (new_string, number_of_subs_made)
// func subn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// 	var pattern starlark.String
// 	if err := starlark.UnpackArgs("subn", args, kwargs, "pattern", &pattern); err != nil {
// 		return starlark.None, err
// 	}

// 	return starlark.None, nil
// }

// func escape(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// 	var pattern starlark.String
// 	if err := starlark.UnpackArgs("escape", args, kwargs, "pattern", &pattern); err != nil {
// 		return starlark.None, err
// 	}

// 	return starlark.None, nil
// }

func newRegex(pattern starlark.String) (*regexp.Regexp, error) {
	return regexp.Compile(string(pattern))
}

func slStrSlice(strs []string) starlark.Tuple {
	var vals starlark.Tuple
	for _, s := range strs {
		vals = append(vals, starlark.String(s))
	}
	return vals
}
