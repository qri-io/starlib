package zipfile

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarktest"
)

func TestFile(t *testing.T) {
	thread := &skylark.Thread{Load: newLoader()}
	skylarktest.SetReporter(thread, t)

	zipBytes, err := ioutil.ReadFile("testdata/hello_world.zip")
	if err != nil {
		t.Fatal(err)
	}

	// Execute test file
	_, err = skylark.ExecFile(thread, "testdata/test.sky", nil, skylark.StringDict{
		"hello_world_zip": skylark.String(zipBytes),
	})
	if err != nil {
		t.Error(err)
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
