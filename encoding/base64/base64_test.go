package base64

import (
	"fmt"
	"testing"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestFile(t *testing.T) {
	resolve.AllowFloat = true
	thread := &starlark.Thread{Load: newLoader()}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.star", nil, nil)
	if err != nil {
		if ee, ok := err.(*starlark.EvalError); ok {
			t.Error(ee.Backtrace())
		} else {
			t.Error(err)
		}
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
