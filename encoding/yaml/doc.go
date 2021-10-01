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
        examples:
          basic
            encode to yaml
            code:
              load("encoding/yaml.star", "yaml")
              data = {"foo": "bar", "baz": True}
              res = yaml.dumps(data)
      loads(source) object
        read a source yaml string to a starlark object
        params:
          source string
            input string of yaml data
        examples:
          basic
            load a yaml string
            code:
              load("encoding/yaml.star", "yaml")
              data = """foo: bar
              baz: true
              """
              d = yaml.loads(data)
              print(d)
              # Output: {"foo": "bar", "baz": True}

*/
package yaml
