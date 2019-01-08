package geo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"go.starlark.net/starlark"
)

func newLine(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var points *starlark.List

	if err = starlark.UnpackArgs("Line", args, kwargs, "points", &points); err != nil {
		return
	}

	return lineFromList(points)
}

func lineFromList(points *starlark.List) (line Line, err error) {
	var (
		v        starlark.Value
		lat, lng float64
		i        int
	)

	line = make(Line, points.Len())
	iter := points.Iterate()
	defer iter.Done()

	for iter.Next(&v) {
		switch pt := v.(type) {
		case Point:
			line[i] = pt
		case *starlark.List:
			if pt.Len() != 2 {
				err = fmt.Errorf("wrong number elements for point. expected 2, got: %d", pt.Len())
			}
			if lat, err = float64FromNumericValue(pt.Index(0)); err != nil {
				return
			}
			if lng, err = float64FromNumericValue(pt.Index(1)); err != nil {
				return
			}
			line[i] = Point{lat, lng}
		default:
			err = fmt.Errorf("invalid type '%s' to create line", pt.Type())
			return
		}
		i++
	}

	return
}

// Line is the starlark geographic line type
type Line []Point

// assert line is a starlark value
var _ starlark.Value = (*Line)(nil)

// String implements the starlark.Value interface
func (l Line) String() string {
	w := strings.Builder{}
	w.WriteRune('[')
	for _, p := range l {
		w.WriteString(p.String() + " ")
	}
	w.WriteRune(']')
	return w.String()
}

// Type implements the starlark.Value interface
func (l Line) Type() string { return "Line" }

// Freeze implements the starlark.Value interface, line is immutable
func (l Line) Freeze() {}

// Truth implements the starlark.Value interface
func (l Line) Truth() starlark.Bool {
	return len(l) > 0
}

// Hash implements the starlark.Value interface
func (l Line) Hash() (h uint32, err error) {
	var hash uint32
	for _, p := range l {
		if hash, err = p.Hash(); err != nil {
			return
		}

		// TODO (b5): this is bad, unique values will report as non-unique. fix
		h += hash
	}
	return
}

// IsClosed checks that first & last points on the line are equal
func (l Line) IsClosed() bool {
	if len(l) == 0 {
		return false
	}
	a := l[0]
	o := l[len(l)-1]
	return a[0] == o[0] && a[1] == o[1]
}

type lineIterator struct {
	idx  int
	line Line
}

// If the iterator is exhausted, Next returns false.
// Otherwise it sets *p to the current element of the sequence,
// advances the iterator, and returns true.
func (li *lineIterator) Next(p *starlark.Value) bool {
	if li.idx == len(li.line) {
		return false
	}

	*p = li.line[li.idx]
	li.idx++
	return true
}

// Done is a no-op
// TODO (b5): this is supposed to be a no-op. right?
func (li *lineIterator) Done() {}

// Iterate implements the starlark iterator interface
func (l Line) Iterate() starlark.Iterator {
	return &lineIterator{line: l}
}

// Index implements the starlark indexable interface
func (l Line) Index(i int) starlark.Value {
	return l[i]
}

// Len implements the starlark indexable interface
func (l Line) Len() int {
	return len(l)
}

// func Slice(start, end, step int) starlark.Value {
// }

// Attr gets an attribute of a line
func (l Line) Attr(name string) (starlark.Value, error) {
	switch name {
	case "buffer":
		return addClosure("buffer", l, l.buffer)
	case "length":
		return addClosure("length", l, l.length)
	case "lengthGeodesic":
		return addClosure("lengthGeodesic", l, l.lengthGeodesic)
	default:
		// attr does not exist
		return nil, nil
	}
}

// AttrNames returns all possible attribute names
func (l Line) AttrNames() []string {
	names := []string{
		// methods
		"buffer",
		"length",
		"lengthGeodesic",
	}
	sort.Strings(names)
	return names
}

func (l Line) buffer(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, fmt.Errorf("not yet implemented: buffer")
}

func (l Line) length(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	length := planar.Length(l.OrbLineString())
	return starlark.Float(length), nil
}

func (l Line) lengthGeodesic(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	length := geo.Length(l.OrbLineString())
	return starlark.Float(length), nil
}

// OrbLineString converts line to a orb.LineString
func (l Line) OrbLineString() orb.LineString {
	pts := make([]orb.Point, len(l))
	for i, pt := range l {
		pts[i] = orb.Point(pt)
	}
	return orb.LineString(pts)
}

// OrbRing converts line to a orb.Ring
func (l Line) OrbRing() orb.Ring {
	pts := make([]orb.Point, len(l))
	for i, pt := range l {
		pts[i] = orb.Point(pt)
	}
	return orb.Ring(pts)
}
