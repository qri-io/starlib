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

	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('math.sky', 'math')
const ModuleName = "math.sky"

var (
	once       sync.Once
	mathModule skylark.StringDict
)

const tau = math.Pi * 2
const oneRad = tau / 360

// LoadModule loads the math module.
// It is concurrency-safe and idempotent.
func LoadModule() (skylark.StringDict, error) {
	once.Do(func() {
		inf := math.Inf(1)
		nan := math.NaN()
		mathModule = skylark.StringDict{
			"math": skylarkstruct.FromStringDict(skylarkstruct.Default, skylark.StringDict{
				"ceil":  skylark.NewBuiltin("ceil", ceil),
				"fabs":  skylark.NewBuiltin("fabs", fabs),
				"floor": skylark.NewBuiltin("floor", floor),

				"exp":  skylark.NewBuiltin("exp", exp),
				"sqrt": skylark.NewBuiltin("sqrt", sqrt),

				"acos":  skylark.NewBuiltin("acos", acos),
				"asin":  skylark.NewBuiltin("asin", asin),
				"atan":  skylark.NewBuiltin("atan", atan),
				"atan2": skylark.NewBuiltin("atan2", atan2),
				"cos":   skylark.NewBuiltin("cos", cos),
				"hypot": skylark.NewBuiltin("hypot", hypot),
				"sin":   skylark.NewBuiltin("sin", sin),
				"tan":   skylark.NewBuiltin("tan", tan),

				"degrees": skylark.NewBuiltin("degrees", degrees),
				"radians": skylark.NewBuiltin("radians", radians),

				"acosh": skylark.NewBuiltin("acosh", acosh),
				"asinh": skylark.NewBuiltin("asinh", asinh),
				"atanh": skylark.NewBuiltin("atanh", atanh),
				"cosh":  skylark.NewBuiltin("cosh", cosh),
				"sinh":  skylark.NewBuiltin("sinh", sinh),
				"tanh":  skylark.NewBuiltin("tanh", tanh),

				"e":   skylark.Float(math.E),
				"pi":  skylark.Float(math.Pi),
				"tau": skylark.Float(tau),
				"phi": skylark.Float(math.Phi),
				"inf": skylark.Float(inf),
				"nan": skylark.Float(nan),
			}),
		}
	})
	return mathModule, nil
}

// Return the floor of x, the largest integer less than or equal to x.
func floor(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("floor", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Floor(float64(x))), nil
}

// Return the ceiling of x, the smallest integer greater than or equal to x
func ceil(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("ceil", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Ceil(float64(x))), nil
}

// Return the absolute value of x
func fabs(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("fabs", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Abs(float64(x))), nil
}

// Return e raised to the power x, where e = 2.718281… is the base of natural logarithms.
func exp(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("exp", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Exp(float64(x))), nil
}

// Return the square root of x
func sqrt(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("sqrt", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Sqrt(float64(x))), nil
}
