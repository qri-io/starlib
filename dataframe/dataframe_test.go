package dataframe

import (
	"testing"
)

func TestDataframeBasic(t *testing.T) {
	runTestScript(t, "testdata/dataframe_basic.star", "testdata/dataframe_basic.expect.txt")
}

func TestDataframeSize(t *testing.T) {
	runTestScript(t, "testdata/dataframe_size.star", "testdata/dataframe_size.expect.txt")
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

func TestDataframeAppend(t *testing.T) {
	runTestScript(t, "testdata/dataframe_append.star", "testdata/dataframe_append.expect.txt")
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

func TestDataframeSetCSV(t *testing.T) {
	runTestScript(t, "testdata/dataframe_set_csv.star", "testdata/dataframe_set_csv.expect.txt")
}

type invalidData struct{}

func TestDataframeFromRows(t *testing.T) {
	// Construct a valid dataframe from a single row of various types of data
	rows := [][]interface{}{}
	record := []interface{}{"test", 31.2, 11.4, "ok", int64(597), "", 107, 6.91}
	rows = append(rows, record)
	df, err := NewDataFrame(rows, nil, nil)
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
	_, err = NewDataFrame(rows, nil, nil)
	if err == nil {
		t.Fatal("expected to get an error, did not get one")
	}
	expectErr := `invalid object &{} of type *dataframe.invalidData`
	if expectErr != err.Error() {
		t.Errorf("error mismatch, expect: %s, got: %s", expectErr, err)
	}

	// Construct a dataframe from multiple rows
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, 17, int64(45)}
	rows = append(rows, record)
	record = []interface{}{"more", 9.8, 62, int64(3)}
	rows = append(rows, record)
	df, err = NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got = df.String()
	expect = `        0     1   2   3
0    test  31.2  17  45
1    more   9.8  62   3
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}

	// Construct a dataframe with non-matching columns, they get casted correctly
	rows = [][]interface{}{}
	record = []interface{}{"test", 31.2, 17, int64(45)}
	rows = append(rows, record)
	record = []interface{}{25, "ok", int64(4), "hi"}
	rows = append(rows, record)
	df, err = NewDataFrame(rows, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got = df.String()
	expect = `        0     1   2   3
0    test  31.2  17  45
1      25    ok   4  hi
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}

func TestDataframeFromList(t *testing.T) {
	ls := []interface{}{1.2, 3.4, 5.6}
	df, err := NewDataFrame(ls, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()
	expect := `       0
0    1.2
1    3.4
2    5.6
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}

func TestDataframeFromSeries(t *testing.T) {
	s := newSeriesFromStrings([]string{"a", "b", "c"}, nil, "")
	df, err := NewDataFrame(s, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := df.String()

	expect := `     0
0    a
1    b
2    c
`
	if got != expect {
		t.Errorf("mismatch: expect %q, got %q", expect, got)
	}
}
