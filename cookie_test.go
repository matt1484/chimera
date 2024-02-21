package chimera_test

import "net/http"

type TestPrimitiveCookieParams struct {
	FormStr  string  `param:"formstr,in=cookie,style=form"`
	FormU8   uint8   `param:"formuint,in=cookie,style=form"`
	FormU16  uint16  `param:"formuint,in=cookie,style=form"`
	FormU32  uint32  `param:"formuint,in=cookie,style=form"`
	FormU64  uint64  `param:"formuint,in=cookie,style=form"`
	FormUint uint    `param:"formuint,in=cookie,style=form"`
	FormI8   int8    `param:"formint,in=cookie,style=form"`
	FormI16  int16   `param:"formint,in=cookie,style=form"`
	FormI32  int32   `param:"formint,in=cookie,style=form"`
	FormI64  int64   `param:"formint,in=cookie,style=form"`
	FormInt  int     `param:"formint,in=cookie,style=form"`
	FormF32  float32 `param:"formfloat,in=cookie,style=form"`
	FormF64  float64 `param:"formfloat,in=cookie,style=form"`

	FormExplodeStr  string  `param:"formexstr,in=cookie,explode,style=form"`
	FormExplodeU8   uint8   `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU16  uint16  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU32  uint32  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU64  uint64  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeUint uint    `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeI8   int8    `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI16  int16   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI32  int32   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI64  int64   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeInt  int     `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeF32  float32 `param:"formexfloat,in=cookie,explode,style=form"`
	FormExplodeF64  float64 `param:"formexfloat,in=cookie,explode,style=form"`
}

type TestComplexCookieParams struct {
	FormStr    []string         `param:"formstr,in=cookie,style=form"`
	FormU8     []uint8          `param:"formuint,in=cookie,style=form"`
	FormU16    []uint16         `param:"formuint,in=cookie,style=form"`
	FormU32    []uint32         `param:"formuint,in=cookie,style=form"`
	FormU64    []uint64         `param:"formuint,in=cookie,style=form"`
	FormUint   []uint           `param:"formuint,in=cookie,style=form"`
	FormI8     []int8           `param:"formint,in=cookie,style=form"`
	FormI16    []int16          `param:"formint,in=cookie,style=form"`
	FormI32    []int32          `param:"formint,in=cookie,style=form"`
	FormI64    []int64          `param:"formint,in=cookie,style=form"`
	FormInt    []int            `param:"formint,in=cookie,style=form"`
	FormF32    []float32        `param:"formfloat,in=cookie,style=form"`
	FormF64    []float64        `param:"formfloat,in=cookie,style=form"`
	FormStruct TestStructParams `param:"formstruct,in=cookie,style=form"`

	FormExplodeStr  []string  `param:"formexstr,in=cookie,explode,style=form"`
	FormExplodeU8   []uint8   `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU16  []uint16  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU32  []uint32  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeU64  []uint64  `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeUint []uint    `param:"formexuint,in=cookie,explode,style=form"`
	FormExplodeI8   []int8    `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI16  []int16   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI32  []int32   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeI64  []int64   `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeInt  []int     `param:"formexint,in=cookie,explode,style=form"`
	FormExplodeF32  []float32 `param:"formexfloat,in=cookie,explode,style=form"`
	FormExplodeF64  []float64 `param:"formexfloat,in=cookie,explode,style=form"`
}

var (
	testValidFormPrimitiveCookieValues = []http.Cookie{
		{Name: "formstr", Value: "ateststring"},
		{Name: "formuint", Value: "123"},
		{Name: "formint", Value: "-123"},
		{Name: "formfloat", Value: "123.45"},
		{Name: "formexstr", Value: "test..."},
		{Name: "formexuint", Value: "255"},
		{Name: "formexint", Value: "0"},
		{Name: "formexfloat", Value: "-123.45"},
	}
	testValidFormComplexCookieValues = []http.Cookie{
		{Name: "formstr", Value: "string1,string2"},
		{Name: "formuint", Value: "0,123"},
		{Name: "formint", Value: "123,-123"},
		{Name: "formfloat", Value: "123.45,0.0"},
		{Name: "formstruct", Value: "stringprop,propstring,intprop,123"},
		{Name: "formexstr", Value: "test..."},
		{Name: "formexuint", Value: "255"},
		{Name: "formexint", Value: "0"},
		{Name: "formexfloat", Value: "-123.45"},
	}

	testPrimitiveCookieParams = TestPrimitiveCookieParams{
		FormStr:  "ateststring",
		FormU8:   123,
		FormU16:  123,
		FormU32:  123,
		FormU64:  123,
		FormUint: 123,
		FormI8:   -123,
		FormI16:  -123,
		FormI32:  -123,
		FormI64:  -123,
		FormInt:  -123,
		FormF32:  123.45,
		FormF64:  123.45,

		FormExplodeStr:  "test...",
		FormExplodeU8:   255,
		FormExplodeU16:  255,
		FormExplodeU32:  255,
		FormExplodeU64:  255,
		FormExplodeUint: 255,
		FormExplodeI8:   0,
		FormExplodeI16:  0,
		FormExplodeI32:  0,
		FormExplodeI64:  0,
		FormExplodeInt:  0,
		FormExplodeF32:  -123.45,
		FormExplodeF64:  -123.45,
	}

	testComplexCookieParams = TestComplexCookieParams{
		FormStr:  []string{"string1", "string2"},
		FormU8:   []uint8{0, 123},
		FormU16:  []uint16{0, 123},
		FormU32:  []uint32{0, 123},
		FormU64:  []uint64{0, 123},
		FormUint: []uint{0, 123},
		FormI8:   []int8{123, -123},
		FormI16:  []int16{123, -123},
		FormI32:  []int32{123, -123},
		FormI64:  []int64{123, -123},
		FormInt:  []int{123, -123},
		FormF32:  []float32{123.45, 0.0},
		FormF64:  []float64{123.45, 0.0},
		FormStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		FormExplodeStr:  []string{"test..."},
		FormExplodeU8:   []uint8{255},
		FormExplodeU16:  []uint16{255},
		FormExplodeU32:  []uint32{255},
		FormExplodeU64:  []uint64{255},
		FormExplodeUint: []uint{255},
		FormExplodeI8:   []int8{0},
		FormExplodeI16:  []int16{0},
		FormExplodeI32:  []int32{0},
		FormExplodeI64:  []int64{0},
		FormExplodeInt:  []int{0},
		FormExplodeF32:  []float32{-123.45},
		FormExplodeF64:  []float64{-123.45},
	}
)
