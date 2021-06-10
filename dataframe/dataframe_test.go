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
	// NOTE: Different from python pandas
/*
0    foo:1
1    bar:2
2    baz:3
3    foo:5
dtype: object

  lkey  value combined
0  foo      1    foo:1
1  bar      2    bar:2
2  baz      3    baz:3
3  foo      5    foo:5

*/
}

func TestDataframeDropDuplicates(t *testing.T) {
	runTestScript(t, "testdata/dataframe_drop_duplicates.star",
		"testdata/dataframe_drop_duplicates.expect.txt")
}

func TestDataframeMerge(t *testing.T) {
	runTestScript(t, "testdata/dataframe_merge.star", "testdata/dataframe_merge.expect.txt")
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
