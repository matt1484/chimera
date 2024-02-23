package chimera

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/invopop/jsonschema"
	"github.com/matt1484/spectagular"
)

// SchemaType describes the type of a schema
type SchemaType int

const (
	structType SchemaType = iota
	primitiveType
	sliceType
	interfaceType
)

// Style denotes the openapi "style" of a parameter
type Style int

const (
	DefaultStyle Style = iota
	SimpleStyle
	LabelStyle
	MatrixStyle
	SpaceDelimitedStyle
	PipeDelimitedStyle
	DeepObjectStyle
	FormStyle
	// Should these be styles or something else?
	// Base64Style
	// JSONStyle
	// JWTStyle
)

// In denotes where a paramter lives (i.e. path, query, header, cookie)
type In int

const (
	NullIn In = iota
	PathIn
	QueryIn
	HeaderIn
	CookieIn
)

// var (
// 	paramsUnmarshalerType = reflect.TypeOf((*ParamsUnmarshaler)(nil)).Elem()
// 	paramUnmarshalerType  = reflect.TypeOf((*ParamUnmarshaler)(nil)).Elem()
// )

// type ParamsUnmarshaler interface {
// 	UnmarshalParams(req *http.Request) error
// }

// type ParamUnmarshaler interface {
// 	UnmarshalParam(req *http.Request) error
// }

// ParamStructTag describes the various parts of the "param" struct tag
type ParamStructTag struct {
	Name            string `structtag:"$name"`
	In              In     `structtag:"in,required"`
	Explode         bool   `structtag:"explode"`
	Style           Style  `structtag:"style"`
	Required        bool   `structtag:"required"`
	Description     string `structtag:"description"`
	Deprecated      bool   `structtag:"deprecated"`
	AllowEmptyValue bool   `structtag:"allowEmptyValue"`
	AllowReserved   bool   `structtag:"allowReserved"`
	prefix          string
	delim           string
	valueDelim      string
	schemaType      SchemaType
	propMap         map[string]*paramProp
	schema          *jsonschema.Schema
	request         bool
}

// OpenAPIParameterSpec returns the Parameter definition of a struct tag
func (p *ParamStructTag) OpenAPIParameterSpec() Parameter {
	in := p.In
	name := p.Name
	if !p.request {
		if p.In == CookieIn {
			if p.Description == "" {
				p.Description = name
			}
			name = "Set-Cookie"
		}
		in = HeaderIn
	}
	return Parameter{
		Name:            name,
		In:              marshalIn(in),
		Style:           marshalParamStyle(p.Style),
		Explode:         p.Explode,
		Required:        p.Required,
		Description:     p.Description,
		Deprecated:      p.Deprecated,
		AllowEmptyValue: p.AllowEmptyValue,
		AllowReserved:   p.AllowReserved,
		Schema:          p.schema,
	}
}

// NewRequiredParamError returns an APIError to denote that a parameter was missing
func NewRequiredParamError(in, name string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       []byte(fmt.Sprintf("missing required %s parameter %s", in, name)),
	}
}

// NewInvalidParamError returns an APIError to denote a parameter was improperly formatted
func NewInvalidParamError(in, paramName, value string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       []byte(fmt.Sprintf("%s parameter %s was improperly formatted %v", in, paramName, value)),
	}
}

// unmarshalParamStyle converts a string to a ParamStyle
func unmarshalParamStyle(value string) Style {
	value = strings.ToLower(value)
	switch value {
	case "simple":
		return SimpleStyle
	case "label":
		return LabelStyle
	case "matrix":
		return MatrixStyle
	case "spacedelimited":
		return SpaceDelimitedStyle
	case "pipedelimited":
		return PipeDelimitedStyle
	case "deepobject":
		return DeepObjectStyle
	case "form":
		return FormStyle
	}
	return DefaultStyle
}

// marshalParamStyle converts a ParamStyle to a string
func marshalParamStyle(value Style) string {
	switch value {
	case SimpleStyle:
		return "simple"
	case LabelStyle:
		return "label"
	case MatrixStyle:
		return "matrix"
	case SpaceDelimitedStyle:
		return "spaceDelimited"
	case PipeDelimitedStyle:
		return "pipeDelimited"
	case DeepObjectStyle:
		return "deepObject"
	case FormStyle:
		return "form"
	}
	return ""
}

