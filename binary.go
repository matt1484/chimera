package chimera

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/invopop/jsonschema"
)

var (
	_ RequestReader  = new(BinaryRequest[Nil])
	_ ResponseWriter = new(BinaryResponse[Nil])
	_ RequestReader  = new(Binary[Nil])
	_ ResponseWriter = new(Binary[Nil])
)

// BinaryRequest[Params any] is a request type that uses a
// []byte as the Body and Params as an user-provided struct
type BinaryRequest[Params any] struct {
	request *http.Request
	Body    []byte
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *BinaryRequest[Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

func readBinaryRequest[Params any](req *http.Request, ctx RouteContext, body *[]byte, params *Params) error {
	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	*body = b

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
func (r *BinaryRequest[Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readBinaryRequest(req, ctx, &r.Body, &r.Params)
}

func binaryRequestSpec[Params any](schema *RequestSpec) {
	schema.RequestBody = &RequestBody{
		Content: map[string]MediaType{
			"application/octet-stream": {
				Schema: &jsonschema.Schema{
					Type:   "string",
					Format: "binary",
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

// OpenAPIRequestSpec returns the Request definition of a BinaryRequest
func (r *BinaryRequest[Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	binaryRequestSpec[Params](&schema)
	return schema
}

// BinaryResponse[Params any] is a response type that uses a
// []byte as the Body and Params as an user-provided struct
type BinaryResponse[Params any] struct {
	Body   []byte
	Params Params
}

func writeBinaryResponse[Params any](w http.ResponseWriter, ctx RouteContext, body *[]byte, params *Params) error {
	w.Header().Add("Content-Type", "application/octet-stream")
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
	_, err = w.Write(*body)
	if err != nil {
		return err
	}
	return nil
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *BinaryResponse[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeBinaryResponse(w, ctx, &r.Body, &r.Params)
}

func binaryResponsesSpec[Params any](schema Responses) {
	response := ResponseSpec{}
	response.Content = map[string]MediaType{
		"application/octet-stream": {
			Schema: &jsonschema.Schema{
				Type:   "string",
				Format: "binary",
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

// OpenAPIResponsesSpec returns the Responses definition of a BinaryResponse
func (r *BinaryResponse[Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	binaryResponsesSpec[Params](schema)
	return schema
}

// ResponseHead returns the status code and header for this response object
func (r *BinaryResponse[Params]) ResponseHead(ctx RouteContext) (ResponseHead, error) {
	head := ResponseHead{
		Header:     make(http.Header),
		StatusCode: ctx.DefaultResponseCode(),
	}
	h, err := MarshalParams(&r.Params)
	if err != nil {
		return head, err
	}
	for k, v := range h {
		for _, x := range v {
			head.Header.Add(k, x)
		}
	}
	return head, nil
}

// NewBinaryResponse creates a BinaryResponse from body and params
func NewBinaryResponse[Params any](body []byte, params Params) *BinaryResponse[Params] {
	return &BinaryResponse[Params]{
		Body:   body,
		Params: params,
	}
}

// Binary[Params] is a helper type that effectively works as both a BinaryRequest[Params] and BinaryResponse[Params]
// This is mostly here for convenience
type Binary[Params any] struct {
	request *http.Request
	Body    []byte
	Params  Params
}

// Context returns the context that was part of the original http.Request
// if this was used in a non-request context it will return nil
func (r *Binary[Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using io.ReadAll.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *Binary[Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	r.request = req
	return readBinaryRequest(req, ctx, &r.Body, &r.Params)
}

// OpenAPIRequestSpec returns the Request definition of a BinaryRequest
func (r *Binary[Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	binaryRequestSpec[Params](&schema)
	return schema
}

// WriteResponse writes the response body, parameters, and response code from context
func (r *Binary[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	if r == nil {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.WriteHeader(ctx.DefaultResponseCode())
		return nil
	}
	return writeBinaryResponse(w, ctx, &r.Body, &r.Params)
}

// OpenAPIResponsesSpec returns the Responses definition of a BinaryResponse
func (r *Binary[Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	binaryResponsesSpec[Params](schema)
	return schema
}

// ResponseHead returns the status code and header for this response object
func (r *Binary[Params]) ResponseHead(ctx RouteContext) (ResponseHead, error) {
	head := ResponseHead{
		Header:     make(http.Header),
		StatusCode: ctx.DefaultResponseCode(),
	}
	h, err := MarshalParams(&r.Params)
	if err != nil {
		return head, err
	}
	for k, v := range h {
		for _, x := range v {
			head.Header.Add(k, x)
		}
	}
	return head, nil
}
