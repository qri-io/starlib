package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	starlark "github.com/google/skylark"
	starlarktest "github.com/google/skylark/skylarktest"
	dataset "github.com/qri-io/dataset"
)

func TestNewModule(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", "Mon, 01 Jun 2000 00:00:00 GMT")
		w.Write([]byte(`{"hello":"world"}`))
	}))
	starlark.Universe["test_server_url"] = starlark.String(ts.URL)

	ds := &dataset.Dataset{
		Transform: &dataset.Transform{
			Syntax: "starlark",
			Config: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	thread := &starlark.Thread{Load: newLoader(ds)}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.sky", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader(ds *dataset.Dataset) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		switch module {
		case ModuleName:
			return LoadModule()
		case "assert.sky":
			return starlarktest.LoadAssertModule()
		}

		return nil, fmt.Errorf("invalid module")
	}
}
