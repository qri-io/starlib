package yaml

import (
	"bytes"
	"sync"

	"github.com/qri-io/starlib/util"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"gopkg.in/yaml.v2"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('encoding/yaml.star', 'yaml')
const ModuleName = "encoding/yaml.star"

var (
	once       sync.Once
	yamlModule starlark.StringDict
)

// LoadModule loads the base64 module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		yamlModule = starlark.StringDict{
			"yaml": &starlarkstruct.Module{
				Name: "yaml",
				Members: starlark.StringDict{
					"loads": starlark.NewBuiltin("loads", Loads),
					"dumps": starlark.NewBuiltin("dumps", Dumps),
				},
			},
		}
	})
	return yamlModule, nil
}

// Loads gets all values from a yaml source
func Loads(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		source starlark.String
		val    interface{}
	)

	err := starlark.UnpackArgs("loads", args, kwargs, "source", &source)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal([]byte(source.GoString()), &val); err != nil {
		return starlark.None, err
	}

	return util.Marshal(val)
}

// Dumps serializes a starlark object to a yaml string
func Dumps(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	values := make([]starlark.Value, len(args))
	valuePtrs := make([]interface{}, len(args))

	for idx, _ := range values {
		valuePtrs[idx] = &values[idx]
	}

	err := starlark.UnpackPositionalArgs("dumps", args, kwargs, 1, valuePtrs...)
	if err != nil {
		return starlark.None, err
	}

	buffer := new(bytes.Buffer)
	for idx, value := range values {

		goValue, err := util.Unmarshal(value)
		if err != nil {
			return starlark.None, err
		}

		rawBytes, err := yaml.Marshal(goValue)
		if err != nil {
			return starlark.None, err
		}

		buffer.Write(rawBytes)

		// Add the YAML document separator if, and only if,
		// there are more values to be serialized.
		if idx < len(values)-1 {
			buffer.WriteString("---\n")
		}
	}

	return starlark.String(buffer.String()), nil
}
