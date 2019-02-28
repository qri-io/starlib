/*Package json implements a json parser for the starlark programming dialect

  outline: json
    json provides functions for working with json data
    path: encoding/json
    functions:
      dumps(obj) string
        serialize obj to a JSON string
        params:
          obj object
            input object
      loads(source) object
        read a source JSON string to a starlark object
        params:
          source string
            input string of json data

*/
package json
