package dataframe

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math"
	"testing"
)

// Test adding 3 ints to get an int series
func TestBuildIntSlice(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(123)
	builder.push(456)
	builder.push(789)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "int64" {
		t.Errorf("expected dtype == int64, got %q", series.dtype)
	}
	if series.which != typeInt {
		t.Errorf("expected which == int, got %v", series.which)
	}
	expect := []interface{}{123, 456, 789}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	n, ok := obj.(int)
	if !ok {
		t.Fatalf("series.At should return a int, got %v", obj)
	}
	if n != 123 {
		t.Errorf("series.At(0) expected 123, got %v", n)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    123
1    456
2    789
dtype: int64`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding 3 floats to get a float series
func TestBuildFloatSlice(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(12.3)
	builder.push(45.6)
	builder.push(78.9)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "float64" {
		t.Errorf("expected dtype == float64, got %q", series.dtype)
	}
	if series.which != typeFloat {
		t.Errorf("expected which == float, got %v", series.which)
	}
	expect := []interface{}{12.3, 45.6, 78.9}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	f, ok := obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float64, got %v", obj)
	}
	if f != 12.3 {
		t.Errorf("series.At(0) expected 12.3, got %v", f)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    12.3
1    45.6
2    78.9
dtype: float64`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding 3 bools to get a bool series
func TestBuildBoolSlice(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(true)
	builder.push(false)
	builder.push(true)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "bool" {
		t.Errorf("expected dtype == bool, got %q", series.dtype)
	}
	if series.which != typeInt {
		t.Errorf("expected which == int, got %v", series.which)
	}
	expect := []interface{}{1, 0, 1}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	b, ok := obj.(bool)
	if !ok {
		t.Fatalf("series.At should return a bool, got %v", obj)
	}
	if !b {
		t.Errorf("series.At(0) expected true")
	}
	obj = series.At(1)
	b, ok = obj.(bool)
	if !ok {
		t.Fatalf("series.At should return a bool, got %v", obj)
	}
	if b {
		t.Errorf("series.At(1) expected false")
	}

	// Verify the string representation
	text := series.String()
	expectText := `0     True
1    False
2     True
dtype: bool`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding 3 strings to get an object series
func TestBuildStrObjSlice(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push("cat")
	builder.push("dog")
	builder.push("eel")
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
	if series.which != typeObj {
		t.Errorf("expected which == obj, got %v", series.which)
	}
	expect := []interface{}{"cat", "dog", "eel"}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	s, ok := obj.(string)
	if !ok {
		t.Fatalf("series.At should return a string, got %v", obj)
	}
	if s != "cat" {
		t.Errorf("series.At(0) expected cat, got %v", s)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    cat
1    dog
2    eel
dtype: object`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding an int then float makes a float series
func TestConvertIntToFloat(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(123)
	builder.push(45.6)
	builder.push(789)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "float64" {
		t.Errorf("expected dtype == float64, got %q", series.dtype)
	}
	if series.which != typeFloat {
		t.Errorf("expected which == float, got %v", series.which)
	}
	expect := []interface{}{123.0, 45.6, 789.0}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	f, ok := obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float, got %v", obj)
	}
	if f != 123.0 {
		t.Errorf("series.At(0) expected 123.0, got %v", f)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    123.0
1     45.6
2    789.0
dtype: float64`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding a float then object makes an object series
func TestConvertFloatToObj(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(12.3)
	builder.push("ok")
	builder.push(78.9)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
	if series.which != typeObj {
		t.Errorf("expected which == object, got %v", series.which)
	}
	expect := []interface{}{12.3, "ok", 78.9}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	f, ok := obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float, got %v", obj)
	}
	if f != 12.3 {
		t.Errorf("series.At(0) expected 123, got %v", f)
	}
	obj = series.At(1)
	s, ok := obj.(string)
	if !ok {
		t.Fatalf("series.At should return a string, got %v", obj)
	}
	if s != "ok" {
		t.Errorf("series.At(1) expected ok, got %v", s)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    12.3
1      ok
2    78.9
dtype: object`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding an int then nil makes an object series
func TestIntWithNone(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(123)
	builder.push(nil)
	builder.push(789)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "float64" {
		t.Errorf("expected dtype == float64, got %q", series.dtype)
	}
	if series.which != typeFloat {
		t.Errorf("expected which == float, got %v", series.which)
	}
	expect := []interface{}{123.0, math.NaN(), 789.0}
	actual := series.values()
	if diff := cmp.Diff(expect, actual, cmpopts.EquateNaNs()); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	f, ok := obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float64, got %v of %T", obj, obj)
	}
	if f != 123.0 {
		t.Errorf("series.At(0) expected 123.0, got %v", f)
	}
	obj = series.At(1)
	f = obj.(float64)
	if !math.IsNaN(f) {
		t.Errorf("series.At(1) expected NaN, got %v", obj)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    123.0
1      NaN
2    789.0
dtype: float64`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding a float then null makes a float series
func TestFloatWithNone(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(12.3)
	builder.push(nil)
	builder.push(78.9)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "float64" {
		t.Errorf("expected dtype == float64, got %q", series.dtype)
	}
	if series.which != typeFloat {
		t.Errorf("expected which == float, got %v", series.which)
	}

	// Because series contains NaN, need to manually compare each element
	actual := series.values()
	expect := []interface{}{12.3, math.NaN(), 78.9}
	if diff := cmp.Diff(expect, actual, cmpopts.EquateNaNs()); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	f, ok := obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float, got %v", obj)
	}
	if f != 12.3 {
		t.Errorf("series.At(0) expected 12.3, got %v", f)
	}
	obj = series.At(1)
	f, ok = obj.(float64)
	if !ok {
		t.Fatalf("series.At should return a float, got %v", obj)
	}
	if !math.IsNaN(f) {
		t.Errorf("series.At(1) expected NaN, got %v", f)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0    12.3
1     NaN
2    78.9
dtype: float64`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding a bool then null makes an object series
func TestBoolWithNone(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(true)
	builder.push(false)
	builder.push(nil)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
	if series.which != typeObj {
		t.Errorf("expected which == object, got %v", series.which)
	}
	expect := []interface{}{true, false, nil}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	b, ok := obj.(bool)
	if !ok {
		t.Fatalf("series.At should return a bool, got %v", obj)
	}
	if !b {
		t.Errorf("series.At(0) expected true")
	}

	// Verify the string representation
	text := series.String()
	expectText := `0     True
1    False
2     None
dtype: object`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test adding an object then null makes an object series
func TestObjectWithNone(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push("cat")
	builder.push(nil)
	builder.push("eel")
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series has the right type and contents
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
	if series.which != typeObj {
		t.Errorf("expected which == object, got %v", series.which)
	}
	expect := []interface{}{"cat", nil, "eel"}
	actual := series.values()
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}

	// Verify that retrieving elements works
	obj := series.At(0)
	s, ok := obj.(string)
	if !ok {
		t.Fatalf("series.At should return a string, got %v", obj)
	}
	if s != "cat" {
		t.Errorf("series.At(0) expected cat, got %v", s)
	}

	// Verify the string representation
	text := series.String()
	expectText := `0     cat
1    None
2     eel
dtype: object`
	if diff := cmp.Diff(expectText, text); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

// Test a string following a bool is coerces to objects
func TestBuildBoolThenString(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(true)
	builder.push("")
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series is objects now
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
}

// Test a string following a bool is coerces to objects
func TestBuildBoolThenInt(t *testing.T) {
	builder := newTypedSliceBuilder(0)
	builder.push(true)
	builder.push(123)
	err := builder.error()
	if err != nil {
		t.Fatal(err)
	}
	series := builder.toSeries(nil, "")

	// Verify that the series is objects now
	if series.dtype != "object" {
		t.Errorf("expected dtype == object, got %q", series.dtype)
	}
}
