package dataframe

import (
	"testing"

	"github.com/qri-io/starlib/testdata"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestFile(t *testing.T) {
	thread := &starlark.Thread{Load: testdata.NewModuleLoader(Module)}
	starlarktest.SetReporter(thread, t)

	for _, filename := range []string{
		"testdata/dataframe.star",
		"testdata/index.star",
	} {
		t.Run(filename, func(t *testing.T) {
			_, err := starlark.ExecFile(thread, filename, nil, nil)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
