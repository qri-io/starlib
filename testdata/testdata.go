package testdata

import (
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
	"path/filepath"
	"runtime"
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
