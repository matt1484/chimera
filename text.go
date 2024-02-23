package chimera

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/invopop/jsonschema"
)

var (
	_ RequestReader  = new(PlainTextRequest[Nil])
	_ ResponseWriter = new(PlainTextResponse[Nil])
	_ RequestReader  = new(PlainText[Nil])
	_ ResponseWriter = new(PlainText[Nil])
)

// PlainTextRequest is any text/plain request that results in a string body
type PlainTextRequest[Params any] struct {
	request *http.Request
	Body    string
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *PlainTextRequest[Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

func readPlainTextRequest[Params any](req *http.Request, ctx RouteContext, body *string, params *Params) error {
	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	*body = string(b)

	if _, ok := any(params).(*Nil); !ok {
		err = UnmarshalParams(req, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using io.ReadAll.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *PlainTextRequest[Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readPlainTextRequest(req, ctx, &r.Body, &r.Params)
}

func textRequestSpec[Params any](schema *RequestSpec) {
	schema.RequestBody = &RequestBody{
		Content: map[string]MediaType{
			"text/plain": {
				Schema: &jsonschema.Schema{
					Type: "string",
				},
			},
		},
	}

	pType := reflect.TypeOf(new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	if pType != reflect.TypeOf(Nil{}) {
		schema.Parameters = CacheRequestParamsType(pType)
	}
}

// OpenAPIRequestSpec describes the RequestSpec for text/plain requests
func (r *PlainTextRequest[Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	textRequestSpec[Params](&schema)
	return schema
}

// PlainTextRequest[Params] is any text/plain response that uses a string body
type PlainTextResponse[Params any] struct {
	Body   string
	Params Params
}

func writeTextResponse[Params any](w http.ResponseWriter, ctx RouteContext, body *string, params *Params) error {
	w.Header().Add("Content-Type", "text/plain")
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
	_, err = w.Write([]byte(*body))
	if err != nil {
		return err
	}
	return nil
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *PlainTextResponse[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeTextResponse[Params](w, ctx, &r.Body, &r.Params)
}

func textResponsesSpec[Params any](schema Responses) {
	response := ResponseSpec{}
	response.Content = map[string]MediaType{
		"text/plain": {
			Schema: &jsonschema.Schema{
				Type: "string",
			},
		},
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

// OpenAPIResponsesSpec describes the Responses for text/plain requests
func (r *PlainTextResponse[Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	textResponsesSpec[Params](schema)
	return schema
}

// NewPlainTextResponse creates a PlainTextResponse from a string and params
func NewPlainTextResponse[Params any](body string, params Params) *PlainTextResponse[Params] {
	return &PlainTextResponse[Params]{
		Body:   body,
		Params: params,
	}
}

// PlainText[Params] is a helper type that effectively works as both a PlainTextRequest[Params] and PlainTextResponse[Params]
// This is mostly here for convenience
type PlainText[Params any] struct {
	request *http.Request
	Body    string
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *PlainText[Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using io.ReadAll.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *PlainText[Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readPlainTextRequest(req, ctx, &r.Body, &r.Params)
}

// OpenAPIRequestSpec describes the RequestSpec for text/plain requests
func (r *PlainText[Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	textRequestSpec[Params](&schema)
	return schema
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *PlainText[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeTextResponse[Params](w, ctx, &r.Body, &r.Params)
}

// OpenAPIResponsesSpec describes the Responses for text/plain requests
func (r *PlainText[Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	textResponsesSpec[Params](schema)
	return schema
}
