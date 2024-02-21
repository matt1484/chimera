package chimera

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	pathParamUnmarshalerType = reflect.TypeOf((*PathParamUnmarshaler)(nil)).Elem()
)

// normalizePathStyle forces a style to match one of the allowed openapi types
func normalizePathStyle(style Style) Style {
	if style == SimpleStyle || style == LabelStyle || style == MatrixStyle {
		return style
	}
	return SimpleStyle
}

// decodePrimitiveString turns a string into a primitive reflect value
func decodePrimitiveString(value string, kind reflect.Kind) (reflect.Value, error) {
	// TODO: support time, ip, email, url, etc?
	switch kind {
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		return reflect.ValueOf(v), err
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(int8))), err
	case reflect.Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(int16))), err
	case reflect.Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(int32))), err
	case reflect.Int, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if kind == reflect.Int64 {
			return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(int64))), err
		}
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(int))), err
	case reflect.Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(uint8))), err
	case reflect.Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(uint16))), err
	case reflect.Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(uint32))), err
	case reflect.Uint, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if kind == reflect.Uint64 {
			return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(uint64))), err
		}
		return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(uint))), err
	case reflect.Float32, reflect.Float64:
		var v float64
		var err error
		if kind == reflect.Float32 {
			v, err = strconv.ParseFloat(value, 32)
			return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(float32))), err
		}
		v, err = strconv.ParseFloat(value, 64)
		return reflect.ValueOf(v), err
	case reflect.Complex64, reflect.Complex128:
		var v complex128
		var err error
		if kind == reflect.Complex64 {
			v, err = strconv.ParseComplex(value, 64)
			return reflect.ValueOf(v).Convert(reflect.TypeOf(*new(complex64))), err
		}
		v, err = strconv.ParseComplex(value, 128)
		return reflect.ValueOf(v), err
	}
	return reflect.ValueOf(nil), errors.New("unable to convert string to kind: " + kind.String())
}

// cutPrefix trims the prefix from a string
func cutPrefix(raw, prefix string) (string, error) {
	if prefix == "" {
		return raw, nil
	}
	if len(raw) < len(prefix) || raw[:len(prefix)] != prefix {
		return raw, nil
	}
	return raw[len(prefix):], nil
}

// unmarshalPathParam converts a path param (string) into a usable value
func unmarshalPathParam(param string, tag *ParamStructTag, addr reflect.Value) error {
	addr = fixPointer(addr)
	if tag.schemaType == interfaceType {
		return addr.Interface().(PathParamUnmarshaler).UnmarshalPathParam(param, *tag)
	}
	return unmarshalStringParam(param, tag, addr)
}

// unmarshalStringParam parses string params
func unmarshalStringParam(param string, tag *ParamStructTag, addr reflect.Value) error {
	switch tag.schemaType {
	case sliceType:
		return unmarshalSliceFromString(param, tag, addr)
	case primitiveType:
		return unmarshalPrimitiveFromString(param, tag, addr)
	case structType:
		return unmarshalStructFromString(param, tag, addr)
	}
	return nil
}

// unmarshalPrimitiveFromString unmarshals a primitive value from a string value
func unmarshalPrimitiveFromString(param string, tag *ParamStructTag, addr reflect.Value) error {
	param, err := cutPrefix(param, tag.prefix)
	if err != nil {
		return err
	}
	val, err := decodePrimitiveString(param, addr.Elem().Kind())
	addr.Elem().Set(val)
	if err != nil {
		err = NewInvalidParamError(marshalIn(tag.In), tag.Name, param)
	}
	return err
}

// unmarshalSliceFromString unmarshals a slice value from a string value
func unmarshalSliceFromString(param string, tag *ParamStructTag, addr reflect.Value) error {
	value := addr.Elem()
	eType := value.Type().Elem()
	value.Set(reflect.MakeSlice(reflect.SliceOf(eType), 0, 0))
	param, err := cutPrefix(param, tag.prefix)
	if err != nil {
		return err
	}
	for _, v := range strings.Split(param, tag.delim) {
		val, err := decodePrimitiveString(v, eType.Kind())
		if err != nil {
			return NewInvalidParamError(marshalIn(tag.In), tag.Name, param)
		}
		value = reflect.Append(value, val)
	}
	addr.Elem().Set(value)
	return nil
}

// PathParamUnmarshaler is used to allow types to implement their own logic for parsing path parameters
type PathParamUnmarshaler interface {
	UnmarshalPathParam(param string, info ParamStructTag) error
}

// type ParamPropUnmarshaler interface {
// 	UnmarshalParamProp(prop string) error
// }

// var (
// 	paramPropUnmarshalerType = reflect.TypeOf((*ParamPropUnmarshaler)(nil)).Elem()
// )

// propsFromString gets the struct properties from a string
// NOTE: this is largely based on of kin-openapi, but may need to change later
func propsFromString(src, propDelim, valueDelim string) (map[string]string, error) {
	props := make(map[string]string)
	pairs := strings.Split(src, propDelim)

	// When propDelim and valueDelim is equal the source string follow the next rule:
	// every even item of pairs is a properties's name, and the subsequent odd item is a property's value.
	if propDelim == valueDelim {
		// Taking into account the rule above, a valid source string must be splitted by propDelim
		// to an array with an even number of items.
		if len(pairs)%2 != 0 {
			return nil, fmt.Errorf("a value must be a list of object's properties in format \"name%svalue\" separated by %s", valueDelim, propDelim)
		}
		for i := 0; i < len(pairs)/2; i++ {
			props[pairs[i*2]] = pairs[i*2+1]
		}
		return props, nil
	}

	// When propDelim and valueDelim is not equal the source string follow the next rule:
	// every item of pairs is a string that follows format <propName><valueDelim><propValue>.
	for _, pair := range pairs {
		prop := strings.Split(pair, valueDelim)
		if len(prop) != 2 {
			return nil, fmt.Errorf("a value must be a list of object's properties in format \"name%svalue\" separated by %s", valueDelim, propDelim)
		}
		props[prop[0]] = prop[1]
	}
	return props, nil
}

// unmarshalStructFromString unmarshals a struct value from a string value
func unmarshalStructFromString(param string, tag *ParamStructTag, addr reflect.Value) error {
	par, err := cutPrefix(param, tag.prefix)
	if err != nil {
		return err
	}
	props, err := propsFromString(par, tag.delim, tag.valueDelim)
	if err != nil {
		return NewInvalidParamError(marshalIn(tag.In), tag.Name, param)
	}
	for name, prop := range tag.propMap {
		if valStr, ok := props[name]; ok {
			f := addr.Elem().Field(prop.fieldIndex)
			switch prop.schemaType {
			case primitiveType:
				v, err := decodePrimitiveString(valStr, f.Kind())
				if err != nil {
					return NewInvalidParamError(marshalIn(tag.In), tag.Name, param)
				}
				f.Set(v)
				// case Interface:
				// 	err = f.Interface().(ParamPropUnmarshaler).UnmarshalParamProp(valStr)
				// 	if err != nil {
				// 		return err
				// 	}
			}
		}
	}
	return nil
}
