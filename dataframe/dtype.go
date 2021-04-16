package dataframe

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// toString returns the string form of value v.
// It may be more efficient than v.String() for larger values.
func toString(v starlark.Value) string {
	buf := new(strings.Builder)
	writeValue(buf, v, nil)
	return buf.String()
}

// writeValue writes x to out.
//
// path is used to detect cycles.
// It contains the list of *List and *Dict values we're currently printing.
// (These are the only potentially cyclic structures.)
// Callers should generally pass nil for path.
// It is safe to re-use the same path slice for multiple calls.
func writeValue(out *strings.Builder, x starlark.Value, path []starlark.Value) {
	switch x := x.(type) {
	case nil:
		out.WriteString("<nil>") // indicates a bug

	// These four cases are duplicates of T.String(), for efficiency.
	case starlark.NoneType:
		out.WriteString("None")

	case starlark.Int:
		out.WriteString(x.String())

	case starlark.Bool:
		if x {
			out.WriteString("True")
		} else {
			out.WriteString("False")
		}

	case starlark.String:
		out.WriteString(syntax.Quote(string(x), false))

	case starlark.Tuple:
		out.WriteByte('(')
		for i, elem := range x {
			if i > 0 {
				out.WriteString(", ")
			}
			writeValue(out, elem, path)
		}
		if len(x) == 1 {
			out.WriteByte(',')
		}
		out.WriteByte(')')

	case *starlark.Function:
		fmt.Fprintf(out, "<function %s>", x.Name())

	case *starlark.Builtin:
		fmt.Fprintf(out, "<built-in function %s>", x.Name())

	case *starlark.Dict:
		out.WriteByte('{')
		if pathContains(path, x) {
			out.WriteString("...") // dict contains itself
		} else {
			sep := ""
			for _, item := range x.Items() {
				k, v := item[0], item[1]
				out.WriteString(sep)
				writeValue(out, k, path)
				out.WriteString(": ")
				writeValue(out, v, append(path, x)) // cycle check
				sep = ", "
			}
		}
		out.WriteByte('}')

	default:
		out.WriteString(x.String())
	}
}

func pathContains(path []starlark.Value, x starlark.Value) bool {
	for _, y := range path {
		if x == y {
			return true
		}
	}
	return false
}
