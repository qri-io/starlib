package dataframe

import (
	"encoding/json"
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

func marshalRowToString(row []string) string {
	data, err := json.Marshal(row)
	if err != nil {
		return "?"
	}
	return string(data)
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
		result = append(result, fmt.Sprintf("%1.1f", f))
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
