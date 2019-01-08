package geo

import (
	"fmt"
	"sort"

	"github.com/paulmach/orb"

	"go.starlark.net/starlark"
)

func newPolygon(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		rings *starlark.List
	)

	if err = starlark.UnpackArgs("Polygon", args, kwargs, "rings", &rings); err != nil {
		return
	}

	i := 0
	p := make(Polygon, rings.Len())
	iter := rings.Iterate()
	defer iter.Done()
	var x starlark.Value

	for iter.Next(&x) {
		switch l := x.(type) {
		case Line:
			if !l.IsClosed() {
				err = fmt.Errorf("line %d is not closed", i)
				return
			}
			p[i] = l
		case *starlark.List:
			line, err := lineFromList(l)
			if err != nil {
				return v, err
			}
			if !line.IsClosed() {
				return v, fmt.Errorf("line %d is not closed", i)
			}
			p[i] = line
		default:
			err = fmt.Errorf("invalid type for creating polygon: %s", x.Type())
			return
		}
		i++
	}
	return p, nil
}

// Polygon is the starlark geographic polygon type
type Polygon []Line

// assert polygon is a starlark value
var _ starlark.Value = (*Polygon)(nil)

// String implements the starlark.Value interface
func (p Polygon) String() string {
	// TODO (b5): finish this
	return fmt.Sprintf("polygon %d rings", len(p))
	// w := strings.Builder{}
	// w.WriteRune('[')
	// for _, p := range p {
	// 	w.WriteString(p.String() + " ")
	// }
	// w.WriteRune(']')
	// return w.String()
}

// Type implements the starlark.Value interface
func (p Polygon) Type() string { return "Polygon" }

// Freeze implements the starlark.Value interface, polygon is immutable
func (p Polygon) Freeze() {}

// Truth implements the starlark.Value interface
func (p Polygon) Truth() starlark.Bool {
	return len(p) > 0
}

// Hash implements the starlark.Value interface
func (p Polygon) Hash() (h uint32, err error) {
	var hash uint32
	for _, l := range p {
		if hash, err = l.Hash(); err != nil {
			return
		}

		// TODO (b5): this is bad, unique values will report as non-unique. fix
		h += hash
	}
	return
}

// OrbPolygon formats polygon as an orb.Polygon
func (p Polygon) OrbPolygon() orb.Polygon {
	ring := make([]orb.Ring, len(p))
	for i, l := range p {
		ring[i] = l.OrbRing()
	}
	return orb.Polygon(ring)
}

type polygonIterator struct {
	idx     int
	polygon Polygon
}

// If the iterator is exhausted, Next returns false.
// Otherwise it sets *p to the current element of the sequence,
// advances the iterator, and returns true.
func (li *polygonIterator) Next(p *starlark.Value) bool {
	if li.idx == len(li.polygon) {
		return false
	}

	*p = li.polygon[li.idx]
	li.idx++
	return true
}

// Done is a no-op
// TODO (b5): this is supposed to be a no-op. right?
func (li *polygonIterator) Done() {}

// Iterate implements the starlark iterator interface
func (p Polygon) Iterate() starlark.Iterator {
	return &polygonIterator{polygon: p}
}

// Index implements the starlark indexable interface
func (p Polygon) Index(i int) starlark.Value {
	return p[i]
}

// Len implements the starlark indexable interface
func (p Polygon) Len() int {
	return len(p)
}

// func Slice(start, end, step int) starlark.Value {
// }

// Attr gets an attribute of a polygon
func (p Polygon) Attr(name string) (starlark.Value, error) {
	switch name {
	default:
		// attr does not exist
		return nil, nil
	}
}

// AttrNames returns all possible attribute names
func (p Polygon) AttrNames() []string {
	names := []string{
		// methods
	}
	sort.Strings(names)
	return names
}
