package chimera_test

import (
	"net/http"

	"github.com/matt1484/chimera"
)

type TestPrimitiveHeaderParams struct {
	SimpleStr  string  `param:"simpstr,in=header,style=simple"`
	SimpleU8   uint8   `param:"simpuint,in=header,style=simple"`
	SimpleU16  uint16  `param:"simpuint,in=header,style=simple"`
	SimpleU32  uint32  `param:"simpuint,in=header,style=simple"`
	SimpleU64  uint64  `param:"simpuint,in=header,style=simple"`
	SimpleUint uint    `param:"simpuint,in=header,style=simple"`
	SimpleI8   int8    `param:"simpint,in=header,style=simple"`
	SimpleI16  int16   `param:"simpint,in=header,style=simple"`
	SimpleI32  int32   `param:"simpint,in=header,style=simple"`
	SimpleI64  int64   `param:"simpint,in=header,style=simple"`
	SimpleInt  int     `param:"simpint,in=header,style=simple"`
	SimpleF32  float32 `param:"simpfloat,in=header,style=simple"`
	SimpleF64  float64 `param:"simpfloat,in=header,style=simple"`

	SimpleExplodeStr  string  `param:"simpexstr,in=header,explode,style=simple"`
	SimpleExplodeU8   uint8   `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU16  uint16  `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU32  uint32  `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU64  uint64  `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeUint uint    `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeI8   int8    `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI16  int16   `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI32  int32   `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI64  int64   `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeInt  int     `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeF32  float32 `param:"simpexfloat,in=header,explode,style=simple"`
	SimpleExplodeF64  float64 `param:"simpexfloat,in=header,explode,style=simple"`
}

type TestComplexHeaderParams struct {
	SimpleStr    []string         `param:"simpstr,in=header,style=simple"`
	SimpleU8     []uint8          `param:"simpuint,in=header,style=simple"`
	SimpleU16    []uint16         `param:"simpuint,in=header,style=simple"`
	SimpleU32    []uint32         `param:"simpuint,in=header,style=simple"`
	SimpleU64    []uint64         `param:"simpuint,in=header,style=simple"`
	SimpleUint   []uint           `param:"simpuint,in=header,style=simple"`
	SimpleI8     []int8           `param:"simpint,in=header,style=simple"`
	SimpleI16    []int16          `param:"simpint,in=header,style=simple"`
	SimpleI32    []int32          `param:"simpint,in=header,style=simple"`
	SimpleI64    []int64          `param:"simpint,in=header,style=simple"`
	SimpleInt    []int            `param:"simpint,in=header,style=simple"`
	SimpleF32    []float32        `param:"simpfloat,in=header,style=simple"`
	SimpleF64    []float64        `param:"simpfloat,in=header,style=simple"`
	SimpleStruct TestStructParams `param:"simpstruct,in=header,style=simple"`

	SimpleExplodeStr    []string         `param:"simpexstr,in=header,explode,style=simple"`
	SimpleExplodeU8     []uint8          `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU16    []uint16         `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU32    []uint32         `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeU64    []uint64         `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeUint   []uint           `param:"simpexuint,in=header,explode,style=simple"`
	SimpleExplodeI8     []int8           `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI16    []int16          `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI32    []int32          `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeI64    []int64          `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeInt    []int            `param:"simpexint,in=header,explode,style=simple"`
	SimpleExplodeF32    []float32        `param:"simpexfloat,in=header,explode,style=simple"`
	SimpleExplodeF64    []float64        `param:"simpexfloat,in=header,explode,style=simple"`
	SimpleExplodeStruct TestStructParams `param:"simpexstruct,in=header,explode,style=simple"`
}

func (c *TestValidCustomParam) UnmarshalHeaderParam(val string, t chimera.ParamStructTag) error {
	*c = "test"
	return nil
}

var (
	testValidSimplePrimitiveHeaderValues = http.Header{
		"simpstr":     []string{"ateststring"},
		"simpuint":    []string{"123"},
		"simpint":     []string{"-123"},
		"simpfloat":   []string{"123.45"},
		"simpexstr":   []string{"test..."},
		"simpexuint":  []string{"255"},
		"simpexint":   []string{"0"},
		"simpexfloat": []string{"-123.45"},
	}
	testValidSimpleComplexHeaderValues = http.Header{
		"simpstr":      []string{"string1,string2"},
		"simpuint":     []string{"0,123"},
		"simpint":      []string{"123,-123"},
		"simpfloat":    []string{"123.45,0.0"},
		"simpstruct":   []string{"stringprop,propstring,intprop,123"},
		"simpexstr":    []string{"test..."},
		"simpexuint":   []string{"255"},
		"simpexint":    []string{"0"},
		"simpexfloat":  []string{"-123.45"},
		"simpexstruct": []string{"stringprop=propstring,intprop=123"},
	}
)

var (
	testPrimitiveHeaderParams = TestPrimitiveHeaderParams{
		SimpleStr:  "ateststring",
		SimpleU8:   123,
		SimpleU16:  123,
		SimpleU32:  123,
		SimpleU64:  123,
		SimpleUint: 123,
		SimpleI8:   -123,
		SimpleI16:  -123,
		SimpleI32:  -123,
		SimpleI64:  -123,
		SimpleInt:  -123,
		SimpleF32:  123.45,
		SimpleF64:  123.45,

		SimpleExplodeStr:  "test...",
		SimpleExplodeU8:   255,
		SimpleExplodeU16:  255,
		SimpleExplodeU32:  255,
		SimpleExplodeU64:  255,
		SimpleExplodeUint: 255,
		SimpleExplodeI8:   0,
		SimpleExplodeI16:  0,
		SimpleExplodeI32:  0,
		SimpleExplodeI64:  0,
		SimpleExplodeInt:  0,
		SimpleExplodeF32:  -123.45,
		SimpleExplodeF64:  -123.45,
	}

	testComplexHeaderParams = TestComplexHeaderParams{
		SimpleStr:  []string{"string1", "string2"},
		SimpleU8:   []uint8{0, 123},
		SimpleU16:  []uint16{0, 123},
		SimpleU32:  []uint32{0, 123},
		SimpleU64:  []uint64{0, 123},
		SimpleUint: []uint{0, 123},
		SimpleI8:   []int8{123, -123},
		SimpleI16:  []int16{123, -123},
		SimpleI32:  []int32{123, -123},
		SimpleI64:  []int64{123, -123},
		SimpleInt:  []int{123, -123},
		SimpleF32:  []float32{123.45, 0.0},
		SimpleF64:  []float64{123.45, 0.0},
		SimpleStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		SimpleExplodeStr:  []string{"test..."},
		SimpleExplodeU8:   []uint8{255},
		SimpleExplodeU16:  []uint16{255},
		SimpleExplodeU32:  []uint32{255},
		SimpleExplodeU64:  []uint64{255},
		SimpleExplodeUint: []uint{255},
		SimpleExplodeI8:   []int8{0},
		SimpleExplodeI16:  []int16{0},
		SimpleExplodeI32:  []int32{0},
		SimpleExplodeI64:  []int64{0},
		SimpleExplodeInt:  []int{0},
		SimpleExplodeF32:  []float32{-123.45},
		SimpleExplodeF64:  []float64{-123.45},
		SimpleExplodeStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},
	}
)
