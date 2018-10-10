/*Package time defines time primitives for starlark, based heavily on the time
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
- unix timet conversions
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

	starlark "github.com/google/starlark"
	starlarkstruct "github.com/google/starlark/starlarkstruct"
	syntax "github.com/google/starlark/syntax"
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
			"duration": starlark.NewBuiltin("duration", duration),
			"location": starlark.NewBuiltin("location", location),
			"now":      starlark.NewBuiltin("now", now),
			"struct":   starlark.NewBuiltin("struct", starlarkstruct.Make),
			"time":     starlark.NewBuiltin("time", time),
			"zero":     Time{},
		}

		// embed file into binary to remove any file dependencies
		file := strings.NewReader(`
time = struct(
	time = time,
	duration = duration,
	location = location,
	now = now,
	zero = zero,
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

func duration(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func location(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func time(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func now(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return Time(NowFunc()), nil
}

// Duration is a starlark representation of a duration
type Duration time.Duration

// String implements the Stringer interface
func (d Duration) String() string { return time.Duration(d).String() }

// Type returns a short string describing the value's type.
func (d Duration) Type() string { return "duration" }

// Freeze renders Duration immutable. required by starlark.Value interface
// because duration is already immutable this is a no-op
func (d Duration) Freeze() {}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface
func (d Duration) Hash() (uint32, error) { return hashString(d.String()), nil }

// Truth returns the truth value of an object required by starlark.Value interface
func (d Duration) Truth() starlark.Bool { return d > 0 }

// Attr gets a value for a string attribute, implementing dot expression support in starklark. required by starlark.HasAttrs interface
func (d Duration) Attr(name string) (starlark.Value, error) {
	return builtinAttr(d, name, durationMethods)
}

// AttrNames lists available dot expression strings. required by starlark.HasAttrs interface
func (d Duration) AttrNames() []string { return builtinAttrNames(durationMethods) }

// Binary implements binary operators, which satisfies the starlark.HasBinary interface
func (d Duration) Binary(op syntax.Token, yV starlark.Value, side starlark.Side) (starlark.Value, error) {
	x := time.Duration(d)
	var y time.Duration
	switch yV.(type) {
	case starlark.Int:
		i, ok := yV.(starlark.Int).Int64()
		if !ok {
			return nil, fmt.Errorf("duration binary operation: couldn't parse int")
		}
		y = time.Duration(i)
	case Duration:
		y = time.Duration(yV.(Duration))
	case Time:
		y := time.Time(yV.(Time))
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

// Time is a starlark representation of a point in time
type Time time.Time

// String implements the Stringer interface
func (t Time) String() string { return time.Time(t).String() }

// Type returns a short string describing the value's type.
func (t Time) Type() string { return "time" }

// Freeze renders time immutable. required by starlark.Value interface
// because Time is already immutable this is a no-op
func (t Time) Freeze() {}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface
func (t Time) Hash() (uint32, error) { return hashString(t.String()), nil }

// Truth returns the truth value of an object required by starlark.Value interface
func (t Time) Truth() starlark.Bool { return starlark.Bool(time.Time(t).IsZero()) }

// Attr gets a value for a string attribute, implementing dot expression support in starklark. required by starlark.HasAttrs interface
func (t Time) Attr(name string) (starlark.Value, error) { return builtinAttr(t, name, timeMethods) }

// AttrNames lists available dot expression strings for time. required by starlark.HasAttrs interface
func (t Time) AttrNames() []string { return builtinAttrNames(timeMethods) }

// CompareSameType implements comparison of two Time values. required by starlark.Comparable interface
func (t Time) CompareSameType(op syntax.Token, yV starlark.Value, depth int) (bool, error) {
	x := time.Time(t)
	y := time.Time(yV.(Time))
	cp := 0
	if x.Before(y) {
		cp = -1
	} else if x.After(y) {
		cp = 1
	}
	return threeway(op, cp), nil
}

// Binary implements binary operators, which satisfies the starlark.HasBinary interface
func (t Time) Binary(op syntax.Token, yV starlark.Value, side starlark.Side) (starlark.Value, error) {
	x := time.Time(t)

	switch yV.(type) {
	case Duration:
		y := time.Duration(yV.(Duration))
		switch op {
		// time + duration = time
		case syntax.PLUS:
			return Time(x.Add(y)), nil
		// time - duration = time
		case syntax.MINUS:
			return Time(x.Add(-y)), nil
		}
	case Time:
		y := time.Time(yV.(Time))
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
	"year":       timeyear,
	"month":      timemonth,
	"day":        timeday,
	"hour":       timehour,
	"minute":     timeminute,
	"second":     timesecond,
	"nanosecond": timenanosecond,
}

// TODO - consider using a higher order function to generate these
func timeyear(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(recv.Year()), nil
}

func timemonth(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(int(recv.Month())), nil
}

func timeday(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(recv.Day()), nil
}

func timehour(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(recv.Hour()), nil
}

func timeminute(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(recv.Minute()), nil
}

func timesecond(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
	return starlark.MakeInt(recv.Second()), nil
}

func timenanosecond(fnname string, recV starlark.Value, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	recv := time.Time(recV.(Time))
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
