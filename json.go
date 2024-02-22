package chimera

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
)

var (
	_ RequestReader  = new(JSONRequest[Nil, Nil])
	_ ResponseWriter = new(JSONResponse[Nil, Nil])
)

// JSONRequest[Body, Params any] is a request type that decodes json request bodies to a
// user-defined struct for the Body and Params
type JSONRequest[Body, Params any] struct {
	request *http.Request
	Body    Body
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *JSONRequest[Body, Params]) Context() context.Context {
	return r.request.Context()
}

// ReadRequest reads the body of an http request and assigns it to the Body field using json.Unmarshal
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *JSONRequest[Body, Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	r.Body = *new(Body)
	if _, ok := any(r.Body).(Nil); !ok {
		err = json.Unmarshal(body, &r.Body)
		if err != nil {
			return err
		}
	}

	r.Params = *new(Params)
	if _, ok := any(r.Params).(Nil); !ok {
		err = UnmarshalParams(req, &r.Params)
		if err != nil {
			return err
		}
	}
	r.request = req
	return nil
}

// standardizedSchemas basically tries to convert all jsonschema.Schema objects to be mapped
// in defs and then replace them with $refs
func standardizedSchemas(schema *jsonschema.Schema, defs map[string]jsonschema.Schema) {
	if schema == nil {
		return
	}
	schema.Version = ""
	standardizedSchemas(schema.AdditionalProperties, defs)
	standardizedSchemas(schema.Contains, defs)
	standardizedSchemas(schema.ContentSchema, defs)
	standardizedSchemas(schema.Else, defs)
	standardizedSchemas(schema.If, defs)
	standardizedSchemas(schema.Items, defs)
	standardizedSchemas(schema.Not, defs)
	standardizedSchemas(schema.PropertyNames, defs)
	standardizedSchemas(schema.Then, defs)
	for _, s := range schema.AllOf {
		standardizedSchemas(s, defs)
	}
	for _, s := range schema.AnyOf {
		standardizedSchemas(s, defs)
	}
	for _, s := range schema.OneOf {
		standardizedSchemas(s, defs)
	}
	for _, s := range schema.DependentSchemas {
		standardizedSchemas(s, defs)
	}
	for _, s := range schema.PatternProperties {
		standardizedSchemas(s, defs)
	}
	for k, s := range schema.Definitions {
		standardizedSchemas(s, defs)
		if _, ok := defs[k]; ok {
			split := strings.Split(s.ID.String(), "/")
			for i, p := range split {
				if i != len(split)-1 {
					k = p + "_" + k
				}
			}
		}
		defs[k] = *s
	}
	schema.Definitions = nil
	for p := schema.Properties.Oldest(); p != nil; p = p.Next() {
		standardizedSchemas(p.Value, defs)
	}
	if strings.HasPrefix(schema.Ref, "#/$defs/") {
		name, _ := strings.CutPrefix(schema.Ref, "#/$defs/")
		schema.Ref = "#/components/schemas/" + name
	}
}

// OpenAPISpec returns the Request definition of a JSONRequest using "invopop/jsonschema"
func (r *JSONRequest[Body, Params]) OpenAPISpec() RequestSpec {
	bType := reflect.TypeOf(new(Body))
	for ; bType.Kind() == reflect.Pointer; bType = bType.Elem() {
	}

	schema := RequestSpec{}
	if bType != reflect.TypeOf(Nil{}) {
		s := (&jsonschema.Reflector{}).Reflect(new(Body))
		schema.RequestBody = &RequestBody{
			Content: map[string]MediaType{
				"application/json": {
					Schema: s,
				},
			},
			Required: reflect.TypeOf(*new(Body)).Kind() != reflect.Pointer,
		}
	}

	pType := reflect.TypeOf(new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	if pType != reflect.TypeOf(Nil{}) {
		schema.Parameters = CacheRequestParamsType(pType)
	}
	return schema
}

// JSONResponse[Body, Params any] is a response type that converts
// user-provided types to json and marshals params to headers
type JSONResponse[Body, Params any] struct {
	Body   Body
	Params Params
}

// WriteResponse writes the response and content-type header, it does not write the status code
func (r *JSONResponse[Body, Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	w.Header().Add("Content-Type", "application/json")
	if r == nil {
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	} else {
		h, err := MarshalParams(&r.Params)
		if err != nil {
			return err
		}
		for k, v := range h {
			for _, x := range v {
				w.Header().Add(k, x)
			}
		}
		w.WriteHeader(ctx.DefaultResponseCode())
		b, err := json.Marshal(&r.Body)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// OpenAPISpec returns the Responses definition of a JSONResponse using "invopop/jsonschema"
func (r *JSONResponse[Body, Params]) OpenAPISpec() Responses {
	schema := make(Responses)
	bType := reflect.TypeOf(new(Body))
	for ; bType.Kind() == reflect.Pointer; bType = bType.Elem() {
	}

	response := ResponseSpec{}
	if bType != reflect.TypeOf(Nil{}) {
		response.Content = map[string]MediaType{
			"application/json": {
				Schema: (&jsonschema.Reflector{
					ExpandedStruct: bType.Kind() == reflect.Struct,
					// DoNotReference: true,
					Lookup: func(t reflect.Type) jsonschema.ID {
						if t.PkgPath() == "" {
							return jsonschema.EmptyID
						}
						return jsonschema.ID("#/components/schemas/" + t.Name())
					},
				}).Reflect(new(Body)),
			},
		}
	}

	pType := reflect.TypeOf(*new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	if pType != reflect.TypeOf(Nil{}) {
		response.Headers = make(map[string]Parameter)
		for _, param := range CacheResponseParamsType(pType) {
			response.Headers[param.Name] = Parameter{
				Schema:          param.Schema,
				Description:     param.Description,
				Deprecated:      param.Deprecated,
				AllowReserved:   param.AllowReserved,
				AllowEmptyValue: param.AllowEmptyValue,
				Required:        param.Required,
				Explode:         param.Explode,
				Example:         param.Example,
				Examples:        param.Examples,
			}
		}
	}
	schema[""] = response
	return schema
}

// NewJSONResponse creates a JSONResponse from body and params
func NewJSONResponse[Body, Params any](body Body, params Params) *JSONResponse[Body, Params] {
	return &JSONResponse[Body, Params]{
		Body:   body,
		Params: params,
	}
}
