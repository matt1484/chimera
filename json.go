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
	_ RequestReader  = new(JSON[Nil, Nil])
	_ ResponseWriter = new(JSON[Nil, Nil])
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
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

func readJSONRequest[Body, Params any](req *http.Request, ctx RouteContext, body *Body, params *Params) error {
	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if _, ok := any(body).(*Nil); !ok {
		err = json.Unmarshal(b, body)
		if err != nil {
			return err
		}
	}

	if _, ok := any(params).(*Nil); !ok {
		err = UnmarshalParams(req, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using json.Unmarshal
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *JSONRequest[Body, Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readJSONRequest(req, ctx, &r.Body, &r.Params)
}

// standardizedSchemas basically tries to convert all jsonschema.Schema objects to be mapped
// in defs and then replace them with $refs
func standardizedSchemas(schema *jsonschema.Schema, defs map[string]jsonschema.Schema) {
	if schema == nil {
		return
	}
	schema.ID = ""
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
	schema.ID = ""
}

func jsonRequestSpec[Body, Params any](schema *RequestSpec) {
	bType := reflect.TypeOf(new(Body))
	for ; bType.Kind() == reflect.Pointer; bType = bType.Elem() {
	}

	if bType != reflect.TypeOf(Nil{}) {
		s := (&jsonschema.Reflector{
			// ExpandedStruct: bType.Kind() == reflect.Struct,
		}).Reflect(new(Body))
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
}

func jsonResponsesSpec[Body, Params any](schema Responses) {
	bType := reflect.TypeOf(new(Body))
	for ; bType.Kind() == reflect.Pointer; bType = bType.Elem() {
	}

	response := ResponseSpec{}
	if bType != reflect.TypeOf(Nil{}) {
		response.Content = map[string]MediaType{
			"application/json": {
				Schema: (&jsonschema.Reflector{
					// ExpandedStruct: bType.Kind() == reflect.Struct,
					// DoNotReference: true,
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
}

// OpenAPIRequestSpec returns the Request definition of a JSONRequest using "invopop/jsonschema"
func (r *JSONRequest[Body, Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	jsonRequestSpec[Body, Params](&schema)
	return schema
}

// JSONResponse[Body, Params any] is a response type that converts
// user-provided types to json and marshals params to headers
type JSONResponse[Body, Params any] struct {
	Body   Body
	Params Params
}

func writeJSONResponse[Body, Params any](w http.ResponseWriter, ctx RouteContext, body *Body, params *Params) error {
	w.Header().Add("Content-Type", "application/json")

	h, err := MarshalParams(params)
	if err != nil {
		return err
	}
	for k, v := range h {
		for _, x := range v {
			w.Header().Add(k, x)
		}
	}
	w.WriteHeader(ctx.DefaultResponseCode())
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *JSONResponse[Body, Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeJSONResponse(w, ctx, &r.Body, &r.Params)
}

// OpenAPIResponsesSpec returns the Responses definition of a JSONResponse using "invopop/jsonschema"
func (r *JSONResponse[Body, Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	jsonResponsesSpec[Body, Params](schema)
	return schema
}

// NewJSONResponse creates a JSONResponse from body and params
func NewJSONResponse[Body, Params any](body Body, params Params) *JSONResponse[Body, Params] {
	return &JSONResponse[Body, Params]{
		Body:   body,
		Params: params,
	}
}

// JSON[Body, Params] is a helper type that effectively works as both a JSONRequest[Body, Params] and JSONResponse[Body, Params]
// This is mostly here for convenience
type JSON[Body, Params any] struct {
	request *http.Request
	Body    Body
	Params  Params
}

// Context returns the context for this request
// NOTE: this type can also be used for responses in which case Context() would be nil
func (r *JSON[Body, Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using json.Unmarshal
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *JSON[Body, Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readJSONRequest(req, ctx, &r.Body, &r.Params)
}

// OpenAPIRequestSpec returns the Request definition of a JSON request using "invopop/jsonschema"
func (r *JSON[Body, Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	jsonRequestSpec[Body, Params](&schema)
	return schema
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *JSON[Body, Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeJSONResponse(w, ctx, &r.Body, &r.Params)
}

// OpenAPIResponsesSpec returns the Responses definition of a JSON response using "invopop/jsonschema"
func (r *JSON[Body, Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	jsonResponsesSpec[Body, Params](schema)
	return schema
}
