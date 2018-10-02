package math

import (
	"math"

	"github.com/google/skylark"
)

// Return the arc cosine of x, in radians.
func acos(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("acos", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Acos(float64(x))), nil
}

// asin(x) - Return the arc sine of x, in radians.
func asin(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("asin", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Asin(float64(x))), nil
}

// atan(x) - Return the arc tangent of x, in radians.
func atan(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("atan", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Atan(float64(x))), nil
}

// atan2(y, x) - Return atan(y / x), in radians. The result is between -pi and pi. The vector in the plane from the origin to point (x, y) makes this angle with the positive X axis. The point of atan2() is that the signs of both inputs are known to it, so it can compute the correct quadrant for the angle. For example, atan(1) and atan2(1, 1) are both pi/4, but atan2(-1, -1) is -3*pi/4.
func atan2(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x, y skylark.Float
	if err := skylark.UnpackArgs("atan2", args, kwargs, "y", &y, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Atan2(float64(y), float64(x))), nil
}

// cos(x) - Return the cosine of x radians.
func cos(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("cos", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Cos(float64(x))), nil
}

// hypot(x, y) - Return the Euclidean norm, sqrt(x*x + y*y). This is the length of the vector from the origin to point (x, y).
func hypot(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x, y skylark.Float
	if err := skylark.UnpackArgs("hypot", args, kwargs, "x", &x, "y", &y); err != nil {
		return nil, err
	}
	return skylark.Float(math.Hypot(float64(x), float64(y))), nil
}

// sin(x) - Return the sine of x radians.
func sin(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("sin", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Sin(float64(x))), nil
}

// tan(x) - Return the tangent of x radians.
func tan(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("tan", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Tan(float64(x))), nil
}

// degrees(x) - Convert angle x from radians to degrees.
func degrees(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("degrees", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(float64(x) / oneRad), nil
}

// radians(x) - Convert angle x from degrees to radians.
func radians(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("radians", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(float64(x) * oneRad), nil
}

// acosh(x) - Return the inverse hyperbolic cosine of x.
func acosh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("acosh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Acosh(float64(x))), nil
}

// asinh(x) - Return the inverse hyperbolic sine of x.
func asinh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("asinh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Asinh(float64(x))), nil
}

// atanh(x) - Return the inverse hyperbolic tangent of x.
func atanh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("atanh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Atanh(float64(x))), nil
}

// cosh(x) - Return the hyperbolic cosine of x.
func cosh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("cosh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Cosh(float64(x))), nil
}

// sinh(x) - Return the hyperbolic sine of x.
func sinh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("sinh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Sinh(float64(x))), nil
}

// tanh(x) - Return the hyperbolic tangent of x.
func tanh(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("tanh", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	return skylark.Float(math.Tanh(float64(x))), nil
}
