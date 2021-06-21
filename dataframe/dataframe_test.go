package dataframe

import (
	"testing"
)

func TestDataframeBasic(t *testing.T) {
	runTestScript(t, "testdata/dataframe_basic.star", "testdata/dataframe_basic.expect.txt")
}

func TestDataframeSetKey(t *testing.T) {
	runTestScript(t, "testdata/dataframe_setkey.star", "testdata/dataframe_setkey.expect.txt")
}

func TestDataframeGetKey(t *testing.T) {
	runTestScript(t, "testdata/dataframe_getkey.star", "testdata/dataframe_getkey.expect.txt")
}

func TestDataframeApply(t *testing.T) {
	runTestScript(t, "testdata/dataframe_apply.star", "testdata/dataframe_apply.expect.txt")
}

func TestDataframeDropDuplicates(t *testing.T) {
	runTestScript(t, "testdata/dataframe_drop_duplicates.star",
		"testdata/dataframe_drop_duplicates.expect.txt")
}

func TestDataframeHead(t *testing.T) {
	runTestScript(t, "testdata/dataframe_head.star", "testdata/dataframe_head.expect.txt")
}

func TestDataframeMerge(t *testing.T) {
	runTestScript(t, "testdata/dataframe_merge.star", "testdata/dataframe_merge.expect.txt")
}

func TestDataframeReadCSV(t *testing.T) {
	runTestScript(t, "testdata/dataframe_read_csv.star", "testdata/dataframe_read_csv.expect.txt")
}

func TestDataframeResetIndex(t *testing.T) {
	runTestScript(t, "testdata/dataframe_reset_index.star", "testdata/dataframe_reset_index.expect.txt")
}

func TestDataframeColumns(t *testing.T) {
	runTestScript(t, "testdata/dataframe_columns.star", "testdata/dataframe_columns.expect.txt")
}

func TestDataframeGroupBy(t *testing.T) {
	runTestScript(t, "testdata/dataframe_groupby.star", "testdata/dataframe_groupby.expect.txt")
}
