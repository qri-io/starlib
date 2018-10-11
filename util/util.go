package util

import (
	"fmt"
	"strconv"

	starlark "github.com/google/skylark"
)

// asString unquotes a starlark string value
func asString(x starlark.Value) (string, error) {
	return strconv.Unquote(x.String())
}

// Unmarshal decodes a starlark.Value into it's golang counterpart
func Unmarshal(x starlark.Value) (val interface{}, err error) {
	switch x.Type() {
	case "NoneType":
		val = nil
	case "bool":
		val = x.Truth() == starlark.True
	case "int":
		val, err = starlark.AsInt32(x)
	case "float":
		if f, ok := starlark.AsFloat(x); !ok {
			err = fmt.Errorf("couldn't parse float")
		} else {
			val = f
		}
	case "string":
		val, err = asString(x)
		// val = x.String()
	case "dict":
		dict, ok := x.(*starlark.Dict)
		if !ok {
			err = fmt.Errorf("error parsing dict. invalid type: %v", x)
			return
		}

		var (
			v     starlark.Value
			pval  interface{}
			value = map[string]interface{}{}
		)

		for _, k := range dict.Keys() {
			v, ok, err = dict.Get(k)
			if err != nil {
				return
			}

			pval, err = Unmarshal(v)
			if err != nil {
				return
			}

			var str string
			str, err = asString(k)
			if err != nil {
				return
			}

			value[str] = pval
		}
		val = value
	case "list":
		list, ok := x.(*starlark.List)
		if !ok {
			err = fmt.Errorf("error parsing list. invalid type: %v", x)
			return
		}

		var (
			i     int
			v     starlark.Value
			iter  = list.Iterate()
			value = make([]interface{}, list.Len())
		)

		defer iter.Done()
		for iter.Next(&v) {
			value[i], err = Unmarshal(v)
			if err != nil {
				return
			}
			i++
		}
		val = value
	case "tuple":
		tuple, ok := x.(starlark.Tuple)
		if !ok {
			err = fmt.Errorf("error parsing dict. invalid type: %v", x)
			return
		}

		var (
			i     int
			v     starlark.Value
			iter  = tuple.Iterate()
			value = make([]interface{}, tuple.Len())
		)

		defer iter.Done()
		for iter.Next(&v) {
			value[i], err = Unmarshal(v)
			if err != nil {
				return
			}
			i++
		}
		val = value
	case "set":
		fmt.Println("errnotdone: SET")
		err = fmt.Errorf("sets aren't yet supported")
	default:
		fmt.Println("errbadtype:", x.Type())
		err = fmt.Errorf("unrecognized starlark type: %s", x.Type())
	}
	return
}

// Marshal turns go values into starlark types
func Marshal(data interface{}) (v starlark.Value, err error) {
	switch x := data.(type) {
	case nil:
		v = starlark.None
	case bool:
		v = starlark.Bool(x)
	case string:
		v = starlark.String(x)
	case int:
		v = starlark.MakeInt(x)
	case int64:
		v = starlark.MakeInt64(int64(x))
	case float64:
		v = starlark.Float(x)
	case []interface{}:
		var elems = make([]starlark.Value, len(x))
		for i, val := range x {
			elems[i], err = Marshal(val)
			if err != nil {
				return
			}
		}
		v = starlark.NewList(elems)
	case map[string]interface{}:
		dict := &starlark.Dict{}
		var elem starlark.Value
		for key, val := range x {
			elem, err = Marshal(val)
			if err != nil {
				return
			}
			if err = dict.Set(starlark.String(key), elem); err != nil {
				return
			}
		}
		v = dict
	default:
		return starlark.None, fmt.Errorf("unrecognized type: %#v", x)
	}
	return
}
