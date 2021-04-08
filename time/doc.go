/*Package time provides time-related constants and functions. The time module
was upstreamed from starlib into go-Starlark. This package exists to add
documentation. The API is locked to strictly match the Starlark module.
Users are encouraged to import the time package directly via:
go.starlark.net/lib/time

For source code see
https://github.com/google/starlark-go/tree/master/lib/time

outline: time
  time is a Starlark module of time-related functions and constants.
  path: time
  constants:
    nanosecond: A duration representing one nanosecond.
    microsecond: A duration representing one microsecond.
    millisecond: A duration representing one millisecond.
    second: A duration representing one second.
    minute: A duration representing one minute.
    hour: duration representing one hour.
  functions:
    from_timestamp(sec, nsec) Time
      Converts the given Unix time corresponding to the number of seconds
      and (optionally) nanoseconds since January 1, 1970 UTC into an object
      of type Time. For more details, refer to https://pkg.go.dev/time#Unix.
    is_valid_timezone(loc) boolean
      Reports whether loc is a valid time zone name.
    now() time
      Returns the current local time
    parse_duration(d) Duration
      Parses the given duration string. For more details, refer to
      https://pkg.go.dev/time#ParseDuration.
    parseTime(x, format, location) Time
      Parses the given time string using a specific time format and location.
      The expected arguments are a time string (mandatory), a time format
      (optional, set to RFC3339 by default, e.g. "2021-03-22T23:20:50.52Z")
      and a name of location (optional, set to UTC by default). For more
      details, refer to https://pkg.go.dev/time#Parse and
      https://pkg.go.dev/time#ParseInLocation.
    time(year?, month?, day?, hour?, minute?, second?, nanosecond?, location?) Time
      Returns the Time corresponding to yyyy-mm-dd hh:mm:ss + nsec nanoseconds
      in the appropriate zone for that time in the given location. All
      parameters are optional.
  types:
    Duration
      fields:
        hours float
        minutes float
        seconds float
        milliseconds int
        microseconds int
        nanoseconds int
      operators:
        duration + duration = duration
        duration + time = time
        duration - duration = duration
        duration / duration = float
        duration / int = duration
        duration / float = duration
        duration // duration = int
        duration * int = duration
    Time
      fields:
        year int
        month int
        day int
        hour int
        minute int
        second int
        nanosecond int
        unix int
        unix_nano int
      functions:
        in_location(locstr) Time
          get time representing the same instant but in a different location
        format() string
          textual representation of time formatted according to the provided
          layout string
      operators:
        time + duration = time
        time - duration = time
        time - time = duration
*/
package time

import "go.starlark.net/lib/time"

// ModuleName declares the intended load import string
// eg: load("time.star", "time")
const ModuleName = "time.star"

// Module exposes the time module. Implementation located at
// https://github.com/google/starlark-go/tree/master/lib/time
var Module = time.Module
