# re
re defines regular expression functions, it's intended to be a drop-in subset of python's re module for starlark: https://docs.python.org/3/library/re.html

## Functions

#### `findall(pattern, text, flags=0)`
Returns all non-overlapping matches of pattern in string, as a list of strings. The string is scanned left-to-right, and matches are returned in the order found. If one or more groups are present in the pattern, return a list of groups; this will be a list of tuples if the pattern has more than one group. Empty matches are included in the result.

**parameters:**

| name | type | description |
|------|------|-------------|
| `pattern` | `string` | regular expression pattern string |
| `text` | `string` | string to find within |
| `flags` | `int` | integer flags to control regex behaviour. reserved for future use |


#### `split(pattern, text, maxsplit=0, flags=0)`
Split text by the occurrences of pattern. If capturing parentheses are used in pattern, then the text of all groups in the pattern are also returned as part of the resulting list. If maxsplit is nonzero, at most maxsplit splits occur, and the remainder of the string is returned as the final element of the list.

**parameters:**

| name | type | description |
|------|------|-------------|
| `pattern` | `string` | regular expression pattern string |
| `text` | `string` | input string to split |
| `maxsplit` | `int` | maximum number of splits to make. default 0 splits all matches |
| `flags` | `int` | integer flags to control regex behaviour. reserved for future use |


#### `sub(pattern, repl, text, count=0, flags=0)`
Return the string obtained by replacing the leftmost non-overlapping occurrences of pattern in string by the replacement repl. If the pattern isn’t found, string is returned unchanged. repl can be a string or a function; if it is a string, any backslash escapes in it are processed. That is, \n is converted to a single newline character, \r is converted to a carriage return, and so forth.

**parameters:**

| name | type | description |
|------|------|-------------|
| `pattern` | `string` | regular expression pattern string |
| `repl` | `string` | string to replace matches with |
| `text` | `string` | input string to replace |
| `count` | `int` | number of replacements to make, default 0 means replace all matches |
| `flags` | `int` | integer flags to control regex behaviour. reserved for future use |


