package dataframe

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qri-io/starlib/testdata"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func runTestScript(t *testing.T, scriptFilename, expectFilename string) {
	output := "\n"
	printCollect := func(thread *starlark.Thread, msg string) {
		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	thread := &starlark.Thread{Load: testdata.NewModuleLoader(Module)}
	thread.Print = printCollect
	starlarktest.SetReporter(thread, t)

	_, err := starlark.ExecFile(thread, scriptFilename, nil, nil)
	if err != nil {
		t.Error(err)
	}
	expect := mustReadFile(t, expectFilename)

	expect = strings.Trim(expect, "\n")
	output = strings.Trim(output, "\n")
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("mismatch. (-want +got):\n%s", diff)
	}
}

func mustReadFile(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
