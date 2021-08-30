package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
)

// AtIndexer is returned by dataframe.At and can index to any cell
type AtIndexer struct {
	owner *DataFrame
}

// compile-time interface assertions
var (
	_ starlark.Value     = (*AtIndexer)(nil)
	_ starlark.Mapping   = (*AtIndexer)(nil)
	_ starlark.HasSetKey = (*AtIndexer)(nil)
)

// NewAtIndexer returns a new AtIndexer with the text values and name
func NewAtIndexer(owner *DataFrame) *AtIndexer {
	return &AtIndexer{owner: owner}
}

// Freeze prevents the index from being mutated
func (ai *AtIndexer) Freeze() {
}

// Hash cannot be used with AtIndexer
func (ai *AtIndexer) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", ai.Type())
}

// String returns the index as a string
func (ai *AtIndexer) String() string {
	return fmt.Sprintf("AtIndexer()")
}

// Truth converts the index into a bool
func (ai *AtIndexer) Truth() starlark.Bool {
	return true
}

// Type returns the type as a string
func (ai *AtIndexer) Type() string {
	return fmt.Sprintf("%s.AtIndexer", Name)
}

// Get returns the value at the given key, which must be an int pair
// implements the Mapping interface
func (ai *AtIndexer) Get(key starlark.Value) (starlark.Value, bool, error) {
	keyOne, keyTwo, ok := keyToIntPair(key)
	if !ok {
		return nil, false, fmt.Errorf("invalid Index: %v", key)
	}
	series := ai.owner.body[keyTwo]
	cell := series.At(keyOne)
	val, err := convertToStarlark(cell)
	if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

// SetKey sets a value at the given key, which must be an int pair
// implements the HasSetKey interface
func (ai *AtIndexer) SetKey(key, val starlark.Value) error {
	keyOne, keyTwo, ok := keyToIntPair(key)
	if !ok {
		return fmt.Errorf("invalid Index: %v", key)
	}
	series := ai.owner.body[keyTwo]
	// Convert to a go native type
	item := toNativeValue(val)
	series.SetAt(keyOne, item)
	return nil
}

func keyToIntPair(key starlark.Value) (int, int, bool) {
	tup, ok := key.(starlark.Tuple)
	if !ok {
		return 0, 0, false
	}
	if len(tup) != 2 {
		return 0, 0, false
	}
	first := tup.Index(0)
	second := tup.Index(1)
	keyOne, err := starlark.AsInt32(first)
	if err != nil {
		return 0, 0, false
	}
	keyTwo, err := starlark.AsInt32(second)
	if err != nil {
		return 0, 0, false
	}
	return keyOne, keyTwo, true
}
