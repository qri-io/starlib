package dataframe

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIndexStringify(t *testing.T) {
	idx := NewIndex([]string{"cat", "dog", "eel"}, "animals")
	actual := idx.String()
	expect := `Index(['cat', 'dog', 'eel'], dtype='object', name='animals')`
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("error mismatch (-want +got):%s\n", diff)
	}
}
