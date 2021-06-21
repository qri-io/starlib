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

func TestDataframeColumns(t *testing.T) {
	runTestScript(t, "testdata/dataframe_columns.star", "testdata/dataframe_columns.expect.txt")
}
