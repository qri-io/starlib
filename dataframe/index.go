package dataframe

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type Index struct {
	name   string
	data   *starlark.List
	frozen bool
}

var (
	_ starlark.Value = (*Index)(nil)
)

func newIndex(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		data                starlark.Value
		name, dtype         starlark.String
		kopy, tupleize_cols starlark.Bool
	)
	if err := starlark.UnpackArgs("Index", args, kwargs,
		"data?", &data,
		"dtype?", &dtype,
		"copy?", &kopy,
		"name?", &name,
		"tupleize_cols?", &tupleize_cols,
	); err != nil {
		return nil, err
	}
	if len(args) > 1 {
		return nil, fmt.Errorf("%s.Index: unexpected positional arguments", Name)
	}

	d, ok := data.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("%s.Index: data arg must be a List", Name)
	}

	return &Index{
		name: maybeExtractName(name, data),
		data: d,
	}, nil
}

// String implements the Stringer interface.
func (idx *Index) String() string { return fmt.Sprintf("Index(%s)", toString(idx.data)) }

// Type returns a short string describing the value's type.
func (Index) Type() string { return fmt.Sprintf("%s.Index", Name) }

// Freeze renders Index immutable. required by starlark.Value interface
// indexes are immitubable
func (idx *Index) Freeze() {
	idx.frozen = true
	idx.data.Freeze()
}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (idx *Index) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s.Index", Name)
}

// Truth reports whether the Index is non-zero.
func (idx *Index) Truth() starlark.Bool { return idx.Len() > 0 }

func (idx *Index) Len() int { return idx.data.Len() }

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (idx *Index) Attr(name string) (starlark.Value, error) {
	if m := idx.method(name); m != nil {
		return m()
	}
	return nil, fmt.Errorf("unrecognized %s attribute %q", idx.Type(), name)
}

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (idx *Index) AttrNames() []string {
	return idx.methodNames()
}

func (idx *Index) methodNames() []string {
	return []string{
		"eq",
	}
}

func (idx *Index) method(name string) func() (starlark.Value, error) {
	switch name {
	case "eq":
		return func() (starlark.Value, error) {
			// // Allocate a closure over 'method'.
			// impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			// 	return method(b.Name(), b.Receiver(), args, kwargs)
			// }
			return starlark.NewBuiltin(name, idx.eq).BindReceiver(idx), nil
		}
	}
	return nil // no such method
}

func (idx *Index) eq(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Value
	if err := starlark.UnpackPositionalArgs("eq", args, kwargs, 1, &x); err != nil {
		return nil, err
	}

	switch xt := x.(type) {
	case starlark.String:
		res := makeFalseIndex(idx.Len())
		i := 0
		var v starlark.Value
		it := idx.data.Iterate()
		for it.Next(&v) {
			if eq, err := starlark.Compare(syntax.EQL, v, xt); err != nil {
				return nil, err
			} else if eq {
				res.data.SetIndex(i, starlark.True)
			}
			i++
		}
		it.Done()
		return res, nil
	}

	return nil, unfinishedError(idx, fmt.Sprintf("eq: %s", x.Type()))
}

// CompareSameType implements comparison of two Index values. required by
// starlark.Comparable interface.
func (idx *Index) CompareSameType(op syntax.Token, y starlark.Value, depth int) (bool, error) {
	return starlark.Compare(op, idx.data, y.(*Index).data)
}

// Binary implements binary operators, which satisfies the starlark.HasBinary
// interface
func (idx *Index) Binary(op syntax.Token, y starlark.Value, side starlark.Side) (starlark.Value, error) {
	return nil, unfinishedError(idx, "Binary")
}

func makeFalseIndex(length int) *Index {
	elems := make([]starlark.Value, 0, length)
	for i := 0; i < length; i++ {
		elems = append(elems, starlark.False)
	}

	return &Index{
		data: starlark.NewList(elems),
	}
}

func (df *Index) Get(key starlark.Value) (value starlark.Value, found bool, err error) {
	return nil, false, unfinishedError(df, "Get")
}

func maybeExtractName(name starlark.String, data starlark.Value) string {
	// TODO(b5) - finish
	return name.GoString()
}
