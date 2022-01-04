package dataframe

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	gotime "time"

	"go.starlark.net/lib/time"
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

func newTypedSliceBuilderNaNFilled(numRows int) *typedSliceBuilder {
	builder := newTypedSliceBuilder(numRows)
	for i := 0; i < numRows; i++ {
		builder.push(math.NaN())
	}
	return builder
}

func (t *typedSliceBuilder) setType(dtype string) {
	if dtype == "" {
		return
	}
	t.dType = dtype
	t.currType = dtype
	if t.dType == "int64" || t.dType == "bool" {
		t.whichVals = typeInt
	} else if t.dType == "datetime64[ns]" {
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
		} else if tim, ok := val.(time.Time); ok {
			t.currType = "datetime64[ns]"
			t.dType = "datetime64[ns]"
			t.whichVals = typeInt
			val = timeToInt(tim)
		} else if val == nil {
			t.currType = "float64"
			t.whichVals = typeFloat
			val = math.NaN()
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
				// Unknown conversion, just use objects
				t.coerceToObjects()
			}
		} else if num, ok := val.(int64); ok {
			if t.currType == "float64" {
				val = float64(num)
			} else if t.currType == "object" {
				// no need to convert
			} else if t.currType != "int64" {
				// Unknown conversion, just use objects
				t.coerceToObjects()
			}
		} else if _, ok := val.(float64); ok {
			// TODO(dustmop): If t.dType != "", is this an error?
			if t.currType == "int64" && t.dType == "" {
				// The list was ints, found a float, coerce the previous list to floats
				t.currType = "float64"
				t.whichVals = typeFloat
				t.valFloats = convertIntsToFloats(t.valInts)
			} else if t.currType == "object" {
				// no need to convert
			} else if t.currType != "float64" {
				// Unknown conversion, just use objects
				t.coerceToObjects()
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
			} else if t.currType == "bool" {
				// The list was bools, found a string, coerce the previous list to objects
				t.currType = "object"
				t.whichVals = typeObj
				t.valObjs = convertBoolsToObjects(t.valInts)
			} else if t.currType == "datetime64[ns]" {
				// no need to convert
				timestamp, err := gotime.Parse("2006-01-02 15:04:05", text)
				if err != nil {
					// TODO(dustmop): Add test
					t.buildError = fmt.Errorf("could not parse timestamp from %s: %w", text, err)
				}
				val = timestamp.UnixNano()
			} else if t.currType != "object" {
				// Unknown conversion, just use objects
				t.coerceToObjects()
			}
		} else if tim, ok := val.(time.Time); ok {
			if t.currType == "datetime64[ns]" {
				val = timeToInt(tim)
			} else {
				// TODO(dustmop): Fix this
				t.buildError = fmt.Errorf("cannot append timestamp to list of type %s", t.currType)
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
				// Unknown conversion, just use objects
				t.coerceToObjects()
			}
		} else if val == nil {
			if t.currType == "float64" {
				val = math.NaN()
			} else if t.currType == "object" {
				// no need to convert
			} else {
				if t.whichVals == typeInt {
					if t.currType == "bool" {
						t.valObjs = convertBoolsToObjects(t.valInts)
						t.whichVals = typeObj
						t.currType = "object"
					} else {
						t.valFloats = convertIntsToFloats(t.valInts)
						t.whichVals = typeFloat
						t.currType = "float64"
						val = math.NaN()
					}
				}
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
	if num, err := strconv.ParseInt(text, 10, 64); err == nil {
		t.push(num)
		return
	} else if f, err := strconv.ParseFloat(text, 64); err == nil {
		t.push(f)
		return
	} else if b, err := strconv.ParseBool(text); err == nil {
		t.push(b)
		return
	}
	t.push(text)
}

func (t *typedSliceBuilder) error() error {
	return t.buildError
}

func (t *typedSliceBuilder) keys() []string {
	return t.keyList
}

func (t *typedSliceBuilder) coerceToObjects() {
	if t.whichVals == typeInt {
		if t.currType == "bool" {
			t.valObjs = convertBoolsToObjects(t.valInts)
		} else {
			t.valObjs = convertIntsToObjects(t.valInts)
		}
	} else if t.whichVals == typeFloat {
		t.valObjs = convertFloatsToObjects(t.valFloats)
	}
	t.currType = "object"
	t.whichVals = typeObj
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

func (t *typedSliceBuilder) Len() int {
	if t.whichVals == typeInt {
		return len(t.valInts)
	} else if t.whichVals == typeFloat {
		return len(t.valFloats)
	}
	return len(t.valObjs)
}
