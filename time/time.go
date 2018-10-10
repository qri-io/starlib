/*
Package sltime defines time primitives for starkark, based heavily on the time
package from the go standard library.

module time

functions
	duration(string) duration                               # parse a duration
	location(string) location                               # parse a location
	time(string, format=..., location=...) time             # parse a time
	now() time # implementations would be able to make this a constant
	zero time # a constant

type duration
operators
	duration - time = duration
	duration + time = time
	duration == duration
	duration < duration
fields
	hours float
	minutes float
	nanoseconds int
	seconds float

type time
operators
	time == time
	time < time
	time + duration = time
	time - duration = time
	time - time = duration
fields
	year int
	month int
	day int
	hour int
	minute int
	second int
	nanosecond int

TODO:
- unix time_t conversions
- timezone stuff
- strftime formatting
- constructor from 6 components + location
*/
package time

import (
	"fmt"
	"go/build"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	starlark "github.com/google/skylark"
	starlarkstruct "github.com/google/skylark/skylarkstruct"
	syntax "github.com/google/skylark/syntax"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('time.sky', 'time')
const ModuleName = "time.sky"

var (
	once       sync.Once
	timeModule starlark.StringDict
	timeError  error
)

// LoadModule loads the time module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		predeclared := starlark.StringDict{
			"duration_": starlark.NewBuiltin("duration", duration_),
			"location_": starlark.NewBuiltin("location", location_),
			"now_":      starlark.NewBuiltin("now", now_),
			"struct":    starlark.NewBuiltin("struct", starlarkstruct.Make),
			"time_":     starlark.NewBuiltin("time", time_),
			"zero_":     Time{},
		}

		// embed file into binary to remove any file dependencies
		file := strings.NewReader(`
time = struct(
	time = time_,
	duration = duration_,
	location = location_,
	now = now_,
	zero = zero_,
)
`)

		// filename := DataFile("time", "time.sky")
		thread := new(starlark.Thread)
		timeModule, timeError = starlark.ExecFile(thread, "time.sky", file, predeclared)
	})
	return timeModule, timeError
}

// NowFunc is a function that generates the current time. Intentionally exported
// so that it can be overridden
var NowFunc = func() time.Time { return time.Now() }

func duration_(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.String
	if err := starlark.UnpackArgs("duration", args, kwargs, "x", &x); err != nil {
		return nil, err
	}

	d, err := time.ParseDuration(string(x))
	if err != nil {
		return nil, err
	}

	return Duration(d), nil
}

func location_(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.String
	if err := starlark.UnpackArgs("location", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(string(x))
	if err != nil {
		return nil, err
	}

	return starlark.String(loc.String()), nil
}

func time_(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		x, location starlark.String
		format      = starlark.String(time.RFC3339)
	)
	if err := starlark.UnpackArgs("time", args, kwargs, "x", &x, "format?", &format, "location", &location); err != nil {
		return nil, err
	}

	if location == "" {
		t, err := time.Parse(format.String(), x.String())
		if err != nil {
			return nil, err
		}
		return Time(t), nil
	}

	loc, err := time.LoadLocation(location.String())
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(format.String(), x.String(), loc)
	if err != nil {
		return nil, err
	}
	return Time(t), nil
}

func now_(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return Time(NowFunc()), nil
}

type Duration time.Duration

func (d Duration) String() string        { return time.Duration(d).String() }
func (d Duration) Type() string          { return "duration" }
func (d Duration) Freeze()               {} // immutable
func (d Duration) Hash() (uint32, error) { return hashString(d.String()), nil }
func (d Duration) Truth() starlark.Bool  { return d > 0 }
func (d Duration) Attr(name string) (starlark.Value, error) {
	return builtinAttr(d, name, durationMethods)
}
func (d Duration) AttrNames() []string { return builtinAttrNames(durationMethods) }
func (d Duration) Binary(op syntax.Token, y_ starlark.Value, side starlark.Side) (starlark.Value, error) {
	x := time.Duration(d)
	var y time.Duration
	switch y_.(type) {
	case starlark.Int:
		i, ok := y_.(starlark.Int).Int64()
		if !ok {
			return nil, fmt.Errorf("duration binary operation: couldn't parse int")
		}
		y = time.Duration(i)
	case Duration:
		y = time.Duration(y_.(Duration))
	case Time:
		y := time.Time(y_.(Time))
		switch op {
		case syntax.PLUS:
			// duration + time = time
			return Time(y.Add(x)), nil
		case syntax.MINUS:
			// duration - time = duration
			return nil, nil
		}
	default:
		return nil, nil
	}

	switch op {
	case syntax.PLUS:
		return Duration(x + y), nil
	case syntax.MINUS:
		return Duration(x - y), nil
	case syntax.SLASH:
		if int64(y) == 0 {
			return nil, fmt.Errorf("cannot divide duration by zero")
		}
		return Duration(x / y), nil
	case syntax.STAR:
		return Duration(x * y), nil
	}

	return nil, nil
}

