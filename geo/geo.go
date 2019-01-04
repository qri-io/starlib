package geo

import (
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('geo.star', 'geo')
const ModuleName = "geo.star"

var (
	once      sync.Once
	geoModule starlark.StringDict
)

// LoadModule loads the geo module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		geoModule = starlark.StringDict{
			"geo": starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{}),
		}
	})
	return geoModule, nil
}
