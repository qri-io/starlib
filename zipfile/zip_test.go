package zipfile

import (
	"fmt"
	"io/ioutil"
	"testing"

	starlark "github.com/google/skylark"
	starlarktest "github.com/google/skylark/skylarktest"
)

func TestFile(t *testing.T) {
	thread := &starlark.Thread{Load: newLoader()}
	starlarktest.SetReporter(thread, t)

	zipBytes, err := ioutil.ReadFile("testdata/hello_world.zip")
	if err != nil {
		t.Fatal(err)
	}

	// Execute test file
	_, err = starlark.ExecFile(thread, "testdata/test.star", nil, starlark.StringDict{
		"hello_world_zip": starlark.String(zipBytes),
	})
	if err != nil {
		t.Error(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader() func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		switch module {
		case ModuleName:
			return LoadModule()
		case "assert.star":
			return starlarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
