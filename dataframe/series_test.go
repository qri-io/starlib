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

func TestSeriesBasic(t *testing.T) {
	output := "\n"
	printCollect := func(thraed *starlark.Thread, msg string) {
		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	thread := &starlark.Thread{Load: testdata.NewModuleLoader(Module)}
	thread.Print = printCollect
	starlarktest.SetReporter(thread, t)

	_, err := starlark.ExecFile(thread, "testdata/series_basic.star", nil, nil)
	if err != nil {
		t.Error(err)
	}
	expect := MustReadFile(t, "testdata/series_basic.expect.txt")

	expect = strings.Trim(expect, "\n")
	output = strings.Trim(output, "\n")
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("mismatch. (-want +got):\n%s", diff)
	}
}

func TestSeriesGet(t *testing.T) {
	output := "\n"
	printCollect := func(thraed *starlark.Thread, msg string) {
		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	thread := &starlark.Thread{Load: testdata.NewModuleLoader(Module)}
	thread.Print = printCollect
	starlarktest.SetReporter(thread, t)

	_, err := starlark.ExecFile(thread, "testdata/series_get.star", nil, nil)
	if err != nil {
		t.Error(err)
	}
	expect := MustReadFile(t, "testdata/series_get.expect.txt")

	expect = strings.Trim(expect, "\n")
	output = strings.Trim(output, "\n")
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("mismatch. (-want +got):\n%s", diff)
	}
}

func TestSeriesPrint(t *testing.T) {
	output := "\n"
	printCollect := func(thraed *starlark.Thread, msg string) {
		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	thread := &starlark.Thread{Load: testdata.NewModuleLoader(Module)}
	thread.Print = printCollect
	starlarktest.SetReporter(thread, t)

	_, err := starlark.ExecFile(thread, "testdata/series_print.star", nil, nil)
	if err != nil {
		t.Error(err)
	}
	expect := MustReadFile(t, "testdata/series_print.expect.txt")

	expect = strings.Trim(expect, "\n")
	output = strings.Trim(output, "\n")
	if diff := cmp.Diff(expect, output); diff != "" {
		t.Errorf("mismatch. (-want +got):\n%s", diff)
	}
}

func TestSeriesRepeatScalar(t *testing.T) {

}

func MustReadFile(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
