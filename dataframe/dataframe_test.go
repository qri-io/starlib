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

type invalidData struct{}

func TestDataframeFromRows(t *testing.T) {
	// Construct a valid dataframe from a row of various types of data
	rows := [][]interface{}{}
	record := []interface{}{"test", 31.2, 11.4, "ok", int64(597), "", 107, 6.91}
	rows = append(rows, record)
	df, err := BuildFromRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()
	expect := `        0     1     2   3    4  5    6    7
0    test  31.2  11.4  ok  597    107  6.9
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}

	// Ensure a dataframe with an invalid type ends up returning an error
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, &invalidData{}}
	rows = append(rows, record)
	_, err = BuildFromRows(rows)
	if err == nil {
		t.Fatal("expected to get an error, did not get one")
	}
	expectErr := `invalid object &{} of type *dataframe.invalidData`
	if expectErr != err.Error() {
		t.Errorf("error mismatch, expect: %s, got: %s", expectErr, err)
	}
}
