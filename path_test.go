package chimera_test

type TestPrimitivePathParams struct {
	SimpleStr  string  `param:"simpstr,in=path,style=simple"`
	SimpleU8   uint8   `param:"simpuint,in=path,style=simple"`
	SimpleU16  uint16  `param:"simpuint,in=path,style=simple"`
	SimpleU32  uint32  `param:"simpuint,in=path,style=simple"`
	SimpleU64  uint64  `param:"simpuint,in=path,style=simple"`
	SimpleUint uint    `param:"simpuint,in=path,style=simple"`
	SimpleI8   int8    `param:"simpint,in=path,style=simple"`
	SimpleI16  int16   `param:"simpint,in=path,style=simple"`
	SimpleI32  int32   `param:"simpint,in=path,style=simple"`
	SimpleI64  int64   `param:"simpint,in=path,style=simple"`
	SimpleInt  int     `param:"simpint,in=path,style=simple"`
	SimpleF32  float32 `param:"simpfloat,in=path,style=simple"`
	SimpleF64  float64 `param:"simpfloat,in=path,style=simple"`

	SimpleExplodeStr  string  `param:"simpexstr,in=path,explode,style=simple"`
	SimpleExplodeU8   uint8   `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU16  uint16  `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU32  uint32  `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU64  uint64  `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeUint uint    `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeI8   int8    `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI16  int16   `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI32  int32   `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI64  int64   `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeInt  int     `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeF32  float32 `param:"simpexfloat,in=path,explode,style=simple"`
	SimpleExplodeF64  float64 `param:"simpexfloat,in=path,explode,style=simple"`

	LabelStr  string  `param:"labelstr,in=path,style=label"`
	LabelU8   uint8   `param:"labeluint,in=path,style=label"`
	LabelU16  uint16  `param:"labeluint,in=path,style=label"`
	LabelU32  uint32  `param:"labeluint,in=path,style=label"`
	LabelU64  uint64  `param:"labeluint,in=path,style=label"`
	LabelUint uint    `param:"labeluint,in=path,style=label"`
	LabelI8   int8    `param:"labelint,in=path,style=label"`
	LabelI16  int16   `param:"labelint,in=path,style=label"`
	LabelI32  int32   `param:"labelint,in=path,style=label"`
	LabelI64  int64   `param:"labelint,in=path,style=label"`
	LabelInt  int     `param:"labelint,in=path,style=label"`
	LabelF32  float32 `param:"labelfloat,in=path,style=label"`
	LabelF64  float64 `param:"labelfloat,in=path,style=label"`

	LabelExplodeStr  string  `param:"labelexstr,in=path,explode,style=label"`
	LabelExplodeU8   uint8   `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU16  uint16  `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU32  uint32  `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU64  uint64  `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeUint uint    `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeI8   int8    `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI16  int16   `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI32  int32   `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI64  int64   `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeInt  int     `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeF32  float32 `param:"labelexfloat,in=path,explod,style=label"`
	LabelExplodeF64  float64 `param:"labelexfloat,in=path,explode,style=label"`

	MatrixStr  string  `param:"matrixstr,in=path,style=matrix"`
	MatrixU8   uint8   `param:"matrixuint,in=path,style=matrix"`
	MatrixU16  uint16  `param:"matrixuint,in=path,style=matrix"`
	MatrixU32  uint32  `param:"matrixuint,in=path,style=matrix"`
	MatrixU64  uint64  `param:"matrixuint,in=path,style=matrix"`
	MatrixUint uint    `param:"matrixuint,in=path,style=matrix"`
	MatrixI8   int8    `param:"matrixint,in=path,style=matrix"`
	MatrixI16  int16   `param:"matrixint,in=path,style=matrix"`
	MatrixI32  int32   `param:"matrixint,in=path,style=matrix"`
	MatrixI64  int64   `param:"matrixint,in=path,style=matrix"`
	MatrixInt  int     `param:"matrixint,in=path,style=matrix"`
	MatrixF32  float32 `param:"matrixfloat,in=path,style=matrix"`
	MatrixF64  float64 `param:"matrixfloat,in=path,style=matrix"`

	MatrixExplodeStr  string  `param:"matrixexstr,in=path,explode,style=matrix"`
	MatrixExplodeU8   uint8   `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU16  uint16  `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU32  uint32  `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU64  uint64  `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeUint uint    `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeI8   int8    `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI16  int16   `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI32  int32   `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI64  int64   `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeInt  int     `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeF32  float32 `param:"matrixexfloat,in=path,explode,style=matrix"`
	MatrixExplodeF64  float64 `param:"matrixexfloat,in=path,explode,style=matrix"`
}

