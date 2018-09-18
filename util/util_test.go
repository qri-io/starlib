package util

import (
	"testing"

	"github.com/google/skylark"
)

func TestAsString(t *testing.T) {
	cases := []struct {
		in       skylark.Value
		got, err string
	}{
		{skylark.String("foo"), "foo", ""},
		{skylark.String("\"foo'"), "\"foo'", ""},
		{skylark.Bool(true), "", "invalid syntax"},
	}

	for i, c := range cases {
		got, err := AsString(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if c.got != got {
			t.Errorf("case %d. expected: '%s', got: '%s'", i, c.got, got)
		}
	}
}
