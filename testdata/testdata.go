package testdata

import (
	"fmt"
	"path/filepath"
	"runtime"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/starlarktest"
)

// Newloader implements the 'load' operation as used in the evaluator tests.
// takes a LoadModule function
// a ModuleName
// and the relative path to the testdata
func NewLoader(loader func() (starlark.StringDict, error), moduleName string) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		switch module {
		case moduleName:
			return loader()
		case "assert.star":
			starlarktest.DataFile = func(pkgdir, filename string) string {
				_, currFileName, _, ok := runtime.Caller(1)
				if !ok {
					return ""
				}
				return filepath.Join(filepath.Dir(currFileName), filename)
			}
			return starlarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}

// NewModuleLoader constructs a loader from a given set of modules
func NewModuleLoader(modules ...*starlarkstruct.Module) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		if module == "assert.star" {
			starlarktest.DataFile = func(pkgdir, filename string) string {
				_, currFileName, _, ok := runtime.Caller(1)
				if !ok {
					return ""
				}
				return filepath.Join(filepath.Dir(currFileName), filename)
			}
			return starlarktest.LoadAssertModule()
		}

		for _, mod := range modules {
			if module == mod.Name+".star" {
				return starlark.StringDict{mod.Name: mod}, nil
			}
		}

		return nil, fmt.Errorf("invalid module")
	}
}
