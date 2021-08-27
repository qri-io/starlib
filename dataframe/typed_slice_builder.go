package dataframe

import (
	"fmt"
	"math"
	"reflect"
)

type typedSliceBuilder struct {
	capHint    int
	keyList    []string
	valInts    []int
	valFloats  []float64
	valObjs    []interface{}
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

func newTypedSliceBuilder(capacityHint int) *typedSliceBuilder {
	return &typedSliceBuilder{
		capHint: capacityHint,
	}
}

func newTypedSliceBuilderFromSeries(series *Series) *typedSliceBuilder {
	return &typedSliceBuilder{
		capHint:   series.Len(),
		whichVals: series.which,
		valInts:   series.valInts,
		valFloats: series.valFloats,
		valObjs:   series.valObjs,
		dType:     series.dtype,
	}
}

func (t *typedSliceBuilder) setType(dtype string) {
	if dtype == "" {
		return
	}
	t.dType = dtype
	t.currType = dtype
	if t.dType == "int64" || t.dType == "bool" {
		t.whichVals = typeInt
	} else if t.dType == "float64" {
		t.whichVals = typeFloat
	} else if t.dType == "object" || t.dType == "datetime64[ns]" {
		t.whichVals = typeObj
	} else if t.dType == "string" {
		// TODO(dustmop): Is string a real type for pandas?
		t.currType = "object"
		t.whichVals = typeObj
	} else {
		t.buildError = fmt.Errorf("invalid dtype: %q", dtype)
	}
}

func (t *typedSliceBuilder) push(val interface{}) {
	if t.currType == "" {
		// Initial data type
		if num, ok := val.(int); ok {
			t.currType = "int64"
			t.whichVals = typeInt
			_ = num
		} else if num, ok := val.(int64); ok {
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
				// no need to convert
			} else if t.currType != "int64" {
				t.buildError = fmt.Errorf("coercion failed, int: %v to %q", num, t.currType)
				return
			}
		} else if num, ok := val.(int64); ok {
			if t.currType == "float64" {
				val = float64(num)
			} else if t.currType == "object" {
				// no need to convert
			} else if t.currType != "int64" {
				t.buildError = fmt.Errorf("coercion failed, int64: %v to %q", num, t.currType)
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
				// no need to convert
			} else if t.currType != "float64" {
				t.buildError = fmt.Errorf("coercion failed, float64: %v to %q", f, t.currType)
				return
			}
		} else if text, ok := val.(string); ok {
			if t.currType == "int64" && t.dType == "" {
				// The list was ints, found a string, coerce the previous list to objects
				t.currType = "object"
				t.whichVals = typeObj
				t.valObjs = convertIntsToObjects(t.valInts)
			} else if t.currType == "float64" && t.dType == "" {
				// The list was floats, found a string, coerce the previous list to objects
				t.currType = "object"
				t.whichVals = typeObj
				t.valObjs = convertFloatsToObjects(t.valFloats)
			} else if t.currType == "datetime64[ns]" {
				// no need to convert
			} else if t.currType != "object" {
				t.buildError = fmt.Errorf("coercion failed, string: %v to %q", text, t.currType)
				return
			}
		} else if b, ok := val.(bool); ok {
			if t.currType == "bool" {
				val = 0
				if b {
					val = 1
				}
			} else if t.currType == "object" {
				// no need to convert
			} else {
				t.buildError = fmt.Errorf("coercion failed, bool: %v to %q", b, t.currType)
				return
			}
		} else if val == nil {
			if t.currType == "float64" {
				val = math.NaN()
			} else {
				if t.whichVals == typeInt {
					if t.currType == "bool" {
						t.valObjs = convertBoolsToObjects(t.valInts)
					} else {
						t.valObjs = convertIntsToObjects(t.valInts)
					}
				}
				t.whichVals = typeObj
				t.currType = "object"
			}
		} else {
			t.buildError = fmt.Errorf("invalid object %v of type %s", val, reflect.TypeOf(val))
			return
		}
	}

	// Add to the appropriate array
	if t.whichVals == typeInt {
		if t.valInts == nil {
			t.valInts = make([]int, 0, t.capHint)
		}
		if n, ok := val.(int); ok {
			t.valInts = append(t.valInts, n)
		} else if n, ok := val.(int64); ok {
			t.valInts = append(t.valInts, int(n))
		} else if val == nil {
			t.whichVals = typeObj
			t.valObjs = convertIntsToObjects(t.valInts)
			t.valObjs = append(t.valObjs, val)
		} else {
			t.buildError = fmt.Errorf("wanted int, got %v of type %s", val, reflect.TypeOf(val))
		}
	} else if t.whichVals == typeFloat {
		if t.valFloats == nil {
			t.valFloats = make([]float64, 0, t.capHint)
		}
		t.valFloats = append(t.valFloats, val.(float64))
	} else if t.whichVals == typeObj {
		if t.valObjs == nil {
			t.valObjs = make([]interface{}, 0, t.capHint)
		}
		t.valObjs = append(t.valObjs, val)
	}
}

func (t *typedSliceBuilder) pushNil() {
	if t.whichVals == typeInt {
		t.valInts = append(t.valInts, 0)
	} else if t.whichVals == typeFloat {
		t.valFloats = append(t.valFloats, 0.0)
	} else if t.whichVals == typeObj {
		t.valObjs = append(t.valObjs, nil)
	}
}

func (t *typedSliceBuilder) pushKeyVal(key string, val interface{}) {
	t.keyList = append(t.keyList, key)
	t.push(val)
}

func (t *typedSliceBuilder) parsePush(text string) {
	// Parse scalar from the text, push it. Used for csv reader.
	// TODO: Actually parse int and float from text, don't assume it must be string.
	t.push(text)
}

func (t *typedSliceBuilder) error() error {
	return t.buildError
}

func (t *typedSliceBuilder) keys() []string {
	return t.keyList
}

// toSeries returns the series that has been built
func (t *typedSliceBuilder) toSeries(index *Index, name string) Series {
	dtype := t.dType
	if dtype == "" && t.currType != "" {
		dtype = t.currType
	}
	return Series{
		dtype:     dtype,
		which:     t.whichVals,
		valInts:   t.valInts,
		valFloats: t.valFloats,
		valObjs:   t.valObjs,
		index:     index,
		name:      name,
	}
}
