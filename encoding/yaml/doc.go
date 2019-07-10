/*Package yaml implements a yaml parser for the starlark programming dialect

  outline: yaml
    yaml provides functions for working with yaml data
    path: encoding/yaml
    functions:
      dumps(obj) string
        serialize obj to a yaml string
        params:
          obj object
            input object
      loads(source) object
        read a source yaml string to a starlark object
        params:
          source string
            input string of yaml data

*/
package yaml
