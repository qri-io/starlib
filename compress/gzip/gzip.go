package gzip

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('compress/gzip.star', 'gzip')
const ModuleName = "compress/gzip.star"

// Module defines the gzip starlark api
var Module = &starlarkstruct.Module{
	Name: "gzip",
	Members: starlark.StringDict{
		"decompress": starlark.NewBuiltin("decompress", decompress),
	},
}

func decompress(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var compressed starlark.Value
	if err := starlark.UnpackPositionalArgs("decompress", args, kwargs, 1, &compressed); err != nil {
		return starlark.None, err
	}

	var rdr io.Reader
	switch v := compressed.(type) {
	case starlark.Bytes:
		rdr = strings.NewReader(string(v))
	case starlark.String:
		rdr = strings.NewReader(string(v))
	}

	r, err := gzip.NewReader(rdr)
	if err != nil {
		return starlark.None, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return starlark.None, err
	}
	return starlark.Bytes(data), nil
}
