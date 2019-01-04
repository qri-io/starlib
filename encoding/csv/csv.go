package csv

import (
	"encoding/csv"
	"io"
	"strings"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('csv.star', 'csv')
const ModuleName = "encoding/csv.star"

var (
	once      sync.Once
	csvModule starlark.StringDict
)

// LoadModule loads the base64 module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		csvModule = starlark.StringDict{
			"csv": starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
				"read_all": starlark.NewBuiltin("read_all", ReadAll),
			}),
		}
	})
	return csvModule, nil
}

// Module joins http tools to a dataset, allowing dataset
// to follow along with http requests
type Module struct {
}

// Struct returns this module's methods as a starlark Struct
func (m *Module) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, m.StringDict())
}

// StringDict returns all module methods in a starlark.StringDict
func (m *Module) StringDict() starlark.StringDict {
	return starlark.StringDict{
		"read_all": starlark.NewBuiltin("read_all", ReadAll),
	}
}

// ReadAll gets all values from a csv source
func ReadAll(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source starlark.Value
	if err := starlark.UnpackPositionalArgs("read_all", args, kwargs, 1, &source); err != nil {
		return nil, err
	}

	var r io.Reader
	switch source.Type() {
	case "string":
		str := string(source.(starlark.String))
		r = strings.NewReader(str)
	}
	csvr := csv.NewReader(r)

	strs, err := csvr.ReadAll()
	if err != nil {
		return starlark.None, err
	}

	vals := make([]starlark.Value, len(strs))
	for i, rowStr := range strs {
		row := make([]starlark.Value, len(rowStr))
		for j, cell := range rowStr {
			row[j] = starlark.String(cell)
		}
		vals[i] = starlark.NewList(row)
	}
	return starlark.NewList(vals), nil
}
