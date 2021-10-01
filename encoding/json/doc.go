/*Package json defines utilities for converting Starlark values to/from JSON
strings. This package exists to add documentation only. The API is locked to
strictly match the Starlark module.
Users are encouraged to import the json package directly via:
go.starlark.net/starlarkjson

For source code see
https://github.com/google/starlark-go/tree/master/lib/starlarkjson

outline: json
  json provides functions for working with json data
  path: encoding/json
  functions:
    encode(obj) string
      Return a JSON string representation of a Starlark data structure
			params:
				 obj any
					obj is a valid Starlark data structure
      examples:
        encode object
          encode a simple object as a JSON string
          code:
            load("encoding/json.star", "json")
            x = json.encode({"foo": ["bar", "baz"]})
            print(x)
            # Output: {"foo":["bar","baz"]}
    decode(src) obj
      Return the Starlark representation of a string instance containing a JSON document. Decoding fails if src is not a valid JSON string.
			params:
				src string
					source string, must be valid JSON string
      examples:
        decode JSON string
          decode a JSON string into a Starlark structure
          code:
            load("encoding/json.star", "json")
            x = json.decode('{"foo": ["bar", "baz"]}')
    indent(src, prefix="", indent="\t") string
      The indent function pretty-prints a valid JSON encoding, and returns a string containing the indented form. It accepts one required positional parameter, the JSON string, and two optional keyword-only string parameters, prefix and indent, that specify a prefix of each new line, and the unit of indentation.
			params:
				src string
					source JSON string to encode
				prefix string
					optional. string prefix that will be prepended to each line. default is ""
				indent string
					optional. string that will be used to represent indentations. default is "\t"
      examples:
				basic
          "pretty print" a valid JSON encoding
          code:
            load("encoding/json.star", "json")
            x = json.indent('{"foo": ["bar", "baz"]}')
						# print(x)
						# {
            #    "foo": [
            #      "bar",
            #      "baz"
            #    ]
            # }
				using prefix & indent
					"pretty print" a valid JSON encoding, including optional prefix and indent parameters
					code:
						load("encoding/json.star", "json")
            x = json.indent('{"foo": ["bar", "baz"]}', prefix='....', indent="____")
						# print(x)
						# {
						# ....____"foo": [
						# ....________"bar",
						# ....________"baz"
						# ....____]
						# ....}
*/
package json

import "go.starlark.net/lib/json"

// ModuleName declares the intended load import string
// eg: load("encoding/json.star", "json")
const ModuleName = "encoding/json.star"

// Module exposes the starlarkjson module. Implementation located at
// https://github.com/google/starlark-go/tree/master/starlarkjson
var Module = json.Module
