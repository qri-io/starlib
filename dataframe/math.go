package dataframe

import (
	"fmt"
	"math"

	"go.starlark.net/starlark"
)

func mathAbs(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var subject starlark.Value

	if err := starlark.UnpackArgs("abs", args, kwargs,
		"subject", &subject,
	); err != nil {
		return nil, err
	}

	series, ok := subject.(*Series)
	if !ok {
		return starlark.None, fmt.Errorf("dataframe.abs can only be called on a Series")
	}

	builder := newTypedSliceBuilder(series.Len())
	if series.which == typeInt {
		// Fast math, all int series
		for i := 0; i < series.Len(); i++ {
			n := series.valInts[i]
			if n < 0 {
				n = -n
			}
			builder.push(n)
		}
	} else {
		for i := 0; i < series.Len(); i++ {
			builder.push(math.Abs(series.FloatAt(i)))
		}
	}
	res := builder.toSeries(series.index, series.name)
	return &res, nil
}
