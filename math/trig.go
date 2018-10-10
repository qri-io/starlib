package math

import (
	"math"

	starlark "github.com/google/skylark"
)

// Return the arc cosine of x, in radians.
func acos(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("acos", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Acos(float64(x))), nil
}

// asin(x) - Return the arc sine of x, in radians.
func asin(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("asin", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Asin(float64(x))), nil
}

// atan(x) - Return the arc tangent of x, in radians.
func atan(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("atan", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Atan(float64(x))), nil
}

// atan2(y, x) - Return atan(y / x), in radians. The result is between -pi and pi. The vector in the plane from the origin to point (x, y) makes this angle with the positive X axis. The point of atan2() is that the signs of both inputs are known to it, so it can compute the correct quadrant for the angle. For example, atan(1) and atan2(1, 1) are both pi/4, but atan2(-1, -1) is -3*pi/4.
func atan2(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x, y starlark.Float
	if err := starlark.UnpackArgs("atan2", args, kwargs, "y", &y, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Atan2(float64(y), float64(x))), nil
}

// cos(x) - Return the cosine of x radians.
func cos(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("cos", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Cos(float64(x))), nil
}

// hypot(x, y) - Return the Euclidean norm, sqrt(x*x + y*y). This is the length of the vector from the origin to point (x, y).
func hypot(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x, y starlark.Float
	if err := starlark.UnpackArgs("hypot", args, kwargs, "x", &x, "y", &y); err != nil {
		return nil, err
	}
	return starlark.Float(math.Hypot(float64(x), float64(y))), nil
}

// sin(x) - Return the sine of x radians.
func sin(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("sin", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Sin(float64(x))), nil
}

// tan(x) - Return the tangent of x radians.
func tan(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("tan", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Tan(float64(x))), nil
}

// degrees(x) - Convert angle x from radians to degrees.
func degrees(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("degrees", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(float64(x) / oneRad), nil
}

// radians(x) - Convert angle x from degrees to radians.
func radians(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("radians", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(float64(x) * oneRad), nil
}

// acosh(x) - Return the inverse hyperbolic cosine of x.
func acosh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("acosh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Acosh(float64(x))), nil
}

// asinh(x) - Return the inverse hyperbolic sine of x.
func asinh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("asinh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Asinh(float64(x))), nil
}

// atanh(x) - Return the inverse hyperbolic tangent of x.
func atanh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("atanh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Atanh(float64(x))), nil
}

// cosh(x) - Return the hyperbolic cosine of x.
func cosh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("cosh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Cosh(float64(x))), nil
}

// sinh(x) - Return the hyperbolic sine of x.
func sinh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("sinh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Sinh(float64(x))), nil
}

// tanh(x) - Return the hyperbolic tangent of x.
func tanh(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("tanh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return starlark.Float(math.Tanh(float64(x))), nil
}
