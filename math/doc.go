/*Package math is a Starlark module of math-related functions and constants.
math was upstreamed from starlib into go-Starlark. This package exists to add
documentation only. The API is locked to strictly match the Starlark module.
Users are encouraged to import the math package directly via:
go.starlark.net/lib/math

For source code see
https://github.com/google/starlark-go/tree/master/lib/math

outline: math
  math defines a Starlark module of mathematical functions. All functions accept
  both int and float values as arguments.
  path: math
  constants:
    e: The base of natural logarithms, approximately 2.71828.
    pi: The ratio of a circle's circumference to its diameter, approximately 3.14159.
  functions:
    acos(x)
      Return the arc cosine of x, in radians.
    acosh(x)
      Return the inverse hyperbolic cosine of x.
    asin(x)
      Return the arc sine of x, in radians.
    asinh(x)
      Return the inverse hyperbolic sine of x.
    atan(x)
      Return the arc tangent of x, in radians.
    atan2(y, x)
      Return atan(y / x), in radians. The result is between -pi and pi.
      The vector in the plane from the origin to point (x, y) makes this angle
      with the positive X axis. The point of atan2() is that the signs of both
      inputs are known to it, so it can compute the correct quadrant for the
      angle. For example, atan(1) and atan2(1, 1) are both pi/4, but
      atan2(-1, -1) is -3*pi/4.
    atanh(x)
      Return the inverse hyperbolic tangent of x.
    ceil(x)
      Return the ceiling of x, the smallest integer greater than or equal to x.
    copysign(x,y)
      Returns a value with the magnitude of x and the sign of y.
    cos(x)
      Return the cosine of x radians.
    cosh(x)
      Return the hyperbolic cosine of x.
    degrees(x)
      Convert angle x from radians to degrees.
    exp(x)
      Returns e raised to the power x, where e = 2.718281… is the base of
      natural logarithms.
    fabs(x)
      Return the absolute value of x.
    floor(x)
      Return the floor of x, the largest integer less than or equal to x.
    gamma(x)
      Returns the Gamma function of x.
    hypot(x, y)
      Return the Euclidean norm, sqrt(x*x + y*y). This is the length of the
      vector from the origin to point (x, y).
    log(x, base)
      Returns the logarithm of x in the given base, or natural logarithm by
      default.
    mod(x, y)
      Returns the floating-point remainder of x/y. The magnitude of the result
      is less than y and its sign agrees with that of x.
    pow(x, y)
      Returns x**y, the base-x exponential of y.
    radians(x)
      Convert angle x from degrees to radians.
    remainder(x, y)
      Returns the IEEE 754 floating-point remainder of x/y.
    round(x)
      Returns the nearest integer, rounding half away from zero.
    sqrt(x)
      Return the square root of x.
    sin(x)
      Return the sine of x radians.
    sinh(x)
      Return the hyperbolic sine of x.
    tan(x)
      Return the tangent of x radians.
    tanh(x)
      Return the hyperbolic tangent of x.
*/
package math

import "go.starlark.net/lib/math"

// ModuleName declares the intended load import string
// eg: load("math.star", "math")
const ModuleName = "math.star"

// Module exposes the time module. Implementation located at
// https://github.com/google/starlark-go/tree/master/lib/math
var Module = math.Module
