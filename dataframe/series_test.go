package dataframe

import (
	"testing"
)

func TestSeriesBasic(t *testing.T) {
	expectScriptOutput(t, "testdata/series_basic.star", "testdata/series_basic.expect.txt")
}

func TestSeriesAttrs(t *testing.T) {
	expectScriptOutput(t, "testdata/series_attrs.star", "testdata/series_attrs.expect.txt")
}

func TestSeriesGet(t *testing.T) {
	expectScriptOutput(t, "testdata/series_get.star", "testdata/series_get.expect.txt")
}

func TestSeriesPrint(t *testing.T) {
	expectScriptOutput(t, "testdata/series_print.star", "testdata/series_print.expect.txt")
}

func TestSeriesBoolSelect(t *testing.T) {
	expectScriptOutput(t, "testdata/series_bool_select.star",
		"testdata/series_bool_select.expect.txt")
}

func TestSeriesIndexWithName(t *testing.T) {
	expectScriptOutput(t, "testdata/series_index_name.star", "testdata/series_index_name.expect.txt")
}

func TestSeriesAsType(t *testing.T) {
	expectScriptOutput(t, "testdata/series_astype.star", "testdata/series_astype.expect.txt")
}

func TestSeriesNotNull(t *testing.T) {
	expectScriptOutput(t, "testdata/series_notnull.star", "testdata/series_notnull.expect.txt")
}

func TestSeriesUnique(t *testing.T) {
	expectScriptOutput(t, "testdata/series_unique.star", "testdata/series_unique.expect.txt")
}
