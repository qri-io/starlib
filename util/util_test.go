package util

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestIsEmptyString(t *testing.T) {
	if !IsEmptyString(starlark.String("")) {
		t.Error("empty string should equal true")
	}

	if IsEmptyString(".") {
		t.Error("non-empty string shouldn't be empty")
	}
}

func TestAsString(t *testing.T) {
	cases := []struct {
		in       starlark.Value
		got, err string
	}{
		{starlark.String("foo"), "foo", ""},
		{starlark.String("\"foo'"), "\"foo'", ""},
		{starlark.Bool(true), "", "invalid syntax"},
	}

	for i, c := range cases {
		got, err := asString(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if c.got != got {
			t.Errorf("case %d. expected: '%s', got: '%s'", i, c.got, got)
		}
	}
}

func TestMarshal(t *testing.T) {
	expectedStringDict := starlark.NewDict(1)
	expectedStringDict.SetKey(starlark.String("foo"), starlark.MakeInt(42))

	expectedIntDict := starlark.NewDict(1)
	expectedIntDict.SetKey(starlark.MakeInt(42*2), starlark.MakeInt(42))

	ct, _ := (&customType{42}).MarshalStarlark()
	expectedStrDictCustomType := starlark.NewDict(2)
	expectedStrDictCustomType.SetKey(starlark.String("foo"), starlark.MakeInt(42))
	expectedStrDictCustomType.SetKey(starlark.String("bar"), ct)

	cases := []struct {
		in   interface{}
		want starlark.Value
		err  string
	}{
		{nil, starlark.None, ""},
		{true, starlark.True, ""},
		{"foo", starlark.String("foo"), ""},
		{42, starlark.MakeInt(42), ""},
		{int8(42), starlark.MakeInt(42), ""},
		{int16(42), starlark.MakeInt(42), ""},
		{int32(42), starlark.MakeInt(42), ""},
		{int64(42), starlark.MakeInt(42), ""},
		{uint(42), starlark.MakeUint(42), ""},
		{uint8(42), starlark.MakeUint(42), ""},
		{uint16(42), starlark.MakeUint(42), ""},
		{uint32(42), starlark.MakeUint(42), ""},
		{uint64(42), starlark.MakeUint64(42), ""},
		{float32(42), starlark.Float(42), ""},
		{42., starlark.Float(42), ""},
		{[]interface{}{42}, starlark.NewList([]starlark.Value{starlark.MakeInt(42)}), ""},
		{map[string]interface{}{"foo": 42}, expectedStringDict, ""},
		{map[interface{}]interface{}{"foo": 42}, expectedStringDict, ""},
		{map[interface{}]interface{}{42 * 2: 42}, expectedIntDict, ""},
		{&customType{42}, ct, ""},
		{map[string]interface{}{"foo": 42, "bar": &customType{42}}, expectedStrDictCustomType, ""},
		{map[interface{}]interface{}{"foo": 42, "bar": &customType{42}}, expectedStrDictCustomType, ""},
		{[]interface{}{42, &customType{42}}, starlark.NewList([]starlark.Value{starlark.MakeInt(42), ct}), ""},
		{&invalidCustomType{42}, starlark.None, "unrecognized type: &util.invalidCustomType{Foo:42}"},
	}

	for i, c := range cases {
		got, err := Marshal(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: %q, got: %q (%T -> %T)", i, c.err, err, c.in, c.want)
			continue
		}

		assert.EqualValues(t, c.want, got, "case %d: %T -> %T", i, c.in, c.want)
	}
}

func TestUnmarshal(t *testing.T) {
	strDict := starlark.NewDict(1)
	strDict.SetKey(starlark.String("foo"), starlark.MakeInt(42))

	intDict := starlark.NewDict(1)
	intDict.SetKey(starlark.MakeInt(42*2), starlark.MakeInt(42))

	ct, _ := (&customType{42}).MarshalStarlark()
	strDictCT := starlark.NewDict(2)
	strDictCT.SetKey(starlark.String("foo"), starlark.MakeInt(42))
	strDictCT.SetKey(starlark.String("bar"), ct)

	cases := []struct {
		in   starlark.Value
		want interface{}
		err  string
	}{
		{starlark.None, nil, ""},
		{starlark.True, true, ""},
		{starlark.String("foo"), "foo", ""},
		{starlark.MakeInt(42), 42, ""},
		{starlark.MakeInt(42), int8(42), ""},
		{starlark.MakeInt(42), int16(42), ""},
		{starlark.MakeInt(42), int32(42), ""},
		{starlark.MakeInt(42), int64(42), ""},
		{starlark.MakeUint(42), uint(42), ""},
		{starlark.MakeUint(42), uint8(42), ""},
		{starlark.MakeUint(42), uint16(42), ""},
		{starlark.MakeUint(42), uint32(42), ""},
		{starlark.MakeUint64(42), uint64(42), ""},
		{starlark.Float(42), float32(42), ""},
		{starlark.Float(42), 42., ""},
		{starlark.NewList([]starlark.Value{starlark.MakeInt(42)}), []interface{}{42}, ""},
		{strDict, map[string]interface{}{"foo": 42}, ""},
		{intDict, map[interface{}]interface{}{42 * 2: 42}, ""},
		{ct, &customType{42}, ""},
		{strDictCT, map[string]interface{}{"foo": 42, "bar": &customType{42}}, ""},
		{starlark.NewList([]starlark.Value{starlark.MakeInt(42), ct}), []interface{}{42, &customType{42}}, ""},
	}

	for i, c := range cases {
		got, err := Unmarshal(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: %q, got: %q %T -> %T", i, c.err, err, c.in, c.want)
			continue
		}

		assert.EqualValues(t, c.want, got, "case %d: %T -> %T", i, c.in, c.want)
	}
}

func TestLifeCycle(t *testing.T) {
	t.Run("once", func(t *testing.T) {
		// golang value
		goVal := &customType{42}
		// starlark value
		slVal, err := Marshal(goVal)
		assert.NoError(t, err)

		assert.IsType(t, &starlarkstruct.Struct{}, slVal)

		gotGoVal, err := Unmarshal(slVal)
		assert.NoError(t, err)
		log.Println(slVal.String())
		assert.EqualValues(t, goVal, gotGoVal)
	})

	t.Run("asDictValue", func(t *testing.T) {
		// golang value
		goVal := map[string]interface{}{
			"foo": &customType{42},
		}

		// starlark value
		slVal, err := Marshal(goVal)
		assert.NoError(t, err)

		wantSlVal := starlark.NewDict(1)
		assert.IsType(t, wantSlVal, slVal)

		wantSlVal.SetKey(starlark.String("foo"), func() starlark.Value { v, _ := Marshal(&customType{42}); return v }())
		assert.EqualValues(t, wantSlVal, slVal)

		gotGoVal, err := Unmarshal(slVal)
		assert.NoError(t, err)
		log.Println(slVal.String())
		assert.EqualValues(t, goVal, gotGoVal)
	})

	t.Run("asListValue", func(t *testing.T) {
		// golang value
		goVal := []interface{}{
			&customType{42},
			&customType{42},
		}

		// starlark value
		slVal, err := Marshal(goVal)
		assert.NoError(t, err)

		wantSlVal := starlark.NewList(nil)
		wantSlVal.Append(func() starlark.Value { v, _ := Marshal(&customType{42}); return v }())
		wantSlVal.Append(func() starlark.Value { v, _ := Marshal(&customType{42}); return v }())
		assert.IsType(t, wantSlVal, slVal)

		assert.EqualValues(t, wantSlVal, slVal)

		gotGoVal, err := Unmarshal(slVal)
		assert.NoError(t, err)
		log.Println(slVal.String())
		assert.EqualValues(t, goVal, gotGoVal)
	})

}

type invalidCustomType struct {
	Foo int64
}

type customType invalidCustomType

func (t *customType) UnmarshalStarlark(v starlark.Value) error {
	// asserts
	if v.Type() != "struct" {
		return fmt.Errorf("not expected top level type, want struct, got %q", v.Type())
	}
	if _, ok := v.(*starlarkstruct.Struct).Constructor().(*customType); !ok {
		return fmt.Errorf("not expected construct type got %T, want %T", v.(*starlarkstruct.Struct).Constructor(), t)
	}

	// TODO: refactoring transform data

	mustInt64 := func(sv starlark.Value) int64 {
		i, _ := sv.(starlark.Int).Int64()
		return i
	}

	data := starlark.StringDict{}
	v.(*starlarkstruct.Struct).ToStringDict(data)

	*t = customType{
		Foo: mustInt64(data["foo"]),
	}
	return nil
}

func (t *customType) MarshalStarlark() (starlark.Value, error) {
	v := starlarkstruct.FromStringDict(&customType{}, starlark.StringDict{
		"foo": starlark.MakeInt64(t.Foo),
	})
	return v, nil
}

func (c customType) String() string {
	return "customType"
}

func (c customType) Type() string { return "test.customType" }

func (customType) Freeze() {}

func (c customType) Truth() starlark.Bool {
	return starlark.True
}

func (c customType) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

var _ Unmarshaler = (*customType)(nil)
var _ Marshaler = (*customType)(nil)
var _ starlark.Value = (*customType)(nil)
