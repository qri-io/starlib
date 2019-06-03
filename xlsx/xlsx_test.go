package xlsx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qri-io/starlib/testdata"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestFromURL(t *testing.T) {
	s := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	starlark.Universe["test_server_url"] = starlark.String(s.URL)

	thread := &starlark.Thread{Load: testdata.NewLoader(LoadModule, ModuleName)}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.star", nil, nil)
	if err != nil {
		t.Error(err)
	}
}
