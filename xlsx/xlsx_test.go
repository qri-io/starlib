package xlsx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarktest"
)

func TestFromURL(t *testing.T) {
	s := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	skylark.Universe["test_server_url"] = skylark.String(s.URL)

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
			return LoadModule()
		case "assert.sky":
			return skylarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
