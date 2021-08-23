package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
)

// RangeIndex represents an Index as a numerical range.
type RangeIndex struct {
	start int
	stop  int
	step  int
}

// compile-time interface assertions
var (
	_ starlark.Value    = (*RangeIndex)(nil)
	_ starlark.HasAttrs = (*RangeIndex)(nil)
	_ starlark.Sequence = (*RangeIndex)(nil)
)

// NewRangeIndex returns a new RangeIndex of the given size
func NewRangeIndex(size int) *RangeIndex {
	return &RangeIndex{start: 0, stop: size, step: 1}
}

// Freeze prevents the rangeIndex from being mutated
func (ri *RangeIndex) Freeze() {
	// pass
}

// Hash cannot be used with RangeIndex
func (ri *RangeIndex) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", ri.Type())
}

// String returns the range index as a string
func (ri *RangeIndex) String() string {
	return fmt.Sprintf("RangeIndex(start=%d, stop=%d, step=%d)", ri.start, ri.stop, ri.step)
}

// Truth converts the range index into a bool
func (ri *RangeIndex) Truth() starlark.Bool {
	return true
}

// Type returns the type as a string
func (ri *RangeIndex) Type() string {
	return fmt.Sprintf("%s.RangeIndex", Name)
}

// Attr gets a value for a string attribute
func (ri *RangeIndex) Attr(name string) (starlark.Value, error) {
	return nil, starlark.NoSuchAttrError(name)
}

// AttrNames lists available dot expression strings
func (ri *RangeIndex) AttrNames() []string {
	return []string{}
}

// Iterate returns an iterator for the rangeIndex
func (ri *RangeIndex) Iterate() starlark.Iterator {
	return &rangeIndexIterator{count: ri.start, limit: ri.stop}
}

// Len returns the length of the index
func (ri *RangeIndex) Len() int {
	return ri.stop - ri.start
}

type rangeIndexIterator struct {
	count int
	limit int
}

// Done does cleanup work when iteration finishes, not needed
func (it *rangeIndexIterator) Done() {}

// Next assigns the next item and returns whether one was found
func (it *rangeIndexIterator) Next(p *starlark.Value) bool {
	if it.count < it.limit {
		*p = starlark.MakeInt(it.count)
		it.count++
		return true
	}
	return false
}
