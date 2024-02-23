package chimera

import (
	"net/http"
	"reflect"
)

var (
	_ ResponseWriter = new(EmptyResponse)
	_ ResponseWriter = new(NoBodyResponse[Nil])
	_ ResponseWriter = new(Response)
)

// ResponseWriter allows chimera to automatically write responses
type ResponseWriter interface {
	WriteResponse(http.ResponseWriter, RouteContext) error
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

// WriteResponse just writes the default status code and no body/headers
func (*EmptyResponse) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	w.WriteHeader(ctx.DefaultResponseCode())
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

// WriteResponse writes the response headers and response code from context
func (r *NoBodyResponse[Params]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
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
}

// Header returns the response headers
func (w *httpResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes to the response body
func (w *httpResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// WriteHeader sets the status code
func (w *httpResponseWriter) WriteHeader(s int) {
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

// WriteResponse writes the exact body, headers, and status code from the struct
func (r *Response) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	for k, v := range r.Headers {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}
	if r.StatusCode > 0 {
		w.WriteHeader(r.StatusCode)
	} else {
		w.WriteHeader(ctx.DefaultResponseCode())
	}
	_, err := w.Write(r.Body)
	return err
}

// OpenAPIResponsesSpec returns an empty Responses object
func (r *Response) OpenAPIResponsesSpec() Responses {
	return Responses{}
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

// RecordResponse copies the response from a ResponseWriter using its WriteResponse method and context
func RecordResponse(w ResponseWriter, ctx RouteContext) (*Response, error) {
	resp := Response{}
	err := w.WriteResponse(&resp, ctx)
	return &resp, err
}
