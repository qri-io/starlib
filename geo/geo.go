package geo

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('geo.star', 'geo')
const ModuleName = "geo.star"

var (
	once      sync.Once
	geoModule starlark.StringDict
)

// LoadModule loads the geo module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		geoModule = starlark.StringDict{
			"geo": starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
				"point": starlark.NewBuiltin("point", newPoint),
			}),
		}
	})
	return geoModule, nil
}

func newPoint(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		x, y     starlark.Value
		lat, lng float64
	)
	v = starlark.None

	if err = starlark.UnpackArgs("encode", args, kwargs, "x", &x, "y", &y); err != nil {
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
func (p Point) Type() string { return "point" }

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

// Line is the starlark geographic line type
type Line struct {
}

// Polygon is the starlark geographic polygon type
type Polygon struct {
}

type builtinMethod func(fnname string, recv starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func addClosure(name string, recv starlark.Value, method builtinMethod) (*starlark.Builtin, error) {
	// Allocate a closure over 'method'.
	impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return method(b.Name(), b.Receiver(), args, kwargs)
	}
	return starlark.NewBuiltin(name, impl).BindReceiver(recv), nil
}

func floatHash(f float64) (uint32, error) {
	// Equal float and int values must yield the same hash.
	if isFinite(f) {
		return finiteFloatToInt(f).Hash()
	}
	return 1618033, nil // NaN, +/-Inf
}

// isFinite reports whether f represents a finite rational value.
// It is equivalent to !math.IsNan(f) && !math.IsInf(f, 0).
func isFinite(f float64) bool {
	return math.Abs(f) <= math.MaxFloat64
}

// finiteFloatToInt converts f to an Int, truncating towards zero.
// f must be finite.
func finiteFloatToInt(f float64) starlark.Int {
	var i big.Int
	if math.MinInt64 <= f && f <= math.MaxInt64 {
		// small values
		i.SetInt64(int64(f))
	} else {
		rat := new(big.Rat).SetFloat64(f)
		if rat == nil {
			panic(f) // non-finite
		}
		i.Div(rat.Num(), rat.Denom())
	}
	return starlark.MakeInt(int(i.Int64()))
}

func float64FromNumericValue(n starlark.Value) (float64, error) {
	switch n.Type() {
	case "int":
		i, ok := n.(starlark.Int).Int64()
		if !ok {
			return 0, fmt.Errorf("invalid int")
		}
		return float64(i), nil
	case "float":
		return float64(n.(starlark.Float)), nil
	default:
		return 0, fmt.Errorf("invalid type '%s' expected int or float", n.Type())
	}
}