// unmarshalIn converts a string to an In enum
func unmarshalIn(value string) In {
	value = strings.ToLower(value)
	switch value {
	case "path":
		return PathIn
	case "query":
		return QueryIn
	case "header":
		return HeaderIn
	case "cookie":
		return CookieIn
	}
	return NullIn
}

// marshalIn converts In to a string
func marshalIn(value In) string {
	switch value {
	case PathIn:
		return "path"
	case QueryIn:
		return "query"
	case HeaderIn:
		return "header"
	case CookieIn:
		return "cookie"
	}
	return ""
}

// UnmarshalTagOption is used to unmarshal a Style found in a struct tag
func (p Style) UnmarshalTagOption(field reflect.StructField, value string) (reflect.Value, error) {
	return reflect.ValueOf(unmarshalParamStyle(value)), nil
}

// UnmarshalTagOption is used to unmarshal a In found in a struct tag
func (i In) UnmarshalTagOption(field reflect.StructField, value string) (reflect.Value, error) {
	in := unmarshalIn(value)
	if in == NullIn {
		return reflect.ValueOf(in), errors.New("invalid param 'in': " + value)
	}
	return reflect.ValueOf(in), nil
}

// initTagCache forces the initialization of a tag cache
// NOTE: maybe add a Must-like function in spectagular
func initTagCache[T any](tag string) *spectagular.StructTagCache[T] {
	cache, err := spectagular.NewFieldTagCache[T](tag)
	if err != nil {
		panic("initializing struct tag cache failed")
	}
	return cache
}

// paramProp stores information about a property in a param struct
type paramProp struct {
	schemaType SchemaType
	fieldIndex int
}

// paramPropTag is just an easy way to the first part of a struct tag
type paramPropTag struct {
	Name string `structtag:"$name"`
}

var (
	requestParamTagCache  = initTagCache[ParamStructTag]("param")
	responseParamTagCache = initTagCache[ParamStructTag]("param")
)

// fixPointer initializes pointer objects and returns the lowest level one
// (I really doubt people are out here making types like ****int but Ive seen jank code in my day)
func fixPointer(value reflect.Value) reflect.Value {
	if value.Elem().Type().Kind() == reflect.Pointer {
		i := 0
		t := value.Type()
		for ; t.Kind() == reflect.Pointer; t = t.Elem() {
			i++
		}
		v := reflect.New(t)
		ret := v
		for ; i > 1; i-- {
			p := reflect.New(v.Type())
			p.Elem().Set(v)
			v = p
		}
		value.Elem().Set(v.Elem())
		return ret
	}
	return value
}

