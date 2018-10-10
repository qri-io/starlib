package html

import (
	"fmt"
	"testing"

	starlark "github.com/google/skylark"
	starlarktest "github.com/google/skylark/skylarktest"
)

func TestModule(t *testing.T) {
	thread := &starlark.Thread{Load: newLoader()}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.sky", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader() func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		switch module {
		case ModuleName:
			return starlark.StringDict{"html": starlark.NewBuiltin("html", NewDocument)}, nil
		case "assert.sky":
			return starlarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
