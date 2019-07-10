package yaml

import (
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
	var (
		source starlark.Value
	)

	err := starlark.UnpackArgs("dumps", args, kwargs, "source", &source)
	if err != nil {
		return starlark.None, err
	}

	val, err := util.Unmarshal(source)
	if err != nil {
		return starlark.None, err
	}

	data, err := yaml.Marshal(val)
	if err != nil {
		return starlark.None, err
	}

	return starlark.String(string(data)), nil
}