type TestComplexPathParams struct {
	SimpleStr    []string         `param:"simpstr,in=path,style=simple"`
	SimpleU8     []uint8          `param:"simpuint,in=path,style=simple"`
	SimpleU16    []uint16         `param:"simpuint,in=path,style=simple"`
	SimpleU32    []uint32         `param:"simpuint,in=path,style=simple"`
	SimpleU64    []uint64         `param:"simpuint,in=path,style=simple"`
	SimpleUint   []uint           `param:"simpuint,in=path,style=simple"`
	SimpleI8     []int8           `param:"simpint,in=path,style=simple"`
	SimpleI16    []int16          `param:"simpint,in=path,style=simple"`
	SimpleI32    []int32          `param:"simpint,in=path,style=simple"`
	SimpleI64    []int64          `param:"simpint,in=path,style=simple"`
	SimpleInt    []int            `param:"simpint,in=path,style=simple"`
	SimpleF32    []float32        `param:"simpfloat,in=path,style=simple"`
	SimpleF64    []float64        `param:"simpfloat,in=path,style=simple"`
	SimpleStruct TestStructParams `param:"simpstruct,in=path,style=simple"`

	SimpleExplodeStr    []string         `param:"simpexstr,in=path,explode,style=simple"`
	SimpleExplodeU8     []uint8          `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU16    []uint16         `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU32    []uint32         `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeU64    []uint64         `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeUint   []uint           `param:"simpexuint,in=path,explode,style=simple"`
	SimpleExplodeI8     []int8           `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI16    []int16          `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI32    []int32          `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeI64    []int64          `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeInt    []int            `param:"simpexint,in=path,explode,style=simple"`
	SimpleExplodeF32    []float32        `param:"simpexfloat,in=path,explode,style=simple"`
	SimpleExplodeF64    []float64        `param:"simpexfloat,in=path,explode,style=simple"`
	SimpleExplodeStruct TestStructParams `param:"simpexstruct,in=path,explode,style=simple"`

	LabelStr    []string         `param:"labelstr,in=path,style=label"`
	LabelU8     []uint8          `param:"labeluint,in=path,style=label"`
	LabelU16    []uint16         `param:"labeluint,in=path,style=label"`
	LabelU32    []uint32         `param:"labeluint,in=path,style=label"`
	LabelU64    []uint64         `param:"labeluint,in=path,style=label"`
	LabelUint   []uint           `param:"labeluint,in=path,style=label"`
	LabelI8     []int8           `param:"labelint,in=path,style=label"`
	LabelI16    []int16          `param:"labelint,in=path,style=label"`
	LabelI32    []int32          `param:"labelint,in=path,style=label"`
	LabelI64    []int64          `param:"labelint,in=path,style=label"`
	LabelInt    []int            `param:"labelint,in=path,style=label"`
	LabelF32    []float32        `param:"labelfloat,in=path,style=label"`
	LabelF64    []float64        `param:"labelfloat,in=path,style=label"`
	LabelStruct TestStructParams `param:"labelstruct,in=path,style=label"`

	LabelExplodeStr    []string         `param:"labelexstr,in=path,explode,style=label"`
	LabelExplodeU8     []uint8          `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU16    []uint16         `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU32    []uint32         `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeU64    []uint64         `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeUint   []uint           `param:"labelexuint,in=path,explode,style=label"`
	LabelExplodeI8     []int8           `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI16    []int16          `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI32    []int32          `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeI64    []int64          `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeInt    []int            `param:"labelexint,in=path,explode,style=label"`
	LabelExplodeStruct TestStructParams `param:"labelexstruct,in=path,explode,style=label"`
	// this creates ambiguity so Im ignoring this for now
	// LabelExplodeF32  []float32 `param:"labelexfloat,in=path,explod,style=label"`
	// LabelExplodeF64  []float64 `param:"labelexfloat,in=path,explode,style=label"`

	MatrixStr    []string         `param:"matrixstr,in=path,style=matrix"`
	MatrixU8     []uint8          `param:"matrixuint,in=path,style=matrix"`
	MatrixU16    []uint16         `param:"matrixuint,in=path,style=matrix"`
	MatrixU32    []uint32         `param:"matrixuint,in=path,style=matrix"`
	MatrixU64    []uint64         `param:"matrixuint,in=path,style=matrix"`
	MatrixUint   []uint           `param:"matrixuint,in=path,style=matrix"`
	MatrixI8     []int8           `param:"matrixint,in=path,style=matrix"`
	MatrixI16    []int16          `param:"matrixint,in=path,style=matrix"`
	MatrixI32    []int32          `param:"matrixint,in=path,style=matrix"`
	MatrixI64    []int64          `param:"matrixint,in=path,style=matrix"`
	MatrixInt    []int            `param:"matrixint,in=path,style=matrix"`
	MatrixF32    []float32        `param:"matrixfloat,in=path,style=matrix"`
	MatrixF64    []float64        `param:"matrixfloat,in=path,style=matrix"`
	MatrixStruct TestStructParams `param:"matrixstruct,in=path,style=matrix"`

	MatrixExplodeStr    []string         `param:"matrixexstr,in=path,explode,style=matrix"`
	MatrixExplodeU8     []uint8          `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU16    []uint16         `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU32    []uint32         `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeU64    []uint64         `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeUint   []uint           `param:"matrixexuint,in=path,explode,style=matrix"`
	MatrixExplodeI8     []int8           `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI16    []int16          `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI32    []int32          `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeI64    []int64          `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeInt    []int            `param:"matrixexint,in=path,explode,style=matrix"`
	MatrixExplodeF32    []float32        `param:"matrixexfloat,in=path,explode,style=matrix"`
	MatrixExplodeF64    []float64        `param:"matrixexfloat,in=path,explode,style=matrix"`
	MatrixExplodeStruct TestStructParams `param:"matrixexstruct,in=path,explode,style=matrix"`
}

