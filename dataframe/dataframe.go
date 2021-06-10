package dataframe

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// Name of the module
	Name = "dataframe"
	// ModuleName is the filename of this module for the loader
	ModuleName = "dataframe.star"
)

// Module exposes the dataframe module
var Module = &starlarkstruct.Module{
	Name: Name,
	Members: starlark.StringDict{
		"Series": starlark.NewBuiltin("Series", newSeries),
	},
}
