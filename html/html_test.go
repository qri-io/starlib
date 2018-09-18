package html

import (
	"fmt"
	"testing"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarktest"
)

func TestModule(t *testing.T) {
	thread := &skylark.Thread{Load: newLoader()}
	skylarktest.SetReporter(thread, t)

	// Execute test file
	_, err := skylark.ExecFile(thread, "testdata/test.sky", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader() func(thread *skylark.Thread, module string) (skylark.StringDict, error) {
	return func(thread *skylark.Thread, module string) (skylark.StringDict, error) {
		switch module {
		case ModuleName:
			return skylark.StringDict{"html": skylark.NewBuiltin("html", NewDocument)}, nil
		case "assert.sky":
			return skylarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
