package chimera

import (
	"net/http"
	"reflect"
)

var (
	cookieParamUnmarshalerType = reflect.TypeOf((*CookieParamUnmarshaler)(nil)).Elem()
	cookieParamMarshalerType   = reflect.TypeOf((*CookieParamMarshaler)(nil)).Elem()
)

// CookieParamUnmarshaler is an interface that supports converting an http.Cookie to a user-defined type
type CookieParamUnmarshaler interface {
	UnmarshalCookieParam(http.Cookie) error
}

// CookieParamMarshaler is an interface that supports converting a user-defined type to an http.Cookie
type CookieParamMarshaler interface {
	MarshalCookieParam(ParamStructTag) (http.Cookie, error)
}

// unmarshalCookieParam converts a cookie to a value using the options in tag
func unmarshalCookieParam(param http.Cookie, tag *ParamStructTag, addr reflect.Value) error {
	addr = fixPointer(addr)
	if tag.schemaType == interfaceType {
		return addr.Interface().(CookieParamUnmarshaler).UnmarshalCookieParam(param)
	}
	return unmarshalStringParam(param.Value, tag, addr)
}

// unmarshalCookieParam converts a value to a http.Cookie using the options in tag
func marshalCookieParam(tag *ParamStructTag, addr reflect.Value) (http.Cookie, error) {
	addr = fixPointer(addr)
	switch tag.schemaType {
	case interfaceType:
		return addr.Interface().(CookieParamMarshaler).MarshalCookieParam(*tag)
	case primitiveType:
		return http.Cookie{
			Name:  tag.Name,
			Value: marshalPrimitiveToString(addr),
		}, nil
	case sliceType:
		return http.Cookie{
			Name:  tag.Name,
			Value: marshalSliceToString(addr),
		}, nil
	case structType:
		return http.Cookie{
			Name:  tag.Name,
			Value: marshalStructToString(addr, tag),
		}, nil
	}
	return http.Cookie{}, nil
}
