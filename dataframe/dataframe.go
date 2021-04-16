package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

const Name = "dataframe"

// Module exposes the dataframe module
var Module = &starlarkstruct.Module{
	Name: Name,
	Members: starlark.StringDict{
		"DataFrame": starlark.NewBuiltin("DataFrame", newDataFrame),
		"Index":     starlark.NewBuiltin("Index", newIndex),
	},
}

func unfinishedError(v starlark.Value, msg string) error {
	return fmt.Errorf("%s %s unfinished implementation: %s", Name, v.Type(), msg)
}

type DataFrame struct {
	frozen bool
}

// compile-time interface assertions
var (
	_ starlark.Value   = (*DataFrame)(nil)
	_ starlark.Mapping = (*DataFrame)(nil)
)

func newDataFrame(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		data, index, columns, dtype starlark.Value
		kopy                        starlark.Bool
	)
	if err := starlark.UnpackArgs("DataFrame", args, kwargs,
		"data?", &data,
		"index?", &index,
		"columns?", &columns,
		"dtype?", &dtype,
		"copy?", &kopy,
	); err != nil {
		return nil, err
	}
	if len(args) > 0 {
		return nil, fmt.Errorf("dataframe: unexpected positional arguments")
	}

	return &DataFrame{}, nil
}

// String implements the Stringer interface.
func (DataFrame) String() string { return "TODO - dataframe string rendering" }

// Type returns a short string describing the value's type.
func (DataFrame) Type() string { return fmt.Sprintf("%s.DataFrame", Name) }

// Freeze renders DataFrame immutable. required by starlark.Value interface
func (df *DataFrame) Freeze() { df.frozen = true }

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (df *DataFrame) Hash() (uint32, error) {
	// TODO (b5) - finish
	return 0, nil
}

// Truth reports whether the DataFrame is non-zero.
func (df *DataFrame) Truth() starlark.Bool {
	// TODO (b5) - finish
	return true
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (df *DataFrame) Attr(name string) (starlark.Value, error) {
	return nil, fmt.Errorf("unrecognized %s attribute %q", df.Type(), name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (df *DataFrame) AttrNames() []string {
	return []string{}
}

// CompareSameType implements comparison of two DataFrame values. required by
// starlark.Comparable interface.
func (df *DataFrame) CompareSameType(op syntax.Token, v starlark.Value, depth int) (bool, error) {
	return false, unfinishedError(df, "CompareSameType")
	// cmp := 0
	// if x, y := d, v.(DataFrame); x < y {
	// 	cmp = -1
	// } else if x > y {
	// 	cmp = 1
	// }
	// return threeway(op, cmp), nil
}

// // Threeway interprets a three-way comparison value cmp (-1, 0, +1)
// // as a boolean comparison (e.g. x < y).
// func threeway(op syntax.Token, cmp int) bool {
// 	switch op {
// 	case syntax.EQL:
// 		return cmp == 0
// 	case syntax.NEQ:
// 		return cmp != 0
// 	case syntax.LE:
// 		return cmp <= 0
// 	case syntax.LT:
// 		return cmp < 0
// 	case syntax.GE:
// 		return cmp >= 0
// 	case syntax.GT:
// 		return cmp > 0
// 	}
// 	panic(op)
// }

// Binary implements binary operators, which satisfies the starlark.HasBinary
// interface
func (df *DataFrame) Binary(op syntax.Token, y starlark.Value, side starlark.Side) (starlark.Value, error) {
	return nil, unfinishedError(df, "Binary")
}

func (df *DataFrame) Get(key starlark.Value) (value starlark.Value, found bool, err error) {
	return nil, false, unfinishedError(df, "Get")
}
