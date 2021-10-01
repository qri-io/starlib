/*Package regexp defines regular expression functions. This package is a shim
while we wait for code to land in starlark/lib: https://github.com/google/starlark-go/pull/369


  outline: regexp
    re defines regular expression functions
    path: re
    functions:
      compile(pattern) Regexp
        Compiles a pattern in RE2 syntax (https://github.com/google/re2/wiki/Syntax) to a value of type 'regexp'.
        Each call to compile returns a distinct regexp value.
        A regexp value can be used for matching using its matches, find, find_all and other methods.
        params:
          pattern string
            regular expression pattern string
        examples:
          basic
            construct a regexp that matches ab one or more times:
            code:
              load("regexp.star", "regexp")
              re_abs = regexp.compile("ab+")
              print(re_abs.find("abbbbb"))
              # Output: abbbbb
    types:
      Regexp
        a compiled regular expression
        functions:
          find(src)
            Returns a string holding the text of the leftmost match in the given string of the regular expression regexp.
            The result is "" if there is no match or if the pattern successfully matches an empty string.
          find_all(src, max)
            Returns a new, mutable list of all successive matches of the regular expression regexp. An empty list indicates no match.
            If max > 0, at most max strings are returned. If max == 0, an empty list is returned. If max < 0, all strings are returned.
            The parameter max is optional: by default no limit is applied.
          find_submatches(src)
            Returns a new, mutable list of strings holding the text of the leftmost match of the regular expression regexp in the given string and the matches, if any, of its subexpressions.
            An empty list indicates no match.
          matches(src)
            Reports whether the given string contains any match of the regular expression regexp.
          replace_all(src, repl)
            Returns a copy of the given string, replacing matches of the regular expression regexp with the replacement string repl. Inside repl, backslash-escaped digits (\1 to \9) can be used to insert text matching corresponding parenthesized group from the pattern.
            \0 in repl refers to the entire matching text.
          replace_all(src, replFunc)
            Returns a copy of the given string in which all matches of the regular expression regexp have been replaced by the return value of the replacement function applied to the matched substring.
          split(src, max)
            Returns a new, mutable list of strings between all the matches of the regular expression regexp.
            If max > 0, at most max strings are returned knowing that the last string is the unsplit remainder. If max == 0, an empty list is returned. If max < 0, all strings are returned.
            The parameter max is optional: by default no limit is applied.
*/
package regexp

// ModuleName declares the intended load import string
// eg: load("regexp.star", "regexp")
const ModuleName = "regexp.star"
