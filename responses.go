package chimera

import (
	"net/http"
	"reflect"
)

var (
	_ ResponseWriter      = new(EmptyResponse)
	_ ResponseWriter      = new(NoBodyResponse[Nil])
	_ ResponseWriter      = new(Response)
	_ http.ResponseWriter = new(Response)
	_ http.ResponseWriter = new(httpResponseWriter)
	_ ResponseWriter      = new(LazyBodyResponse)
)

// ResponseHead contains the head of an HTTP Response
type ResponseHead struct {
	StatusCode int
	Headers    http.Header
}

type BodyWriteFunc func(body []byte) (int, error)

type ResponseBodyWriter interface {
	WriteBody(write BodyWriteFunc) error
}

type ResponseHeadWriter interface {
	WriteHead(*ResponseHead) error
}

// ResponseWriter allows chimera to automatically write responses
type ResponseWriter interface {
	ResponseBodyWriter
	ResponseHeadWriter
	OpenAPIResponsesSpec() Responses
}

// ResponseWriterPtr is just a workaround to allow chimera to accept a pointer
// to a ResponseWriter and convert to the underlying type
type ResponseWriterPtr[T any] interface {
	ResponseWriter
	*T
}

// EmptyResponse is an empty response, effectively a no-op
// (mostly used for DELETE requests)
type EmptyResponse struct{}

// WriteHead does nothing
func (*EmptyResponse) WriteHead(*ResponseHead) error {
	return nil
}

// WriteBody does nothing
func (*EmptyResponse) WriteBody(BodyWriteFunc) error {
	return nil
}

// OpenAPIResponsesSpec returns an empty Responses definition
func (*EmptyResponse) OpenAPIResponsesSpec() Responses {
	return Responses{}
}

// NoBodyResponse is a response with no body, but has parameters
// (mostly used for DELETE requests)
type NoBodyResponse[Params any] struct {
	Params Params
}

// WriteBody does nothing
func (r *NoBodyResponse[Params]) WriteBody(BodyWriteFunc) error {
	return nil
}

// OpenAPIResponsesSpec returns the parameter definitions of this object
func (r *NoBodyResponse[Params]) OpenAPIResponsesSpec() Responses {
	schema := make(Responses)
	pType := reflect.TypeOf(*new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	response := ResponseSpec{}
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

// WriteHead writes the headers for this response
func (r *NoBodyResponse[Params]) WriteHead(head *ResponseHead) error {
	h, err := MarshalParams(&r.Params)
	if err != nil {
		return err
	}
	for k, v := range h {
		for _, x := range v {
			head.Headers.Add(k, x)
		}
	}
	return nil
}

// NewBinaryResponse creates a NoBodyResponse from params
func NewNoBodyResponse[Params any](params Params) *NoBodyResponse[Params] {
	return &NoBodyResponse[Params]{
		Params: params,
	}
}

// httpResponseWriter is the interal struct that overrides the default http.ResponseWriter
type httpResponseWriter struct {
	writer    http.ResponseWriter
	respError error
	response  ResponseWriter
	route     *route
	dirty     bool
}

// Header returns the response headers
func (w *httpResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes to the response body
func (w *httpResponseWriter) Write(b []byte) (int, error) {
	w.dirty = true
	return w.writer.Write(b)
}

// WriteHeader sets the status code
func (w *httpResponseWriter) WriteHeader(s int) {
	w.dirty = true
	w.writer.WriteHeader(s)
}

// Response is a simple response type to support creating responses on the fly
// it is mostly useful for middleware where execution needs to halt and an
// undefined response needs to be returned
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// WriteBody writes the exact body from the struct
func (r *Response) WriteBody(write BodyWriteFunc) error {
	_, err := write(r.Body)
	return err
}

// OpenAPIResponsesSpec returns an empty Responses object
func (r *Response) OpenAPIResponsesSpec() Responses {
	return Responses{}
}

// WriteHead returns the status code and header for this response object
func (r *Response) WriteHead(head *ResponseHead) error {
	if r.StatusCode > 0 {
		head.StatusCode = r.StatusCode
	}
	for k, v := range r.Headers {
		head.Headers[k] = v
	}
	return nil
}

// Write stores the body in the Reponse object for use later
func (r *Response) Write(body []byte) (int, error) {
	if r.Body == nil {
		r.Body = body
	} else {
		r.Body = append(r.Body, body...)
	}
	return len(body), nil
}

// WriteHeader stores the status code in the Reponse object for use later
func (r *Response) WriteHeader(status int) {
	r.StatusCode = status
}

// Header returns the current header for http.ResponseWriter compatibility
func (r *Response) Header() http.Header {
	return r.Headers
}

// NewResponse creates a response with the body, status, and header
func NewResponse(body []byte, statusCode int, header http.Header) *Response {
	return &Response{
		Body:       body,
		StatusCode: statusCode,
		Headers:    header,
	}
}

// LazyBodyResponse is a response that effectively wraps another ReponseWriter with predefined header/status code
type LazyBodyResponse struct {
	StatusCode int
	Body       ResponseBodyWriter
	Headers    http.Header
}

// WriteBody writes the exact body from the struct
func (r *LazyBodyResponse) WriteBody(write BodyWriteFunc) error {
	return r.Body.WriteBody(write)
}

// OpenAPIResponsesSpec returns an empty Responses object
func (r *LazyBodyResponse) OpenAPIResponsesSpec() Responses {
	return Responses{}
}

// WriteHead returns the status code and header for this response object
func (r *LazyBodyResponse) WriteHead(head *ResponseHead) error {
	if r.StatusCode > 0 {
		head.StatusCode = r.StatusCode
	}
	for k, v := range r.Headers {
		head.Headers[k] = v
	}
	return nil
}

// NewLazyBodyResponse creates a response with predefined headers and a lazy body
func NewLazyBodyResponse(head ResponseHead, resp ResponseBodyWriter) *LazyBodyResponse {
	return &LazyBodyResponse{
		Body:       resp,
		Headers:    head.Headers,
		StatusCode: head.StatusCode,
	}
}
