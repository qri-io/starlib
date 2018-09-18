package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarktest"
	"github.com/qri-io/dataset"
)

func TestNewModule(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", "Mon, 01 Jun 2000 00:00:00 GMT")
		w.Write([]byte(`{"hello":"world"}`))
	}))
	skylark.Universe["test_server_url"] = skylark.String(ts.URL)

	ds := &dataset.Dataset{
		Transform: &dataset.Transform{
			Syntax: "skylark",
			Config: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	thread := &skylark.Thread{Load: newLoader(ds)}
	skylarktest.SetReporter(thread, t)

	// Execute test file
	_, err := skylark.ExecFile(thread, "testdata/test.sky", nil, nil)
	if err != nil {
		t.Error(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func newLoader(ds *dataset.Dataset) func(thread *skylark.Thread, module string) (skylark.StringDict, error) {
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
