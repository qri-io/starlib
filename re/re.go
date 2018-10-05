/*
Package re defines regular expression functions, it's intented to be a drop-in
subset of python's re module for starlark:
https://docs.python.org/3/library/re.html

	TODO - compile

	TODO - search
	match
	TODO - fullmatch
	split
	TODO - findall
	TODO - finditer
	sub
	TODO - subn
	TODO - escape
*/
package re

import (
	"regexp"
	"sync"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('re.sky', 're')
const ModuleName = "re.sky"

var (
	once     sync.Once
	reModule skylark.StringDict
)

// LoadModule loads the re module.
// It is concurrency-safe and idempotent.
func LoadModule() (skylark.StringDict, error) {
	once.Do(func() {
		reModule = skylark.StringDict{
			"re": skylarkstruct.FromStringDict(skylark.String("re"), skylark.StringDict{
				// TODO
				// "compile" : skylark.NewBuiltin("complile", compile),

				// "search":    skylark.NewBuiltin("search", search),
				"match": skylark.NewBuiltin("match", match),
				// "fullmatch": skylark.NewBuiltin("fullmatch", fullmatch),
				"split":   skylark.NewBuiltin("split", split),
				"findall": skylark.NewBuiltin("findall", findall),
				// "finditer":  skylark.NewBuiltin("finditer", finditer),
				"sub": skylark.NewBuiltin("sub", sub),
				// "subn":      skylark.NewBuiltin("subn", subn),
				// "escape":    skylark.NewBuiltin("escape", escape),
			}),
		}
	})
	return reModule, nil
}

// search(pattern,string,flags=0)
// Scan through string looking for the first location where the regular expression pattern produces a match,
// and return a corresponding match object. Return None if no position in the string matches the pattern;
// note that this is different from finding a zero-length match at some point in the string.
// func search(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// 	var (
// 		pattern, str skylark.String
// 		flags skylark.Int
// 	)
// 	if err := skylark.UnpackArgs("search", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
// 		return skylark.None, err
// 	}
// 	r, err := newRegex(pattern)
// 	if err != nil {
// 		return skylark.None, err
// 	}
// 	// r.FindStringIndex(string(str))
// 	return skylark.None, nil
// }

// match(pattern, string, flags=0)
// If zero or more characters at the beginning of string match the regular expression pattern,
// return a corresponding match object. Return None if the string does not match the pattern;
// note that this is different from a zero-length match.
// Note that even in MULTILINE mode, re.match() will only match at the beginning of the string and not at the beginning of each line.
// If you want to locate a match anywhere in string, use search() instead
func match(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var (
		pattern, str skylark.String
		flags        skylark.Int
	)
	if err := skylark.UnpackArgs("match", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
		return skylark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return skylark.None, err
	}

	vals := skylark.NewList(nil)
	for _, match := range r.FindAllStringSubmatch(string(str), -1) {
		if err = vals.Append(slStrSlice(match)); err != nil {
			return skylark.None, err
		}
	}

	return vals, nil
}

// fullmatch(pattern, string, flags=0)¶
// If the whole string matches the regular expression pattern, return a corresponding match object.
// Return None if the string does not match the pattern; note that this is different from a zero-length match.
// func fullmatch(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// 	var pattern skylark.String
// 	if err := skylark.UnpackArgs("fullmatch", args, kwargs, "pattern", &pattern); err != nil {
// 		return skylark.None, err
// 	}

// 	return skylark.None, nil
// }

// split(pattern, string, maxsplit=0, flags=0)
// Split string by the occurrences of pattern. If capturing parentheses are used in pattern,
// then the text of all groups in the pattern are also returned as part of the resulting list.
// If maxsplit is nonzero, at most maxsplit splits occur, and the remainder of the string
// is returned as the final element of the list.
func split(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var (
		pattern  skylark.String
		str      skylark.String
		maxSplit skylark.Int
		flags    skylark.Int
	)
	if err := skylark.UnpackArgs("split", args, kwargs, "pattern", &pattern, "string", &str, "maxsplit?", &maxSplit, "flags", &flags); err != nil {
		return skylark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return skylark.None, err
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
func findall(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var (
		pattern skylark.String
		str     skylark.String
		flags   skylark.Int
	)
	if err := skylark.UnpackArgs("findall", args, kwargs, "pattern", &pattern, "string", &str, "flags?", &flags); err != nil {
		return skylark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return skylark.None, err
	}
	res := r.FindAllString(string(str), -1)
	return slStrSlice(res), nil
}

// func finditer(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// 	var pattern skylark.String
// 	if err := skylark.UnpackArgs("finditer", args, kwargs, "pattern", &pattern); err != nil {
// 		return skylark.None, err
// 	}

// 	return skylark.None, nil
// }

// sub(pattern, repl, string, count=0, flags=0)
// Return the string obtained by replacing the leftmost non-overlapping occurrences of pattern
// in string by the replacement repl. If the pattern isn’t found, string is returned unchanged.
// repl can be a string or a function; if it is a string, any backslash escapes in it are processed.
// That is, \n is converted to a single newline character, \r is converted to a carriage return, and so forth.
// Unknown escapes such as \& are left alone. Backreferences, such as \6, are replaced with the substring matched by group 6 in the pattern.
func sub(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var (
		pattern, repl, str skylark.String
		count, flags       skylark.Int
	)
	if err := skylark.UnpackArgs("sub", args, kwargs, "pattern", &pattern, "repl", &repl, "string", &str, "count?", &count, "flags", &flags); err != nil {
		return skylark.None, err
	}

	r, err := newRegex(pattern)
	if err != nil {
		return skylark.None, nil
	}
	res := r.ReplaceAllString(string(str), string(repl))
	return skylark.String(res), nil
}

// subn(pattern, repl, string, count=0, flags=0)
// Perform the same operation as sub(), but return a tuple (new_string, number_of_subs_made)
// func subn(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// 	var pattern skylark.String
// 	if err := skylark.UnpackArgs("subn", args, kwargs, "pattern", &pattern); err != nil {
// 		return skylark.None, err
// 	}

// 	return skylark.None, nil
// }

// func escape(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
// 	var pattern skylark.String
// 	if err := skylark.UnpackArgs("escape", args, kwargs, "pattern", &pattern); err != nil {
// 		return skylark.None, err
// 	}

// 	return skylark.None, nil
// }

func newRegex(pattern skylark.String) (*regexp.Regexp, error) {
	return regexp.Compile(string(pattern))
}

func slStrSlice(strs []string) skylark.Tuple {
	var vals skylark.Tuple
	for _, s := range strs {
		vals = append(vals, skylark.String(s))
	}
	return vals
}
