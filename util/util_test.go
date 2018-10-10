package util

import (
	"testing"

	starlark "github.com/google/skylark"
)

func TestAsString(t *testing.T) {
	cases := []struct {
		in       starlark.Value
		got, err string
	}{
		{starlark.String("foo"), "foo", ""},
		{starlark.String("\"foo'"), "\"foo'", ""},
		{starlark.Bool(true), "", "invalid syntax"},
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
