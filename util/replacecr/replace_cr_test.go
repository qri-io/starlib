package replacecr

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	input := []byte("foo\r\rbar\r\nbaz\r\r")
	expect := []byte("foo\r\n\r\nbar\r\nbaz\r\n\r\n")

	got := make([]byte, 19)
	n, err := Reader(bytes.NewReader(input)).Read(got)
	if err != nil && err.Error() != "EOF" {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 19 {
		t.Errorf("length error. expected: %d, got: %d", 19, n)
	}
	if !bytes.Equal(expect, got) {
		t.Errorf("byte mismatch. expected:\n%v\ngot:\n%v", expect, got)
	}
}
