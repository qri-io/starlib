package math

import (
	"fmt"
	"testing"

	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
	"github.com/google/skylark/skylarktest"
)

func TestFile(t *testing.T) {
	resolve.AllowFloat = true
	thread := &skylark.Thread{Load: newLoader()}
	skylarktest.SetReporter(thread, t)

	// Execute test file
	_, err := skylark.ExecFile(thread, "testdata/test.sky", nil, nil)
	if err != nil {
		if ee, ok := err.(*skylark.EvalError); ok {
			t.Error(ee.Backtrace())
		} else {
			t.Error(err)
		}
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader() func(thread *skylark.Thread, module string) (skylark.StringDict, error) {
	return func(thread *skylark.Thread, module string) (skylark.StringDict, error) {
		switch module {
		case ModuleName:
			return LoadModule()
		case "assert.sky":
			return skylarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
