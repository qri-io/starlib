package dataframe

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func loadModule(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	switch module {
	case "dataframe.star":
		return starlark.StringDict{"dataframe": Module}, nil
	case "time.star":
		return starlark.StringDict{"time": time.Module}, nil
	case "assert.star":
		starlarktest.DataFile = func(pkgdir, filename string) string {
			_, currFileName, _, ok := runtime.Caller(1)
			if !ok {
				return ""
			}
			return filepath.Join(filepath.Dir(currFileName), filename)
		}
		return starlarktest.LoadAssertModule()
	}
	return nil, fmt.Errorf("invalid module")
}

func runScript(t *testing.T, scriptFilename string) (string, error) {
	t.Helper()
	output := "\n"
	printCollect := func(thread *starlark.Thread, msg string) {
		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	thread := &starlark.Thread{Load: loadModule}
	thread.Print = printCollect
	starlarktest.SetReporter(thread, t)
	thread.SetLocal(keyOutputConfig, &OutputConfig{})

	_, err := starlark.ExecFile(thread, scriptFilename, nil, nil)
	return strings.Trim(output, "\n"), err
}

func expectScriptOutput(t *testing.T, scriptFilename, expectFilename string) {
	t.Helper()
	output, err := runScript(t, scriptFilename)
	if err != nil {
		t.Fatal(err)
	}
	expect := mustReadFile(t, expectFilename)

	expect = strings.Trim(expect, "\n")
	output = strings.Trim(output, "\n")
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("mismatch. (-want +got):\n%s", diff)
	}
}

func mustReadFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