// CacheRequestParamsType adds a type to the internal tag cache and returns the resulting Paramter objects
func CacheRequestParamsType(t reflect.Type) []Parameter {
	var params []Parameter
	pTag, err := requestParamTagCache.GetOrAdd(t)
	if err != nil {
		panic("chimera: failed to parse parameter field of type " + t.Name())
	}
	v := reflect.New(t)
	for i, tag := range pTag {
		t := v.Elem().Field(tag.FieldIndex).Addr()
		t = fixPointer(t)
		switch tag.Value.In {
		case PathIn:
			tag.Value.Style = normalizePathStyle(Style(tag.Value.Style))
			tag.Value.Required = true
			if t.Type().Implements(pathParamUnmarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		case QueryIn:
			tag.Value.Style = normalizeQueryStyle(Style(tag.Value.Style))
			if t.Type().Implements(queryParamUnmarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		case HeaderIn:
			// TODO: support JSON/JWT/base64
			tag.Value.Style = SimpleStyle
			if t.Type().Implements(headerParamUnmarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		case CookieIn:
			// TODO: support JSON/JWT/base64
			tag.Value.Style = FormStyle
			if t.Type().Implements(cookieParamUnmarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		}
		if tag.Value.schemaType != interfaceType {
			if t.Elem().Kind() == reflect.Slice {
				tag.Value.schemaType = sliceType
			} else if t.Elem().Kind() == reflect.Struct {
				tag.Value.schemaType = structType
			} else {
				tag.Value.schemaType = primitiveType
			}
		}
		switch tag.Value.schemaType {
		case primitiveType:
			tag.Value.schema = (&jsonschema.Reflector{
				ExpandedStruct: false,
				DoNotReference: true,
				FieldNameTag:   "prop",
			}).Reflect(t.Interface())
			switch Style(tag.Value.Style) {
			case LabelStyle:
				tag.Value.prefix = "."
			case MatrixStyle:
				tag.Value.prefix = ";" + tag.Value.Name + "="
			}
		case structType:
			tag.Value.propMap = make(map[string]*paramProp)
			schema := (&jsonschema.Reflector{
				ExpandedStruct: true,
				FieldNameTag:   "prop",
			}).Reflect(t.Elem().Interface())
			tag.Value.schema = schema
			fieldTags, _ := spectagular.ParseTagsForType[paramPropTag]("prop", t.Elem().Type())
			for _, ft := range fieldTags {
				v := t.Elem().Field(ft.FieldIndex).Addr()
				v = fixPointer(v)
				field := t.Elem().Type().Field(ft.FieldIndex)
				if field.Anonymous {
					continue
				}
				jsonSchemaTags := strings.Split(field.Tag.Get("jsonschema"), ",")
				if jsonSchemaTags[0] == "-" {
					continue
				}
				name := ft.Value.Name
				// should we only unmarshal props that are in the spec?
				if _, ok := schema.Properties.Get(name); !ok {
					continue
				}
				if tag.Value.Style == DeepObjectStyle {
					name = tag.Value.Name + "[" + name + "]"
				}
				tag.Value.propMap[name] = &paramProp{
					fieldIndex: ft.FieldIndex,
				}
				// if v.Type().Implements(paramPropUnmarshalerType) {
				// 	tag.Value.propMap[name].schemaType = Interface
				// } else
				if v.Elem().Kind() != reflect.Struct &&
					v.Elem().Kind() != reflect.Array &&
					v.Elem().Kind() != reflect.Slice &&
					v.Elem().Kind() != reflect.Chan &&
					v.Elem().Kind() != reflect.Interface &&
					v.Elem().Kind() != reflect.Map &&
					v.Elem().Kind() != reflect.UnsafePointer &&
					v.Elem().Kind() != reflect.Ptr {
					tag.Value.propMap[name].schemaType = primitiveType
				}
			}
			switch tag.Value.Style {
			case SimpleStyle:
				tag.Value.delim = ","
				if tag.Value.Explode {
					tag.Value.valueDelim = "="
				} else {
					tag.Value.valueDelim = ","
				}
			case LabelStyle:
				tag.Value.prefix = "."
				if tag.Value.Explode {
					tag.Value.delim = "."
					tag.Value.valueDelim = "="
				} else {
					tag.Value.delim = ","
					tag.Value.valueDelim = ","
				}
			case MatrixStyle:
				tag.Value.prefix = ";"
				if tag.Value.Explode {
					tag.Value.delim = ";"
					tag.Value.valueDelim = "="
				} else {
					tag.Value.prefix += tag.Value.Name + "="
					tag.Value.delim = ","
					tag.Value.valueDelim = ","
				}
			case FormStyle:
				if tag.Value.Explode {
					tag.Value.valueDelim = "="
					tag.Value.delim = "&"
				} else {
					tag.Value.valueDelim = ","
					tag.Value.delim = ","
				}
			}
		case sliceType:
			tag.Value.schema = (&jsonschema.Reflector{
				ExpandedStruct: false,
				FieldNameTag:   "prop",
			}).Reflect(t.Interface())
			switch Style(tag.Value.Style) {
			case SimpleStyle:
				tag.Value.delim = ","
			case LabelStyle:
				tag.Value.prefix = "."
				if tag.Value.Explode {
					tag.Value.delim = "."
				} else {
					tag.Value.delim = ","
				}
			case MatrixStyle:
				tag.Value.prefix = ";" + tag.Value.Name + "="
				if tag.Value.Explode {
					tag.Value.delim = ";" + tag.Value.Name + "="
				} else {
					tag.Value.delim = ","
				}
			case FormStyle:
				// TODO: flatten schemas
				tag.Value.prefix = tag.Value.Name + "="
				if tag.Value.Explode {
					tag.Value.delim = "&" + tag.Value.Name + "="
				} else {
					tag.Value.prefix = ""
					tag.Value.delim = ","
				}
			case PipeDelimitedStyle:
				tag.Value.prefix = tag.Value.Name + "="
				if tag.Value.Explode {
					tag.Value.delim = "&" + tag.Value.Name + "="
				} else {
					tag.Value.delim = "|"
				}
			case SpaceDelimitedStyle:
				tag.Value.prefix = tag.Value.Name + "="
				if tag.Value.Explode {
					tag.Value.delim = "&" + tag.Value.Name + "="
				} else {
					tag.Value.delim = " "
				}
			}
		}
		tag.Value.request = true
		pTag[i] = tag
		params = append(params, tag.Value.OpenAPIParameterSpec())
	}
	return params
}

// UnmarshalParams gets all the parameters out of a request object
// (headers, cookies, query, path)
func UnmarshalParams(request *http.Request, obj any) error {
	value := reflect.ValueOf(obj).Elem()
	paramType := value.Type()
	paramTags, found := requestParamTagCache.Get(paramType)
	if !found {
		CacheRequestParamsType(paramType)
		paramTags, _ = requestParamTagCache.GetOrAdd(paramType)
	}
	reqCtx := chi.RouteContext(request.Context())
	for _, tag := range paramTags {
		addr := value.Field(tag.FieldIndex).Addr()
		switch tag.Value.In {
		case PathIn:
			if err := unmarshalPathParam(reqCtx.URLParam(tag.Value.Name), &tag.Value, addr); err != nil {
				return err
			}
		case HeaderIn:
			if err := unmarshalHeaderParam(request.Header.Values(tag.Value.Name), &tag.Value, addr); err != nil {
				return err
			}
		case CookieIn:
			cookie, err := request.Cookie(tag.Value.Name)
			if err != nil {
				if err == http.ErrNoCookie && tag.Value.Required {
					return NewRequiredParamError("cookie", tag.Value.Name)
				}
				return err
			}
			if cookie != nil {
				// TODO: support raw cookie parsing? not sure how useful that is
				if err := unmarshalCookieParam(*cookie, &tag.Value, addr); err != nil {
					return err
				}
			} else if tag.Value.Required {
				return NewRequiredParamError("cookie", tag.Value.Name)
			}
		case QueryIn:
			if err := unmarshalQueryParam(request.URL.Query(), &tag.Value, addr); err != nil {
				return err
			}
		}
	}
	return nil
}

// CacheResponseParamsType adds a type to the internal tag cache and returns the resulting Paramter objects
func CacheResponseParamsType(t reflect.Type) []Parameter {
	var params []Parameter
	pTag, err := responseParamTagCache.GetOrAdd(t)
	if err != nil {
		panic("chimera: failed to parse parameter field of type " + t.Name())
	}
	v := reflect.New(t)
	cookieCount := 0
	for i, tag := range pTag {
		t := v.Elem().Field(tag.FieldIndex).Addr()
		t = fixPointer(t)
		switch tag.Value.In {
		case PathIn:
		case QueryIn:
			// TODO: maybe throw an error here? not sure how strict I should be
			continue
		case HeaderIn:
			// TODO: support JSON/JWT/base64
			tag.Value.Style = SimpleStyle
			if t.Type().Implements(headerParamMarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		case CookieIn:
			// TODO: support JSON/JWT/base64
			tag.Value.Style = FormStyle
			if t.Type().Implements(cookieParamMarshalerType) {
				tag.Value.schemaType = interfaceType
			}
		}
		if tag.Value.schemaType != interfaceType {
			if t.Elem().Kind() == reflect.Slice {
				tag.Value.schemaType = sliceType
			} else if t.Elem().Kind() == reflect.Struct {
				tag.Value.schemaType = structType
			} else {
				tag.Value.schemaType = primitiveType
			}
		}
		switch tag.Value.schemaType {
		case primitiveType:
			tag.Value.schema = (&jsonschema.Reflector{
				ExpandedStruct: false,
				DoNotReference: true,
				FieldNameTag:   "prop",
			}).Reflect(t.Interface())
		case structType:
			tag.Value.propMap = make(map[string]*paramProp)
			if tag.Value.Style == DeepObjectStyle {
				tag.Value.Explode = true
			}
			schema := (&jsonschema.Reflector{
				ExpandedStruct: true,
				FieldNameTag:   "prop",
			}).Reflect(t.Elem().Interface())
			tag.Value.schema = schema
			fieldTags, _ := spectagular.ParseTagsForType[paramPropTag]("prop", t.Elem().Type())
			for _, ft := range fieldTags {
				v := t.Elem().Field(ft.FieldIndex).Addr()
				v = fixPointer(v)
				field := t.Elem().Type().Field(ft.FieldIndex)
				if field.Anonymous {
					continue
				}
				jsonSchemaTags := strings.Split(field.Tag.Get("jsonschema"), ",")
				if jsonSchemaTags[0] == "-" {
					continue
				}
				name := ft.Value.Name
				// should we only unmarshal props that are in the spec?
				if _, ok := schema.Properties.Get(name); !ok {
					continue
				}
				if tag.Value.Style == DeepObjectStyle {
					name = tag.Value.Name + "[" + name + "]"
				}
				tag.Value.propMap[name] = &paramProp{
					fieldIndex: ft.FieldIndex,
				}
				if v.Elem().Kind() != reflect.Struct &&
					v.Elem().Kind() != reflect.Array &&
					v.Elem().Kind() != reflect.Slice &&
					v.Elem().Kind() != reflect.Chan &&
					v.Elem().Kind() != reflect.Interface &&
					v.Elem().Kind() != reflect.Map &&
					v.Elem().Kind() != reflect.UnsafePointer &&
					v.Elem().Kind() != reflect.Ptr {
					tag.Value.propMap[name].schemaType = primitiveType
				}
			}
			switch tag.Value.Style {
			case SimpleStyle:
				tag.Value.delim = ","
				if tag.Value.Explode {
					tag.Value.valueDelim = "="
				} else {
					tag.Value.valueDelim = ","
				}
			case FormStyle:
				if tag.Value.Explode {
					tag.Value.valueDelim = "="
					tag.Value.delim = "&"
				} else {
					tag.Value.valueDelim = ","
					tag.Value.delim = ","
				}
			}
		case sliceType:
			tag.Value.schema = (&jsonschema.Reflector{
				ExpandedStruct: false,
				FieldNameTag:   "prop",
			}).Reflect(t.Interface())
			switch Style(tag.Value.Style) {
			case SimpleStyle:
				tag.Value.delim = ","
			case FormStyle:
				tag.Value.prefix = tag.Value.Name + "="
				if tag.Value.Explode {
					tag.Value.delim = "&" + tag.Value.Name + "="
				} else {
					tag.Value.prefix = ""
					tag.Value.delim = ","
				}
			}
		}
		pTag[i] = tag
		param := tag.Value.OpenAPIParameterSpec()
		if tag.Value.In == CookieIn {
			param.Name = strings.Repeat("\000", cookieCount) + param.Name
			cookieCount++
		}
		params = append(params, param)
	}
	return params
}

// MarshalParams turns an object into headers
// this technically supports cookies and headers but the result is all headers
func MarshalParams(obj any) (http.Header, error) {
	value := reflect.ValueOf(obj).Elem()
	paramType := value.Type()
	paramTags, found := responseParamTagCache.Get(paramType)
	if !found {
		CacheResponseParamsType(paramType)
		paramTags, _ = responseParamTagCache.GetOrAdd(paramType)
	}
	header := http.Header{}
	for _, tag := range paramTags {
		addr := value.Field(tag.FieldIndex).Addr()
		switch tag.Value.In {
		case HeaderIn:
			h, err := marshalHeaderParam(&tag.Value, addr)
			if err != nil {
				return nil, err
			}
			for n, v := range h {
				header[n] = append(header[n], v...)
			}
		case CookieIn:
			cookie, err := marshalCookieParam(&tag.Value, addr)
			if err != nil {
				return nil, err
			}
			header.Add("set-cookie", cookie.String())
		}
	}
	return header, nil
}
