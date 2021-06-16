package dataframe

import (
	"fmt"
	"reflect"
	"strconv"
)

type typedArrayBuilder struct {
	size       int
	keyList    []string
	valInts    []int
	valFloats  []float64
	valObjs    []string
	whichVals  int
	dType      string
	currType   string
	buildError error
}

// A Series contains values of one of these three types. The which field uses these
// constants to determine which slice holds the actual values of the Series.
const (
	typeInt   = 1
	typeFloat = 2
	typeObj   = 3
)

func newTypedArrayBuilder(size int) *typedArrayBuilder {
	return &typedArrayBuilder{
		size: size,
	}
}

func (t *typedArrayBuilder) setType(dtype string) {
	if dtype == "" {
		return
	}
	t.dType = dtype
	t.currType = dtype
	if t.dType == "int64" || t.dType == "bool" {
		t.whichVals = typeInt
	} else if t.dType == "float64" {
		t.whichVals = typeFloat
	} else if t.dType == "object" {
		t.whichVals = typeObj
	} else if t.dType == "string" {
		// TODO(dustmop): Is string a real type for pandas?
		t.currType = "object"
		t.whichVals = typeObj
	} else {
		t.buildError = fmt.Errorf("invalid dtype: %q", dtype)
	}
}

func (t *typedArrayBuilder) push(val interface{}) {
	if t.currType == "" {
		// Initial data type
		if num, ok := val.(int); ok {
			t.currType = "int64"
			t.whichVals = typeInt
			_ = num
		} else if f, ok := val.(float64); ok {
			t.currType = "float64"
			t.whichVals = typeFloat
			_ = f
		} else if text, ok := val.(string); ok {
			t.currType = "object"
			t.whichVals = typeObj
			_ = text
		} else if b, ok := val.(bool); ok {
			t.currType = "bool"
			t.whichVals = typeInt
			val = 0
			if b {
				val = 1
			}
		} else {
			t.buildError = fmt.Errorf("invalid object %v of type %s", val, reflect.TypeOf(val))
			return
		}
	} else {
		// Coerce types as needed
		if num, ok := val.(int); ok {
			if t.currType == "float64" {
				val = float64(num)
			} else if t.currType == "object" {
				val = strconv.Itoa(num)
			} else if t.currType != "int64" {
				t.buildError = fmt.Errorf("handle coercion: %v to %q", num, t.currType)
				return
			}
		} else if f, ok := val.(float64); ok {
			// TODO(dustmop): If t.dType != "", is this an error?
			if t.currType == "int64" && t.dType == "" {
				// The list was ints, found a float, coerce the previous list to floats
				t.currType = "float64"
				t.whichVals = typeFloat
				t.valFloats = convertIntsToFloats(t.valInts)
			} else if t.currType == "object" {
				//
				val = stringifyFloat(f)
			} else if t.currType != "float64" {
				t.buildError = fmt.Errorf("handle coercion: %v to %q", f, t.currType)
				return
			}
		} else if text, ok := val.(string); ok {
			if t.currType == "int64" && t.dType == "" {
				// The list was ints, found a string, coerce the previous list to objects
				t.currType = "object"
				t.whichVals = typeObj
				t.valObjs = convertIntsToStrings(t.valInts)
			} else if t.currType == "float64" && t.dType == "" {
				// The list was floats, found a string, coerce the previous list to objects
				t.currType = "object"
				t.whichVals = typeObj
				t.valObjs = convertFloatsToStrings(t.valFloats)
			} else if t.currType != "object" {
				t.buildError = fmt.Errorf("handle coercion: %v to %q", text, t.currType)
				return
			}
		} else if b, ok := val.(bool); ok {
			if t.currType == "bool" {
				val = 0
				if b {
					val = 1
				}
			} else if t.currType != "bool" {
				t.buildError = fmt.Errorf("handle coercion: %v to %q", text, t.currType)
				return
			}
		} else {
			t.buildError = fmt.Errorf("invalid object %v of type %s", val, reflect.TypeOf(val))
			return
		}
	}

	// Add to the appropriate array
	if t.whichVals == typeInt {
		if t.valInts == nil {
			t.valInts = make([]int, 0, t.size)
		}
		t.valInts = append(t.valInts, val.(int))
	} else if t.whichVals == typeFloat {
		if t.valFloats == nil {
			t.valFloats = make([]float64, 0, t.size)
		}
		t.valFloats = append(t.valFloats, val.(float64))
	} else if t.whichVals == typeObj {
		if t.valObjs == nil {
			t.valObjs = make([]string, 0, t.size)
		}
		t.valObjs = append(t.valObjs, val.(string))
	}
}

func (t *typedArrayBuilder) pushKeyVal(key string, val interface{}) {
	t.keyList = append(t.keyList, key)
	t.push(val)
}

func (t *typedArrayBuilder) parsePush(text string) {
	// TODO: Parse a scalar from the text, push it. Used for csv reader.
}

func (t *typedArrayBuilder) error() error {
	return t.buildError
}

func (t *typedArrayBuilder) keys() []string {
	return t.keyList
}

func (t *typedArrayBuilder) dtype() string {
	if t.dType == "" {
		if t.currType != "" {
			t.dType = t.currType
		}
	}
	return t.dType
}

func (t *typedArrayBuilder) which() int {
	return t.whichVals
}

func (t *typedArrayBuilder) asIntSlice() []int {
	return t.valInts
}

func (t *typedArrayBuilder) asFloatSlice() []float64 {
	return t.valFloats
}

func (t *typedArrayBuilder) asObjSlice() []string {
	return t.valObjs
}
