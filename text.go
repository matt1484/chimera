package chimera

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/invopop/jsonschema"
)

// PlainTextRequest is any text/plain request that results in a string body
type PlainTextRequest[Params any] struct {
	request *http.Request
	Body    string
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *PlainTextRequest[Params]) Context() context.Context {
	return r.request.Context()
}

// ReadRequest reads the body of an http request and assigns it to the Body field using io.ReadAll.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *PlainTextRequest[Params]) ReadRequest(req *http.Request) error {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	r.Body = string(body)

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

// OpenAPISpec describes the RequestSpec for text/plain requests
func (r *PlainTextRequest[Params]) OpenAPISpec() RequestSpec {
	schema := RequestSpec{}
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
	return schema
}

// PlainTextRequest is any text/plain response that uses a string body
type PlainTextResponse[Params any] struct {
	Body   string
	Params Params
}

// WriteResponse writes the response and content-type header, it does not write the status code
func (r *PlainTextResponse[Params]) WriteResponse(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "text/plain")
	if r == nil {
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
		_, err = w.Write([]byte(r.Body))
		if err != nil {
			return err
		}
	}
	return nil
}

// OpenAPISpec describes the Responses for text/plain requests
func (r *PlainTextResponse[Params]) OpenAPISpec() Responses {
	schema := make(Responses)
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
	return schema
}

// NewPlainTextResponse creates a PlainTextResponse from a string and params
func NewPlainTextResponse[Params any](body string, params Params) *PlainTextResponse[Params] {
	return &PlainTextResponse[Params]{
		Body:   body,
		Params: params,
	}
}
