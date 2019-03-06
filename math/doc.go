/*Package math defines mathematical functions, it's intended to be a drop-in
subset of python's math module for starlark: https://docs.python.org/3/library/math.html

Currently defined functions are as follows:

  outline: math
    math defines mathematical functions, it's intended to be a drop-in
    subset of python's math module for starlark: https://docs.python.org/3/library/math.html

    functions:
      ceil(x)
        Return the ceiling of x, the smallest integer greater than or equal to x.
      fabs(x)
        Return the absolute value of x.
      floor(x)
        Return the floor of x, the largest integer less than or equal to x.
      round(x)
        Returns the nearest integer, rounding half away from zero.
      exp(x)
        Return e raised to the power x, where e = 2.718281… is the base of natural logarithms
      sqrt(x)
        Return the square root of x.
      asin(x)
        Return the arc sine of x, in radians.
      acos(x)
        Return the arc cosine of x, in radians.
      atan(x)
        Return the arc tangent of x, in radians.
      atan2(y, x)
        Return atan(y / x), in radians. The result is between -pi and pi. The vector in the plane from the origin to point (x, y) makes this angle with the positive X axis. The point of atan2() is that the signs of both inputs are known to it, so it can compute the correct quadrant for the angle. For example, atan(1) and atan2(1, 1) are both pi/4, but atan2(-1, -1) is -3*pi/4.
      cos(x)
        Return the cosine of x radians.
      hypot(x, y)
        Return the Euclidean norm, sqrt(x*x + y*y). This is the length of the vector from the origin to point (x, y).
      sin(x)
        Return the sine of x radians.
      tan(x)
        Return the tangent of x radians.
      degrees(x)
        Convert angle x from radians to degrees.
      radians(x)
        Convert angle x from degrees to radians.
      acosh(x)
        Return the inverse hyperbolic cosine of x.
      asinh(x)
        Return the inverse hyperbolic sine of x.
      atanh(x)
        Return the inverse hyperbolic tangent of x.
      cosh(x)
        Return the hyperbolic cosine of x.
      sinh(x)
        Return the hyperbolic sine of x.
      tanh(x)
        Return the hyperbolic tangent of x.
*/
package math
