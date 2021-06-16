package dataframe

import (
	"fmt"
	"strconv"

	"go.starlark.net/starlark"
)

// convert starlark value to a string fit for printing
func toStr(val starlark.Value) string {
	if text, ok := val.(starlark.String); ok {
		return string(text)
	}
	if val == nil {
		return "<nil>"
	}
	return val.String()
}

// convert starlark value to a string if it has string type, and the empty string otherwise
func toStrOrEmpty(v starlark.Value) string {
	if text, ok := v.(starlark.String); ok {
		return string(text)
	}
	return ""
}

// convert starlark value to a list of strings, fit for printing, or nil if not possible
func toStrListOrNil(v starlark.Value) []string {
	if v == nil {
		return nil
	}
	if v == starlark.None {
		return nil
	}

	list, ok := v.(*starlark.List)
	if ok {
		result := make([]string, 0, list.Len())
		for i := 0; i < list.Len(); i++ {
			result = append(result, toStr(list.Index(i)))
		}
		return result
	}

	tup, ok := v.(starlark.Tuple)
	if ok {
		result := make([]string, 0, tup.Len())
		for i := 0; i < tup.Len(); i++ {
			result = append(result, toStr(tup.Index(i)))
		}
		return result
	}

	return nil
}

// convert starlark value to a list of any values, or nil if not a list
func toInterfaceListOrNil(v starlark.Value) []interface{} {
	if v == nil {
		return nil
	}
	if v == starlark.None {
		return nil
	}

	list, ok := v.(*starlark.List)
	if ok {
		result := make([]interface{}, 0, list.Len())
		for i := 0; i < list.Len(); i++ {
			elem, ok := toScalarMaybe(list.Index(i))
			if !ok {
				return nil
			}
			result = append(result, elem)
		}
		return result
	}

	tup, ok := v.(starlark.Tuple)
	if ok {
		result := make([]interface{}, 0, tup.Len())
		for i := 0; i < tup.Len(); i++ {
			elem, ok := toScalarMaybe(tup.Index(i))
			if !ok {
				return nil
			}
			result = append(result, elem)
		}
		return result
	}

	return nil
}

// convert starlark value to a go native int if it has the right type
func toIntMaybe(v starlark.Value) (int, bool) {
	n, err := starlark.AsInt32(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// convert starlark value to a go native float if it has the right type
func toFloatMaybe(v starlark.Value) (float64, bool) {
	return starlark.AsFloat(v)
}

// convert starlark value to a go native string if it has the right type
func toStrMaybe(v starlark.Value) (string, bool) {
	if str, ok := v.(starlark.String); ok {
		return string(str), true
	}
	return "", false
}

// convert starlark value to a go native bool if it has the right type
func toBoolMaybe(v starlark.Value) (bool, bool) {
	if b, ok := v.(starlark.Bool); ok {
		return bool(b), true
	}
	return false, false
}

// convert starlark value to go native int, bool, float, or string
func toScalarMaybe(v starlark.Value) (interface{}, bool) {
	if num, ok := toIntMaybe(v); ok {
		return num, true
	}
	if f, ok := toFloatMaybe(v); ok {
		return f, true
	}
	if text, ok := toStrMaybe(v); ok {
		return text, true
	}
	if b, ok := toBoolMaybe(v); ok {
		return b, true
	}
	return nil, false
}

func toIndexMaybe(v starlark.Value) (*Index, bool) {
	texts := toStrListOrNil(v)
	if texts != nil {
		return NewIndex(texts, ""), true
	}
	if index, ok := v.(*Index); ok {
		return index, true
	}
	return nil, false
}

func stringifyFloat(f float64) string {
	return fmt.Sprintf("%1.1f", f)
}

// convert a list of ints to a list of floats
func convertIntsToFloats(vals []int) []float64 {
	result := make([]float64, 0, len(vals))
	for _, n := range vals {
		result = append(result, float64(n))
	}
	return result
}

// convert a list of ints to a list of strings
func convertIntsToStrings(vals []int) []string {
	result := make([]string, 0, len(vals))
	for _, n := range vals {
		result = append(result, strconv.Itoa(n))
	}
	return result
}

// convert a list of floats to a list of strings
func convertFloatsToStrings(vals []float64) []string {
	result := make([]string, 0, len(vals))
	for _, f := range vals {
		result = append(result, stringifyFloat(f))
	}
	return result
}

// convert one of the supported go native data types into a starlark value
func convertToStarlark(it interface{}) (starlark.Value, error) {
	if num, ok := it.(int); ok {
		return starlark.MakeInt(num), nil
	}
	if f, ok := it.(float64); ok {
		return starlark.Float(f), nil
	}
	if str, ok := it.(string); ok {
		return starlark.String(str), nil
	}
	return starlark.None, fmt.Errorf("unknown type of %v", it)
}

func coerceToDatatype(text, dtype string) string {
	if dtype == "bool" {
		if text == "1" {
			return " True"
		}
		return "False"
	}
	return text
}
