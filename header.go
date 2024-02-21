package chimera

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var (
	headerParamUnmarshalerType = reflect.TypeOf((*HeaderParamUnmarshaler)(nil)).Elem()
	headerParamMarshalerType   = reflect.TypeOf((*HeaderParamMarshaler)(nil)).Elem()
)

// HeaderParamUnmarshaler is an interface that supports converting a header value ([]string) to a user-defined type
type HeaderParamUnmarshaler interface {
	UnmarshalHeaderParam(value []string, info ParamStructTag) error
}

// HeaderParamMarshaler is an interface that supports converting a user-defined type to a header value ([]string)
type HeaderParamMarshaler interface {
	MarshalHeaderParam(ParamStructTag) (http.Header, error)
}

// unmarshalHeaderParam converts a []string to a value using the options in tag
func unmarshalHeaderParam(param []string, tag *ParamStructTag, addr reflect.Value) error {
	addr = fixPointer(addr)
	if tag.schemaType == interfaceType {
		return addr.Interface().(HeaderParamUnmarshaler).UnmarshalHeaderParam(param, *tag)
	}
	if tag.Required && len(param) == 0 {
		return NewRequiredParamError("header", tag.Name)
	}
	if len(param) == 0 {
		return nil
	}
	return unmarshalStringParam(param[0], tag, addr)
}

// marshalHeaderParam converts a value to a http.Header using the options in tag
func marshalHeaderParam(tag *ParamStructTag, addr reflect.Value) (http.Header, error) {
	addr = fixPointer(addr)
	switch tag.schemaType {
	case interfaceType:
		return addr.Interface().(HeaderParamMarshaler).MarshalHeaderParam(*tag)
	case primitiveType:
		return http.Header{
			tag.Name: []string{marshalPrimitiveToString(addr)},
		}, nil
	case sliceType:
		return http.Header{
			tag.Name: []string{marshalSliceToString(addr)},
		}, nil
	case structType:
		return http.Header{
			tag.Name: []string{marshalStructToString(addr, tag)},
		}, nil
	}
	return nil, nil
}

// marshalPrimitiveToString converts a value to a string
func marshalPrimitiveToString(addr reflect.Value) string {
	return fmt.Sprint(addr.Elem().Interface())
}

// marshalSliceToString converts a slice/array value to a string
func marshalSliceToString(addr reflect.Value) string {
	value := ""
	addr = addr.Elem()
	for i := 0; i < addr.Len(); i++ {
		if i != 0 {
			value += ","
		}
		value += marshalPrimitiveToString(addr.Index(i).Addr())
	}
	return value
}

// marshalStructToString converts a struct value to a string
func marshalStructToString(addr reflect.Value, tag *ParamStructTag) string {
	values := make([]string, 0)
	for name, prop := range tag.propMap {
		f := addr.Elem().Field(prop.fieldIndex)
		if f.Type().Kind() == reflect.Pointer {
			f = fixPointer(f)
		}
		v := name + tag.valueDelim
		v += marshalPrimitiveToString(f.Addr())
		values = append(values, v)
	}
	return strings.Join(values, tag.delim)
}
