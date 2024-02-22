package chimera

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/invopop/jsonschema"
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
	return r.request.Context()
}

// ReadRequest reads the body of an http request and assigns it to the Body field using io.ReadAll.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *BinaryRequest[Params]) ReadRequest(req *http.Request, ctx RouteContext) error {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	r.Body = body

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

// OpenAPISpec returns the Request definition of a BinaryRequest
func (r *BinaryRequest[Params]) OpenAPISpec() RequestSpec {
	schema := RequestSpec{}
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
	return schema
}

// BinaryResponse[Params any] is a response type that uses a
// []byte as the Body and Params as an user-provided struct
type BinaryResponse[Params any] struct {
	Body   []byte
	Params Params
}

// WriteResponse writes the response and content-type header, it does not write the status code
func (r *BinaryResponse[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	w.Header().Add("Content-Type", "application/octet-stream")
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
		_, err = w.Write(r.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// OpenAPISpec returns the Responses definition of a BinaryResponse
func (r *BinaryResponse[Params]) OpenAPISpec() Responses {
	schema := make(Responses)
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
	return schema
}

// NewBinaryResponse creates a BinaryResponse from body and params
func NewBinaryResponse[Params any](body []byte, params Params) *BinaryResponse[Params] {
	return &BinaryResponse[Params]{
		Body:   body,
		Params: params,
	}
}
