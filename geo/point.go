package geo

import (
	"fmt"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"go.starlark.net/starlark"
)

func newPoint(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		x, y     starlark.Value
		lat, lng float64
	)

	if err = starlark.UnpackArgs("Point", args, kwargs, "x", &x, "y", &y); err != nil {
		return
	}

	if lat, err = float64FromNumericValue(x); err != nil {
		return
	}
	if lng, err = float64FromNumericValue(y); err != nil {
		return
	}

	return Point{lat, lng}, nil
}

// Point is the starlark geographic point type
type Point [2]float64

// assert point is a starlark value
var _ starlark.Value = (*Point)(nil)

// String implements the starlark.Value interface
func (p Point) String() string { return fmt.Sprintf("(%f,%f)", p[0], p[1]) }

// Type implements the starlark.Value interface
func (p Point) Type() string { return "Point" }

// Freeze implements the starlark.Value interface, point is immutable
func (p Point) Freeze() {}

// Truth implements the starlark.Value interface
func (p Point) Truth() starlark.Bool {
	return starlark.Bool(p[0] != 0 && p[1] != 0)
}

// Hash implements the starlark.Value interface
func (p Point) Hash() (uint32, error) {
	x, _ := floatHash(p[0])
	y, _ := floatHash(p[1])
	// TODO (b5): this is bad. fix
	return x - y, nil
}

// Attr gets an attribute of point
func (p Point) Attr(name string) (starlark.Value, error) {
	switch name {
	case "x", "lat":
		return starlark.Float(p[0]), nil
	case "y", "lng":
		return starlark.Float(p[1]), nil
	case "buffer":
		return addClosure("buffer", p, p.buffer)
	case "distance":
		return addClosure("distance", p, p.distance)
	case "distanceGeodesic":
		return addClosure("distanceGeodesic", p, p.distanceGeodesic)
	case "KNN":
		return addClosure("KNN", p, p.knn)
	default:
		// attr does not exist
		return nil, nil
	}
}

// AttrNames returns all possible attribute names
func (p Point) AttrNames() []string {
	names := []string{
		// attributes
		"x", "y", "lat", "lng",

		// methods
		"buffer",
		"distance",
		"distanceGeodesic",
		"KNN",
	}
	sort.Strings(names)
	return names
}

func (p Point) buffer(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, fmt.Errorf("not yet implemented: buffer")
}

func (p Point) distance(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var p2 Point
	if err := starlark.UnpackArgs("distance", args, kwargs, "p2", &p2); err != nil {
		return nil, err
	}

	d := planar.Distance(orb.Point(p), orb.Point(p2))
	return starlark.Float(d), nil
}

func (p Point) distanceGeodesic(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var p2 Point
	if err := starlark.UnpackArgs("distanceGeodesic", args, kwargs, "p2", &p2); err != nil {
		return nil, err
	}

	d := geo.Distance(orb.Point(p), orb.Point(p2))
	return starlark.Float(d), nil
}

func (p Point) knn(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, fmt.Errorf("not yet implemented: knn")
}
