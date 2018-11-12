package xlsx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestFromURL(t *testing.T) {
	s := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	starlark.Universe["test_server_url"] = starlark.String(s.URL)

	thread := &starlark.Thread{Load: newLoader()}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.star", nil, nil)
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
