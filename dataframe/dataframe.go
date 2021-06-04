package dataframe

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Name of the module
const Name = "dataframe"
// ModuleName is the filename of this module for the loader
const ModuleName = "dataframe.star"

// Module exposes the dataframe module
var Module = &starlarkstruct.Module{
	Name: Name,
	Members: starlark.StringDict{
		"Series": starlark.NewBuiltin("Series", newSeries),
	},
}
