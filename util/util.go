package util

import (
	"fmt"
	"strconv"

	"github.com/google/skylark"
)

// AsString unquotes a skylark string value
func AsString(x skylark.Value) (string, error) {
	return strconv.Unquote(x.String())
}

// Unmarshal decodes a skylark.Value into it's golang counterpart
func Unmarshal(x skylark.Value) (val interface{}, err error) {
	switch x.Type() {
	case "NoneType":
		val = nil
	case "bool":
		val = x.Truth() == skylark.True
	case "int":
		val, err = skylark.AsInt32(x)
	case "float":
		if f, ok := skylark.AsFloat(x); ok {
			val = f
		} else {
			err = fmt.Errorf("couldn't parse float")
		}
	case "string":
		val, err = AsString(x)
		// val = x.String()
	case "dict":
		if dict, ok := x.(*skylark.Dict); ok {
			var (
				v     skylark.Value
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
				str, err = AsString(k)
				if err != nil {
					return
				}

				value[str] = pval
			}
			val = value
		} else {
			err = fmt.Errorf("error parsing dict. invalid type: %v", x)
		}
	case "list":
		if list, ok := x.(*skylark.List); ok {
			var (
				i     int
				v     skylark.Value
				iter  = list.Iterate()
				value = make([]interface{}, list.Len())
			)

			for iter.Next(&v) {
				value[i], err = Unmarshal(v)
				if err != nil {
					return
				}
				i++
			}
			iter.Done()
			val = value
		} else {
			err = fmt.Errorf("error parsing list. invalid type: %v", x)
		}
	case "tuple":
		if tuple, ok := x.(skylark.Tuple); ok {
			var (
				i     int
				v     skylark.Value
				iter  = tuple.Iterate()
				value = make([]interface{}, tuple.Len())
			)

			for iter.Next(&v) {
				value[i], err = Unmarshal(v)
				if err != nil {
					return
				}
				i++
			}
			iter.Done()
			val = value
		} else {
			err = fmt.Errorf("error parsing dict. invalid type: %v", x)
		}
	case "set":
		fmt.Println("errnotdone: SET")
		err = fmt.Errorf("sets aren't yet supported")
	default:
		fmt.Println("errbadtype:", x.Type())
		err = fmt.Errorf("unrecognized skylark type: %s", x.Type())
	}
	return
}

// Marshal turns go values into skylark types
func Marshal(data interface{}) (v skylark.Value, err error) {
	switch x := data.(type) {
	case nil:
		v = skylark.None
	case bool:
		v = skylark.Bool(x)
	case string:
		v = skylark.String(x)
	case int:
		v = skylark.MakeInt(x)
	case int64:
		v = skylark.MakeInt(int(x))
	case float64:
		v = skylark.Float(x)
	case []interface{}:
		var elems = make([]skylark.Value, len(x))
		for i, val := range x {
			elems[i], err = Marshal(val)
			if err != nil {
				return
			}
		}
		v = skylark.NewList(elems)
	case map[string]interface{}:
		dict := &skylark.Dict{}
		var elem skylark.Value
		for key, val := range x {
			elem, err = Marshal(val)
			if err != nil {
				return
			}
			if err = dict.Set(skylark.String(key), elem); err != nil {
				return
			}
		}
		v = dict
	default:
		return skylark.None, fmt.Errorf("unrecognized type: %#v", x)
	}
	return
}
