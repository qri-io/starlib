/*Package time defines time primitives for starlark, based heavily on the time
package from the go standard library.

  outline: time
    time defines time primitives for starlark
    functions:
      duration(string) duration
        parse a duration
      location(string) location
        parse a location
      time(string, format=..., location=...) time
        parse a time
      now() time
        implementations would be able to make this a constant
      zero() time
        a constant

    types:
      duration
        fields:
          hours float
          minutes float
          nanoseconds int
          seconds float
        operators:
          duration - time = duration
          duration + time = time
          duration == duration = boolean
          duration < duration = booleans
      time
        fields:
          year int
          month int
          day int
          hour int
          minute int
          second int
          nanosecond int
        operators:
          time == time = boolean
          time < time = boolean
          time + duration = time
          time - duration = time
          time - time = duration
*/
package time
