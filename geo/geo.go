package geo

import (
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/paulmach/orb"
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
				// constructors
				"Point":   starlark.NewBuiltin("Point", newPoint),
				"Line":    starlark.NewBuiltin("Line", newLine),
				"Polygon": starlark.NewBuiltin("Polygon", newPolygon),

				// geographic joins
				"within":       starlark.NewBuiltin("within", within),
				"intersects":   starlark.NewBuiltin("intersects", intersects),
				"parseGeoJSON": starlark.NewBuiltin("parseGeoJSON", parseGeoJSON),
			}),
		}
	})
	return geoModule, nil
}

func within(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		a starlark.Value
		b Polygon
	)
	v = starlark.None

	if err = starlark.UnpackArgs("within", args, kwargs, "a", &a, "b", &b); err != nil {
		return
	}

	poly := b.OrbPolygon()
	switch geom := a.(type) {
	case Line:
		for _, pt := range geom.OrbLineString() {
			if !planar.PolygonContains(poly, pt) {
				return starlark.Bool(false), nil
			}
			return starlark.Bool(true), nil
		}
	case Point:
		within := planar.PolygonContains(poly, orb.Point{geom[0], geom[1]})
		return starlark.Bool(within), nil
	case Polygon:
		err = fmt.Errorf("checking polygons-within-polygons is not yet supported")
		return
	default:
		err = fmt.Errorf("unrecognized type: %s", a.Type())
	}

	return starlark.None, fmt.Errorf("not finished: within")
}

func intersects(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	return starlark.None, fmt.Errorf("not finished: intersects")
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

func float64FromNumericValue(v starlark.Value) (float64, error) {
	switch n := v.(type) {
	case starlark.Int:
		i, ok := n.Int64()
		if !ok {
			return 0, fmt.Errorf("invalid int")
		}
		return float64(i), nil
	case starlark.Float:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("invalid type '%s' expected int or float", n.Type())
	}
}
