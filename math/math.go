/*
Package math defines mathimatical functions, it's intented to be a drop-in
subset of python's math module for starlark:
https://docs.python.org/3/library/math.html

Currently defined functions are as follows:

Number-theoretic and representation functions:

	ceil(x) - Return the ceiling of x, the smallest integer greater than or equal to x.
	fabs(x) - Return the absolute value of x.
	floor(x) - Return the floor of x, the largest integer less than or equal to x.

Power and logarithmic functions

	exp(x) - Return e raised to the power x, where e = 2.718281… is the base of natural logarithms
	sqrt(x) - Return the square root of x.

Trigonometric functions

	acos(x) - Return the arc cosine of x, in radians.
	asin(x) - Return the arc sine of x, in radians.
	atan(x) - Return the arc tangent of x, in radians.
	atan2(y, x) - Return atan(y / x), in radians. The result is between -pi and pi. The vector in the plane from the origin to point (x, y) makes this angle with the positive X axis. The point of atan2() is that the signs of both inputs are known to it, so it can compute the correct quadrant for the angle. For example, atan(1) and atan2(1, 1) are both pi/4, but atan2(-1, -1) is -3*pi/4.
	cos(x) - Return the cosine of x radians.
	hypot(x, y) - Return the Euclidean norm, sqrt(x*x + y*y). This is the length of the vector from the origin to point (x, y).
	sin(x) - Return the sine of x radians.
	tan(x) - Return the tangent of x radians.

Angular conversion

	degrees(x) - Convert angle x from radians to degrees.
	radians(x) - Convert angle x from degrees to radians.

Hyperbolic functions - Hyperbolic functions are analogs of trigonometric functions that are based on hyperbolas instead of circles.

	acosh(x) - Return the inverse hyperbolic cosine of x.
	asinh(x) - Return the inverse hyperbolic sine of x.
	atanh(x) - Return the inverse hyperbolic tangent of x.
	cosh(x) - Return the hyperbolic cosine of x.
	sinh(x) - Return the hyperbolic sine of x.
	tanh(x) - Return the hyperbolic tangent of x.
*/
package math

import (
	"math"
	"sync"

	starlark "github.com/google/skylark"
	starlarkstruct "github.com/google/skylark/skylarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('math.star', 'math')
const ModuleName = "math.star"

var (
	once       sync.Once
	mathModule starlark.StringDict
)

const tau = math.Pi * 2
const oneRad = tau / 360

// LoadModule loads the math module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		inf := math.Inf(1)
		nan := math.NaN()
		mathModule = starlark.StringDict{
			"math": starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
				"ceil":  starlark.NewBuiltin("ceil", ceil),
				"fabs":  starlark.NewBuiltin("fabs", fabs),
				"floor": starlark.NewBuiltin("floor", floor),

				"exp":  starlark.NewBuiltin("exp", exp),
				"sqrt": starlark.NewBuiltin("sqrt", sqrt),

				"acos":  starlark.NewBuiltin("acos", acos),
				"asin":  starlark.NewBuiltin("asin", asin),
				"atan":  starlark.NewBuiltin("atan", atan),
				"atan2": starlark.NewBuiltin("atan2", atan2),
				"cos":   starlark.NewBuiltin("cos", cos),
				"hypot": starlark.NewBuiltin("hypot", hypot),
				"sin":   starlark.NewBuiltin("sin", sin),
				"tan":   starlark.NewBuiltin("tan", tan),

				"degrees": starlark.NewBuiltin("degrees", degrees),
				"radians": starlark.NewBuiltin("radians", radians),

				"acosh": starlark.NewBuiltin("acosh", acosh),
				"asinh": starlark.NewBuiltin("asinh", asinh),
				"atanh": starlark.NewBuiltin("atanh", atanh),
				"cosh":  starlark.NewBuiltin("cosh", cosh),
				"sinh":  starlark.NewBuiltin("sinh", sinh),
				"tanh":  starlark.NewBuiltin("tanh", tanh),

				"e":   starlark.Float(math.E),
				"pi":  starlark.Float(math.Pi),
				"tau": starlark.Float(tau),
				"phi": starlark.Float(math.Phi),
				"inf": starlark.Float(inf),
				"nan": starlark.Float(nan),
			}),
		}
	})
	return mathModule, nil
}

// floatFunc unpacks a starlark function call, calls a passed in float64 function
// and returns the result as a starlark value
func floatFunc(name string, args starlark.Tuple, kwargs []starlark.Tuple, fn func(float64) float64) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs(name, args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(fn(float64(x))), nil
}

// floatFunc2 is a 2-argument float func
func floatFunc2(name string, args starlark.Tuple, kwargs []starlark.Tuple, fn func(float64, float64) float64) (starlark.Value, error) {
	var x, y starlark.Float
	if err := starlark.UnpackArgs(name, args, kwargs, "x", &x, "y", &y); err != nil {
		return nil, err
	}
	return starlark.Float(fn(float64(x), float64(y))), nil
}

// Return the floor of x, the largest integer less than or equal to x.
func floor(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return floatFunc("floor", args, kwargs, math.Floor)
}

// Return the ceiling of x, the smallest integer greater than or equal to x
func ceil(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return floatFunc("ceil", args, kwargs, math.Ceil)
}

// Return the absolute value of x
func fabs(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return floatFunc("fabs", args, kwargs, math.Abs)
}

// Return e raised to the power x, where e = 2.718281… is the base of natural logarithms.
func exp(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return floatFunc("exp", args, kwargs, math.Exp)
}

// Return the square root of x
func sqrt(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return floatFunc("sqrt", args, kwargs, math.Sqrt)
}