const (
	testValidSimplePath                = "/{simpstr}/{simpuint}/{simpint}/{simpfloat}/{simpstruct}/{simpexstr}/{simpexuint}/{simpexint}/{simpexfloat}/{simpexstruct}"
	testValidLabelPath                 = "/{labelstr}/{labeluint}/{labelint}/{labelfloat}/{labelstruct}/{labelexstr}/{labelexuint}/{labelexint}/{labelexfloat}/{labelexstruct}"
	testValidMatrixPath                = "/{matrixstr}/{matrixuint}/{matrixint}/{matrixfloat}/{matrixstruct}/{matrixexstr}/{matrixexuint}/{matrixexint}/{matrixexfloat}/{matrixexstruct}"
	testValidPrimitiveSimplePathValues = "/ateststring/123/-123/123.45/NA/test.../255/0/-123.45/NA"
	testValidPrimitiveLabelPathValues  = "/.ateststring/.123/.-123/.123.45/NA/.test.../.255/.0/.-123.45/NA"
	testValidPrimitiveMatrixPathValues = "/;matrixstr=ateststring/;matrixuint=123/;matrixint=-123/;matrixfloat=123.45/NA/;matrixexstr=test.../;matrixexuint=255/;matrixexint=0/;matrixexfloat=-123.45/NA"
	testValidComplexSimplePathValues   = "/string1,string2/0,123/123,-123/123.45,0.0/stringprop,propstring,intprop,123/test.../255/0/-123.45/stringprop=propstring,intprop=123"
	testValidComplexLabelPathValues    = "/.string1,string2/.0,123/.123,-123/.123.45,0.0/.stringprop,propstring,intprop,123/.string3.string4/.255.0/.0.-123/NA/.stringprop=propstring.intprop=123"
	testValidComplexMatrixPathValues   = "/;matrixstr=string1,string2/;matrixuint=0,123/;matrixint=123,-123/;matrixfloat=123.45,0.0/;matrixstruct=stringprop,propstring,intprop,123/;matrixexstr=string3;matrixexstr=string4/;matrixexuint=255;matrixexuint=0/;matrixexint=0;matrixexint=-123/;matrixexfloat=-123.45;matrixexfloat=1/;stringprop=propstring;intprop=123"
)

