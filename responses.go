package chimera

import (
	"net/http"
	"reflect"
)

// ResponseWriter allows chimera to automatically write responses
type ResponseWriter interface {
	WriteResponse(http.ResponseWriter) error
	OpenAPISpecifier[Responses]
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

// WriteResponse does nothing
func (*EmptyResponse) WriteResponse(w http.ResponseWriter) error {
	return nil
}

// OpenAPISpec returns an empty Responses definition
func (*EmptyResponse) OpenAPISpec() Responses {
	return Responses{}
}

// EmptyResponse is a response with no body, but has parameters
// (mostly used for DELETE requests)
type NoBodyResponse[Params any] struct {
	Params Params
}

// WriteResponse writes the reponse headers (but not status code)
func (r *NoBodyResponse[Params]) WriteResponse(w http.ResponseWriter) error {
	h, err := MarshalParams(&r.Params)
	if err != nil {
		return err
	}
	for k, v := range h {
		for _, x := range v {
			w.Header().Add(k, x)
		}
	}
	return nil
}

// OpenAPISpec returns the parameter definitions of this object
func (r *NoBodyResponse[Params]) OpenAPISpec() Responses {
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
	route     *Route
	respCode  int
}

// Header returns the response headers
func (w *httpResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes to the response body
func (w *httpResponseWriter) Write(b []byte) (int, error) {
	if w.respCode != 0 {
		w.WriteHeader(w.respCode)
	}
	w.respCode = 0
	return w.writer.Write(b)
}

// WriteHeader sets the status code
func (w *httpResponseWriter) WriteHeader(s int) {
	w.respCode = 0
	w.writer.WriteHeader(s)
}

// Response is a simple response type to support creating responses on the fly
// it is mostly useful for middleware where execution needs to halt and an
// undefined response needs to be returned
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// WriteResponse writes the exact body, headers, and status code
func (r *Response) WriteResponse(w http.ResponseWriter) error {
	for k, v := range r.Header {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}
	if r.StatusCode > 0 {
		w.WriteHeader(r.StatusCode)
	}
	_, err := w.Write(r.Body)
	return err
}

// OpenAPISpec returns an empty Responses object
func (r *Response) OpenAPISpec() Responses {
	return Responses{}
}

func NewResponse(body []byte, statusCode int, header http.Header) *Response {
	return &Response{
		Body:       body,
		StatusCode: statusCode,
		Header:     header,
	}
}
