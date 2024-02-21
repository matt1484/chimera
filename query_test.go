package chimera_test

import (
	"net/url"
)

type TestPrimitiveQueryParams struct {
	FormStr  string  `param:"formstr,in=query,style=form"`
	FormU8   uint8   `param:"formuint,in=query,style=form"`
	FormU16  uint16  `param:"formuint,in=query,style=form"`
	FormU32  uint32  `param:"formuint,in=query,style=form"`
	FormU64  uint64  `param:"formuint,in=query,style=form"`
	FormUint uint    `param:"formuint,in=query,style=form"`
	FormI8   int8    `param:"formint,in=query,style=form"`
	FormI16  int16   `param:"formint,in=query,style=form"`
	FormI32  int32   `param:"formint,in=query,style=form"`
	FormI64  int64   `param:"formint,in=query,style=form"`
	FormInt  int     `param:"formint,in=query,style=form"`
	FormF32  float32 `param:"formfloat,in=query,style=form"`
	FormF64  float64 `param:"formfloat,in=query,style=form"`

	FormExplodeStr  string  `param:"formexstr,in=query,explode,style=form"`
	FormExplodeU8   uint8   `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU16  uint16  `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU32  uint32  `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU64  uint64  `param:"formexuint,in=query,explode,style=form"`
	FormExplodeUint uint    `param:"formexuint,in=query,explode,style=form"`
	FormExplodeI8   int8    `param:"formexint,in=query,explode,style=form"`
	FormExplodeI16  int16   `param:"formexint,in=query,explode,style=form"`
	FormExplodeI32  int32   `param:"formexint,in=query,explode,style=form"`
	FormExplodeI64  int64   `param:"formexint,in=query,explode,style=form"`
	FormExplodeInt  int     `param:"formexint,in=query,explode,style=form"`
	FormExplodeF32  float32 `param:"formexfloat,in=query,explode,style=form"`
	FormExplodeF64  float64 `param:"formexfloat,in=query,explode,style=form"`
}

type TestComplexQueryParams struct {
	FormStr    []string         `param:"formstr,in=query,style=form"`
	FormU8     []uint8          `param:"formuint,in=query,style=form"`
	FormU16    []uint16         `param:"formuint,in=query,style=form"`
	FormU32    []uint32         `param:"formuint,in=query,style=form"`
	FormU64    []uint64         `param:"formuint,in=query,style=form"`
	FormUint   []uint           `param:"formuint,in=query,style=form"`
	FormI8     []int8           `param:"formint,in=query,style=form"`
	FormI16    []int16          `param:"formint,in=query,style=form"`
	FormI32    []int32          `param:"formint,in=query,style=form"`
	FormI64    []int64          `param:"formint,in=query,style=form"`
	FormInt    []int            `param:"formint,in=query,style=form"`
	FormF32    []float32        `param:"formfloat,in=query,style=form"`
	FormF64    []float64        `param:"formfloat,in=query,style=form"`
	FormStruct TestStructParams `param:"formstruct,in=query,style=form"`

	FormExplodeStr    []string         `param:"formexstr,in=query,explode,style=form"`
	FormExplodeU8     []uint8          `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU16    []uint16         `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU32    []uint32         `param:"formexuint,in=query,explode,style=form"`
	FormExplodeU64    []uint64         `param:"formexuint,in=query,explode,style=form"`
	FormExplodeUint   []uint           `param:"formexuint,in=query,explode,style=form"`
	FormExplodeI8     []int8           `param:"formexint,in=query,explode,style=form"`
	FormExplodeI16    []int16          `param:"formexint,in=query,explode,style=form"`
	FormExplodeI32    []int32          `param:"formexint,in=query,explode,style=form"`
	FormExplodeI64    []int64          `param:"formexint,in=query,explode,style=form"`
	FormExplodeInt    []int            `param:"formexint,in=query,explode,style=form"`
	FormExplodeF32    []float32        `param:"formexfloat,in=query,explode,style=form"`
	FormExplodeF64    []float64        `param:"formexfloat,in=query,explode,style=form"`
	FormExplodeStruct TestStructParams `param:"formexstruct,in=query,explode,style=form"`

	SpaceStr  []string  `param:"spacestr,in=query,style=spaceDelimited"`
	SpaceU8   []uint8   `param:"spaceuint,in=query,style=spaceDelimited"`
	SpaceU16  []uint16  `param:"spaceuint,in=query,style=spaceDelimited"`
	SpaceU32  []uint32  `param:"spaceuint,in=query,style=spaceDelimited"`
	SpaceU64  []uint64  `param:"spaceuint,in=query,style=spaceDelimited"`
	SpaceUint []uint    `param:"spaceuint,in=query,style=spaceDelimited"`
	SpaceI8   []int8    `param:"spaceint,in=query,style=spaceDelimited"`
	SpaceI16  []int16   `param:"spaceint,in=query,style=spaceDelimited"`
	SpaceI32  []int32   `param:"spaceint,in=query,style=spaceDelimited"`
	SpaceI64  []int64   `param:"spaceint,in=query,style=spaceDelimited"`
	SpaceInt  []int     `param:"spaceint,in=query,style=spaceDelimited"`
	SpaceF32  []float32 `param:"spacefloat,in=query,style=spaceDelimited"`
	SpaceF64  []float64 `param:"spacefloat,in=query,style=spaceDelimited"`

	SpaceExplodeStr  []string  `param:"spaceexstr,in=query,explode,style=spaceDelimited"`
	SpaceExplodeU8   []uint8   `param:"spaceexuint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeU16  []uint16  `param:"spaceexuint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeU32  []uint32  `param:"spaceexuint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeU64  []uint64  `param:"spaceexuint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeUint []uint    `param:"spaceexuint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeI8   []int8    `param:"spaceexint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeI16  []int16   `param:"spaceexint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeI32  []int32   `param:"spaceexint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeI64  []int64   `param:"spaceexint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeInt  []int     `param:"spaceexint,in=query,explode,style=spaceDelimited"`
	SpaceExplodeF32  []float32 `param:"spaceexfloat,in=query,explode,style=spaceDelimited"`
	SpaceExplodeF64  []float64 `param:"spaceexfloat,in=query,explode,style=spaceDelimited"`

	PipeStr  []string  `param:"pipestr,in=query,style=pipeDelimited"`
	PipeU8   []uint8   `param:"pipeuint,in=query,style=pipeDelimited"`
	PipeU16  []uint16  `param:"pipeuint,in=query,style=pipeDelimited"`
	PipeU32  []uint32  `param:"pipeuint,in=query,style=pipeDelimited"`
	PipeU64  []uint64  `param:"pipeuint,in=query,style=pipeDelimited"`
	PipeUint []uint    `param:"pipeuint,in=query,style=pipeDelimited"`
	PipeI8   []int8    `param:"pipeint,in=query,style=pipeDelimited"`
	PipeI16  []int16   `param:"pipeint,in=query,style=pipeDelimited"`
	PipeI32  []int32   `param:"pipeint,in=query,style=pipeDelimited"`
	PipeI64  []int64   `param:"pipeint,in=query,style=pipeDelimited"`
	PipeInt  []int     `param:"pipeint,in=query,style=pipeDelimited"`
	PipeF32  []float32 `param:"pipefloat,in=query,style=pipeDelimited"`
	PipeF64  []float64 `param:"pipefloat,in=query,style=pipeDelimited"`

	PipeExplodeStr  []string  `param:"pipeexstr,in=query,explode,style=pipeDelimited"`
	PipeExplodeU8   []uint8   `param:"pipeexuint,in=query,explode,style=pipeDelimited"`
	PipeExplodeU16  []uint16  `param:"pipeexuint,in=query,explode,style=pipeDelimited"`
	PipeExplodeU32  []uint32  `param:"pipeexuint,in=query,explode,style=pipeDelimited"`
	PipeExplodeU64  []uint64  `param:"pipeexuint,in=query,explode,style=pipeDelimited"`
	PipeExplodeUint []uint    `param:"pipeexuint,in=query,explode,style=pipeDelimited"`
	PipeExplodeI8   []int8    `param:"pipeexint,in=query,explode,style=pipeDelimited"`
	PipeExplodeI16  []int16   `param:"pipeexint,in=query,explode,style=pipeDelimited"`
	PipeExplodeI32  []int32   `param:"pipeexint,in=query,explode,style=pipeDelimited"`
	PipeExplodeI64  []int64   `param:"pipeexint,in=query,explode,style=pipeDelimited"`
	PipeExplodeInt  []int     `param:"pipeexint,in=query,explode,style=pipeDelimited"`
	PipeExplodeF32  []float32 `param:"pipeexfloat,in=query,explode,style=pipeDelimited"`
	PipeExplodeF64  []float64 `param:"pipeexfloat,in=query,explode,style=pipeDelimited"`

	DeepStruct TestStructParams `param:"deepstruct,in=query,style=deepObject"`
}

var (
	testValidFormPrimitiveQueryValues = url.Values{
		"formstr":     []string{"ateststring"},
		"formuint":    []string{"123"},
		"formint":     []string{"-123"},
		"formfloat":   []string{"123.45"},
		"formexstr":   []string{"test..."},
		"formexuint":  []string{"255"},
		"formexint":   []string{"0"},
		"formexfloat": []string{"-123.45"},
	}
	testValidFormComplexQueryValues = url.Values{
		"formstr":     []string{"string1,string2"},
		"formuint":    []string{"0,123"},
		"formint":     []string{"123,-123"},
		"formfloat":   []string{"123.45,0.0"},
		"formstruct":  []string{"stringprop,propstring,intprop,123"},
		"formexstr":   []string{"string3", "string4"},
		"formexuint":  []string{"255", "1"},
		"formexint":   []string{"0", "-1"},
		"formexfloat": []string{"-123.45", "1"},
		"stringprop":  []string{"propstring"},
		"intprop":     []string{"123"},

		"spacestr":     []string{"string1 string2"},
		"spaceuint":    []string{"0 123"},
		"spaceint":     []string{"123 -123"},
		"spacefloat":   []string{"123.45 0.0"},
		"spaceexstr":   []string{"string3", "string4"},
		"spaceexuint":  []string{"255", "1"},
		"spaceexint":   []string{"0", "-1"},
		"spaceexfloat": []string{"-123.45", "1"},

		"pipestr":     []string{"string1|string2"},
		"pipeuint":    []string{"0|123"},
		"pipeint":     []string{"123|-123"},
		"pipefloat":   []string{"123.45|0.0"},
		"pipeexstr":   []string{"string3", "string4"},
		"pipeexuint":  []string{"255", "1"},
		"pipeexint":   []string{"0", "-1"},
		"pipeexfloat": []string{"-123.45", "1"},

		"deepstruct[stringprop]": []string{"propstring"},
		"deepstruct[intprop]":    []string{"123"},
	}

	testPrimitiveQueryParams = TestPrimitiveQueryParams{
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

	testComplexQueryParams = TestComplexQueryParams{
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

		FormExplodeStr:  []string{"string3", "string4"},
		FormExplodeU8:   []uint8{255, 1},
		FormExplodeU16:  []uint16{255, 1},
		FormExplodeU32:  []uint32{255, 1},
		FormExplodeU64:  []uint64{255, 1},
		FormExplodeUint: []uint{255, 1},
		FormExplodeI8:   []int8{0, -1},
		FormExplodeI16:  []int16{0, -1},
		FormExplodeI32:  []int32{0, -1},
		FormExplodeI64:  []int64{0, -1},
		FormExplodeInt:  []int{0, -1},
		FormExplodeF32:  []float32{-123.45, 1},
		FormExplodeF64:  []float64{-123.45, 1},
		FormExplodeStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		SpaceStr:  []string{"string1", "string2"},
		SpaceU8:   []uint8{0, 123},
		SpaceU16:  []uint16{0, 123},
		SpaceU32:  []uint32{0, 123},
		SpaceU64:  []uint64{0, 123},
		SpaceUint: []uint{0, 123},
		SpaceI8:   []int8{123, -123},
		SpaceI16:  []int16{123, -123},
		SpaceI32:  []int32{123, -123},
		SpaceI64:  []int64{123, -123},
		SpaceInt:  []int{123, -123},
		SpaceF32:  []float32{123.45, 0.0},
		SpaceF64:  []float64{123.45, 0.0},

		SpaceExplodeStr:  []string{"string3", "string4"},
		SpaceExplodeU8:   []uint8{255, 1},
		SpaceExplodeU16:  []uint16{255, 1},
		SpaceExplodeU32:  []uint32{255, 1},
		SpaceExplodeU64:  []uint64{255, 1},
		SpaceExplodeUint: []uint{255, 1},
		SpaceExplodeI8:   []int8{0, -1},
		SpaceExplodeI16:  []int16{0, -1},
		SpaceExplodeI32:  []int32{0, -1},
		SpaceExplodeI64:  []int64{0, -1},
		SpaceExplodeInt:  []int{0, -1},
		SpaceExplodeF32:  []float32{-123.45, 1},
		SpaceExplodeF64:  []float64{-123.45, 1},

		PipeStr:  []string{"string1", "string2"},
		PipeU8:   []uint8{0, 123},
		PipeU16:  []uint16{0, 123},
		PipeU32:  []uint32{0, 123},
		PipeU64:  []uint64{0, 123},
		PipeUint: []uint{0, 123},
		PipeI8:   []int8{123, -123},
		PipeI16:  []int16{123, -123},
		PipeI32:  []int32{123, -123},
		PipeI64:  []int64{123, -123},
		PipeInt:  []int{123, -123},
		PipeF32:  []float32{123.45, 0.0},
		PipeF64:  []float64{123.45, 0.0},

		PipeExplodeStr:  []string{"string3", "string4"},
		PipeExplodeU8:   []uint8{255, 1},
		PipeExplodeU16:  []uint16{255, 1},
		PipeExplodeU32:  []uint32{255, 1},
		PipeExplodeU64:  []uint64{255, 1},
		PipeExplodeUint: []uint{255, 1},
		PipeExplodeI8:   []int8{0, -1},
		PipeExplodeI16:  []int16{0, -1},
		PipeExplodeI32:  []int32{0, -1},
		PipeExplodeI64:  []int64{0, -1},
		PipeExplodeInt:  []int{0, -1},
		PipeExplodeF32:  []float32{-123.45, 1},
		PipeExplodeF64:  []float64{-123.45, 1},

		DeepStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},
	}
)