var (
	testPrimitivePathParams = TestPrimitivePathParams{
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

		LabelStr:  "ateststring",
		LabelU8:   123,
		LabelU16:  123,
		LabelU32:  123,
		LabelU64:  123,
		LabelUint: 123,
		LabelI8:   -123,
		LabelI16:  -123,
		LabelI32:  -123,
		LabelI64:  -123,
		LabelInt:  -123,
		LabelF32:  123.45,
		LabelF64:  123.45,

		LabelExplodeStr:  "test...",
		LabelExplodeU8:   255,
		LabelExplodeU16:  255,
		LabelExplodeU32:  255,
		LabelExplodeU64:  255,
		LabelExplodeUint: 255,
		LabelExplodeI8:   0,
		LabelExplodeI16:  0,
		LabelExplodeI32:  0,
		LabelExplodeI64:  0,
		LabelExplodeInt:  0,
		LabelExplodeF32:  -123.45,
		LabelExplodeF64:  -123.45,

		MatrixStr:  "ateststring",
		MatrixU8:   123,
		MatrixU16:  123,
		MatrixU32:  123,
		MatrixU64:  123,
		MatrixUint: 123,
		MatrixI8:   -123,
		MatrixI16:  -123,
		MatrixI32:  -123,
		MatrixI64:  -123,
		MatrixInt:  -123,
		MatrixF32:  123.45,
		MatrixF64:  123.45,

		MatrixExplodeStr:  "test...",
		MatrixExplodeU8:   255,
		MatrixExplodeU16:  255,
		MatrixExplodeU32:  255,
		MatrixExplodeU64:  255,
		MatrixExplodeUint: 255,
		MatrixExplodeI8:   0,
		MatrixExplodeI16:  0,
		MatrixExplodeI32:  0,
		MatrixExplodeI64:  0,
		MatrixExplodeInt:  0,
		MatrixExplodeF32:  -123.45,
		MatrixExplodeF64:  -123.45,
	}

	testComplexPathParams = TestComplexPathParams{
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

		LabelStr:  []string{"string1", "string2"},
		LabelU8:   []uint8{0, 123},
		LabelU16:  []uint16{0, 123},
		LabelU32:  []uint32{0, 123},
		LabelU64:  []uint64{0, 123},
		LabelUint: []uint{0, 123},
		LabelI8:   []int8{123, -123},
		LabelI16:  []int16{123, -123},
		LabelI32:  []int32{123, -123},
		LabelI64:  []int64{123, -123},
		LabelInt:  []int{123, -123},
		LabelF32:  []float32{123.45, 0.0},
		LabelF64:  []float64{123.45, 0.0},
		LabelStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		LabelExplodeStr:  []string{"string3", "string4"},
		LabelExplodeU8:   []uint8{255, 0},
		LabelExplodeU16:  []uint16{255, 0},
		LabelExplodeU32:  []uint32{255, 0},
		LabelExplodeU64:  []uint64{255, 0},
		LabelExplodeUint: []uint{255, 0},
		LabelExplodeI8:   []int8{0, -123},
		LabelExplodeI16:  []int16{0, -123},
		LabelExplodeI32:  []int32{0, -123},
		LabelExplodeI64:  []int64{0, -123},
		LabelExplodeInt:  []int{0, -123},
		LabelExplodeStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		MatrixStr:  []string{"string1", "string2"},
		MatrixU8:   []uint8{0, 123},
		MatrixU16:  []uint16{0, 123},
		MatrixU32:  []uint32{0, 123},
		MatrixU64:  []uint64{0, 123},
		MatrixUint: []uint{0, 123},
		MatrixI8:   []int8{123, -123},
		MatrixI16:  []int16{123, -123},
		MatrixI32:  []int32{123, -123},
		MatrixI64:  []int64{123, -123},
		MatrixInt:  []int{123, -123},
		MatrixF32:  []float32{123.45, 0.0},
		MatrixF64:  []float64{123.45, 0.0},
		MatrixStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},

		MatrixExplodeStr:  []string{"string3", "string4"},
		MatrixExplodeU8:   []uint8{255, 0},
		MatrixExplodeU16:  []uint16{255, 0},
		MatrixExplodeU32:  []uint32{255, 0},
		MatrixExplodeU64:  []uint64{255, 0},
		MatrixExplodeUint: []uint{255, 0},
		MatrixExplodeI8:   []int8{0, -123},
		MatrixExplodeI16:  []int16{0, -123},
		MatrixExplodeI32:  []int32{0, -123},
		MatrixExplodeI64:  []int64{0, -123},
		MatrixExplodeInt:  []int{0, -123},
		MatrixExplodeF32:  []float32{-123.45, 1},
		MatrixExplodeF64:  []float64{-123.45, 1},
		MatrixExplodeStruct: TestStructParams{
			StringProp: "propstring",
			IntProp:    123,
		},
	}
)