var durationMethods = map[string]builtinMethod{
	// "hours" :
	// "minutes" :
	// "nanoseconds" :
	// "seconds" :
}

type Time time.Time

func (t Time) String() string                           { return time.Time(t).String() }
func (t Time) Type() string                             { return "time" }
func (t Time) Freeze()                                  {} // immutable
func (t Time) Hash() (uint32, error)                    { return hashString(t.String()), nil }
func (t Time) Truth() starlark.Bool                     { return starlark.Bool(time.Time(t).IsZero()) }
func (t Time) Attr(name string) (starlark.Value, error) { return builtinAttr(t, name, timeMethods) }
func (t Time) AttrNames() []string                      { return builtinAttrNames(timeMethods) }
func (t Time) CompareSameType(op syntax.Token, y_ starlark.Value, depth int) (bool, error) {
	x := time.Time(t)
	y := time.Time(y_.(Time))
	cp := 0
	if x.Before(y) {
		cp = -1
	} else if x.After(y) {
		cp = 1
	}
	return threeway(op, cp), nil
}
func (t Time) Binary(op syntax.Token, y_ starlark.Value, side starlark.Side) (starlark.Value, error) {
	x := time.Time(t)

	switch y_.(type) {
	case Duration:
		y := time.Duration(y_.(Duration))
		switch op {
		// time + duration = time
		case syntax.PLUS:
			return Time(x.Add(y)), nil
		// time - duration = time
		case syntax.MINUS:
			return Time(x.Add(-y)), nil
		}
	case Time:
		y := time.Time(y_.(Time))
		switch op {
		// time - time = duration
		case syntax.MINUS:
			if side == starlark.Left {
				return Duration(x.Sub(y)), nil
			}
			return Duration(y.Sub(x)), nil
		}
	}

	// dunno, bail
	return nil, nil
}

var timeMethods = map[string]builtinMethod{
	"year":       time_year,
	"month":      time_month,
	"day":        time_day,
	"hour":       time_hour,
	"minute":     time_minute,
	"second":     time_second,
	"nanosecond": time_nanosecond,
}

// TODO - consider using a higher order function to generate these
func time_year(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Year()), nil
}

func time_month(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(int(recv.Month())), nil
}

func time_day(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Day()), nil
}

func time_hour(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Hour()), nil
}

func time_minute(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Minute()), nil
}

func time_second(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Second()), nil
}

func time_nanosecond(fnname string, recv_ starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recv_.(Time))
	return starlark.MakeInt(recv.Nanosecond()), nil
}

type builtinMethod func(fnname string, recv starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func builtinAttr(recv starlark.Value, name string, methods map[string]builtinMethod) (starlark.Value, error) {
	method := methods[name]
	if method == nil {
		return nil, nil // no such method
	}

	// Allocate a closure over 'method'.
	impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return method(b.Name(), b.Receiver(), args, kwargs)
	}
	return starlark.NewBuiltin(name, impl).BindReceiver(recv), nil
}

func builtinAttrNames(methods map[string]builtinMethod) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// hashString computes the FNV hash of s.
func hashString(s string) uint32 {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// threeway interprets a three-way comparison value cmp (-1, 0, +1)
// as a boolean comparison (e.g. x < y).
func threeway(op syntax.Token, cmp int) bool {
	switch op {
	case syntax.EQL:
		return cmp == 0
	case syntax.NEQ:
		return cmp != 0
	case syntax.LE:
		return cmp <= 0
	case syntax.LT:
		return cmp < 0
	case syntax.GE:
		return cmp >= 0
	case syntax.GT:
		return cmp > 0
	}
	panic(op)
}

// DataFile returns the effective filename of the specified
// test data resource.  The function abstracts differences between
// 'go build', under which a test runs in its package directory,
// and Blaze, under which a test runs in the root of the tree.
var DataFile = func(pkgdir, filename string) string {
	return filepath.Join(build.Default.GOPATH, "src/github.com/qri-io", pkgdir, filename)
}
