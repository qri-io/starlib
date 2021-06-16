package dataframe

import (
	"testing"
)

func TestSeriesBasic(t *testing.T) {
	runTestScript(t, "testdata/series_basic.star", "testdata/series_basic.expect.txt")
}

func TestSeriesAttrs(t *testing.T) {
	runTestScript(t, "testdata/series_attrs.star", "testdata/series_attrs.expect.txt")
}

func TestSeriesGet(t *testing.T) {
	runTestScript(t, "testdata/series_get.star", "testdata/series_get.expect.txt")
}

func TestSeriesPrint(t *testing.T) {
	runTestScript(t, "testdata/series_print.star", "testdata/series_print.expect.txt")
}

func TestSeriesIndexWithName(t *testing.T) {
	runTestScript(t, "testdata/series_index_name.star", "testdata/series_index_name.expect.txt")
}
