package chimera

import (
	"net/url"
	"reflect"
	"strings"
)

var (
	queryParamUnmarshalerType = reflect.TypeOf((*QueryParamUnmarshaler)(nil)).Elem()
)

// QueryParamUnmarshaler allows a type to add custom validation or parsing logic based on query parameters
type QueryParamUnmarshaler interface {
	UnmarshalQueryParam(value url.Values, info ParamStructTag) error
}

// unmarshalQueryParam attempts to turn query values into a value
func unmarshalQueryParam(param url.Values, tag *ParamStructTag, addr reflect.Value) error {
	// TODO: we need to find a way to handle multiple nested query params
	// its hard because form encoded bodies cant repeat keys but query strings
	// are the wild west. at this point, I am just going to follow the OpenAPI
	// standard but it may be preferable later to support multiple options like
	// X.Y.Z.0 or X[Y][Z][0] or X.Y.Z[0] or just always treat repeats as arrays
	// may also want to disable the use of explode for form styles
	addr = fixPointer(addr)
	if tag.schemaType == interfaceType {
		return addr.Interface().(QueryParamUnmarshaler).UnmarshalQueryParam(param, *tag)
	}
	switch tag.schemaType {
	case sliceType:
		return unmarshalSliceFromQuery(param, tag, addr)
	case primitiveType:
		if val, ok := param[tag.Name]; ok {
			if len(val) > 0 {
				return unmarshalPrimitiveFromString(val[0], tag, addr)
			}
		} else if tag.Required {
			return NewRequiredParamError("query", tag.Name)
		}
	case structType:
		return unmarshalStructFromQuery(param, tag, addr)
	}
	return nil
}

// unmarshalSliceFromQuery unmarshals slice types from query values
func unmarshalSliceFromQuery(param url.Values, tag *ParamStructTag, addr reflect.Value) error {
	value := addr.Elem()
	eType := value.Type().Elem()
	value.Set(reflect.MakeSlice(reflect.SliceOf(eType), 0, 0))
	vals, ok := param[tag.Name]
	if tag.Required && (!ok || len(vals) == 0) {
		return NewRequiredParamError("query", tag.Name)
	}
	if len(vals) == 0 {
		return nil
	}
	if !tag.Explode {
		vals = strings.Split(vals[0], tag.delim)
	}

	for _, v := range vals {
		val, err := decodePrimitiveString(v, eType.Kind())
		if err != nil {
			return NewInvalidParamError(marshalIn(tag.In), tag.Name, v)
		}
		value = reflect.Append(value, val)
	}
	addr.Elem().Set(value)
	return nil
}

// unmarshalStructFromQuery unmarshals struct types from query values
func unmarshalStructFromQuery(param url.Values, tag *ParamStructTag, addr reflect.Value) error {
	switch tag.Style {
	case FormStyle:
		if tag.Explode {
			requiredCheck := !tag.Required
			for name, prop := range tag.propMap {
				switch prop.schemaType {
				case primitiveType:
					if val, ok := param[name]; ok && len(val) > 0 {
						f := addr.Elem().Field(prop.fieldIndex)
						v, err := decodePrimitiveString(val[0], f.Kind())
						if err != nil {
							return NewInvalidParamError(marshalIn(tag.In), tag.Name, val[0])
						}
						f.Set(v)
						requiredCheck = true
					}
					// case Interface:
					// 	if val, ok := param[name]; ok && len(val) > 0 {
					// 		f := addr.Elem().Field(prop.fieldIndex)
					// 		err := f.Interface().(ParamPropUnmarshaler).UnmarshalParamProp(val[0])
					// 		if err != nil {
					// 			return nil
					// 		}
					// 		requiredCheck = true
					// 	}
				}
			}
			if !requiredCheck {
				return NewRequiredParamError("query", tag.Name)
			}
		} else {
			if val, ok := param[tag.Name]; ok {
				if len(val) > 0 {
					return unmarshalStructFromString(param.Get(tag.Name), tag, addr)
				}
			} else if tag.Required {
				return NewRequiredParamError("query", tag.Name)
			}
		}
	case DeepObjectStyle:
		requiredCheck := !tag.Required
		for name, prop := range tag.propMap {
			f := addr.Elem().Field(prop.fieldIndex)
			if val, ok := param[name]; ok && len(val) > 0 {
				v, err := decodePrimitiveString(val[0], f.Kind())
				if err != nil {
					return NewInvalidParamError(marshalIn(tag.In), tag.Name, val[0])
				}
				f.Set(v)
				requiredCheck = true
			}
		}
		if !requiredCheck {
			return NewRequiredParamError("query", tag.Name)
		}
	}
	return nil
}

// normalizeQueryStyle ensures the style is appropriate for a query param
func normalizeQueryStyle(style Style) Style {
	if style == PipeDelimitedStyle || style == FormStyle || style == SpaceDelimitedStyle || style == DeepObjectStyle {
		return style
	}
	return FormStyle
}
