package geo

import (
	"fmt"
	"sort"

	"github.com/paulmach/orb"
	"go.starlark.net/starlark"
)

func newMultiPolygon(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		polygons *starlark.List
		x        starlark.Value
	)

	if err = starlark.UnpackArgs("MultiPolygon", args, kwargs, "polygons", &polygons); err != nil {
		return
	}

	i := 0
	p := make(MultiPolygon, polygons.Len())
	iter := polygons.Iterate()
	defer iter.Done()

	for iter.Next(&x) {
		switch geom := x.(type) {
		case Polygon:
			p[i] = geom
		default:
			err = fmt.Errorf("invalid type for creating MultiPolygon: %s", x.Type())
			return
		}
		i++
	}
	return p, nil
}

// MultiPolygon is the starlark geographic collection-of-polygons type
type MultiPolygon []Polygon

// assert polygon is a starlark value
var _ starlark.Value = (*MultiPolygon)(nil)

// String implements the starlark.Value interface
func (p MultiPolygon) String() string {
	// TODO (b5): finish this
	return fmt.Sprintf("multi polygon %d polygons", len(p))
	// w := strings.Builder{}
	// w.WriteRune('[')
	// for _, p := range p {
	// 	w.WriteString(p.String() + " ")
	// }
	// w.WriteRune(']')
	// return w.String()
}

// Type implements the starlark.Value interface
func (p MultiPolygon) Type() string { return "MultiPolygon" }

// Freeze implements the starlark.Value interface, polygon is immutable
func (p MultiPolygon) Freeze() {}

// Truth implements the starlark.Value interface
func (p MultiPolygon) Truth() starlark.Bool {
	return len(p) > 0
}

// Hash implements the starlark.Value interface
func (p MultiPolygon) Hash() (h uint32, err error) {
	var hash uint32
	for _, poly := range p {
		if hash, err = poly.Hash(); err != nil {
			return
		}
		// TODO (b5): this is bad, unique values will report as non-unique. fix
		h += hash
	}
	return
}

// OrbMultiPolygon formats polygon as an orb.MultiPolygon
func (p MultiPolygon) OrbMultiPolygon() orb.MultiPolygon {
	polys := make([]orb.Polygon, len(p))
	for i, l := range p {
		polys[i] = l.OrbPolygon()
	}
	return orb.MultiPolygon(polys)
}

type multiPolygonIterator struct {
	idx     int
	polygon MultiPolygon
}

// If the iterator is exhausted, Next returns false.
// Otherwise it sets *p to the current element of the sequence,
// advances the iterator, and returns true.
func (li *multiPolygonIterator) Next(p *starlark.Value) bool {
	if li.idx == len(li.polygon) {
		return false
	}

	*p = li.polygon[li.idx]
	li.idx++
	return true
}

// Done is a no-op
// TODO (b5): this is supposed to be a no-op. right?
func (li *multiPolygonIterator) Done() {}

// Iterate implements the starlark iterator interface
func (p MultiPolygon) Iterate() starlark.Iterator {
	return &multiPolygonIterator{polygon: p}
}

// Index implements the starlark indexable interface
func (p MultiPolygon) Index(i int) starlark.Value {
	return p[i]
}

// Len implements the starlark indexable interface
func (p MultiPolygon) Len() int {
	return len(p)
}

// func Slice(start, end, step int) starlark.Value {
// }

// Attr gets an attribute of a polygon
func (p MultiPolygon) Attr(name string) (starlark.Value, error) {
	switch name {
	default:
		// attr does not exist
		return nil, nil
	}
}

// AttrNames returns all possible attribute names
func (p MultiPolygon) AttrNames() []string {
	names := []string{
		// methods
	}
	sort.Strings(names)
	return names
}
