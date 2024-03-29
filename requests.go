package chimera

import (
	"net/http"
	"reflect"
)

var (
	_ RequestReader = new(EmptyRequest)
	_ RequestReader = new(NoBodyRequest[Nil])
	_ RequestReader = new(Request)
)

// RequestReader is used to allow chimera to automatically read/parse requests
// as well as describe the parts of a request via openapi
type RequestReader interface {
	ReadRequest(*http.Request) error
	OpenAPIRequestSpec() RequestSpec
}

// RequestReaderPtr is just a workaround to allow chimera to accept a pointer
// to a RequestReader and convert to the underlying type
type RequestReaderPtr[T any] interface {
	RequestReader
	*T
}

// EmptyRequest is an empty request, effectively a no-op
// (mostly used for GET requests)
type EmptyRequest struct{}

// ReadRequest does nothing
func (*EmptyRequest) ReadRequest(*http.Request) error {
	return nil
}

// OpenAPIRequestSpec returns an empty RequestSpec
func (*EmptyRequest) OpenAPIRequestSpec() RequestSpec {
	return RequestSpec{}
}

// NoBodyRequest is a request with only parameters and an empty body
// (mostly used for GET requests)
type NoBodyRequest[Params any] struct {
	Params Params
}

// ReadRequest parses the params of the request
func (r *NoBodyRequest[Params]) ReadRequest(req *http.Request) error {
	r.Params = *new(Params)
	if _, ok := any(r.Params).(Nil); !ok {
		err := UnmarshalParams(req, &r.Params)
		if err != nil {
			return err
		}
	}
	return nil
}

// OpenAPIRequestSpec returns the parameter definitions of this object
func (r *NoBodyRequest[Params]) OpenAPIRequestSpec() RequestSpec {
	schema := RequestSpec{}
	pType := reflect.TypeOf(new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	if pType != reflect.TypeOf(Nil{}) {
		schema.Parameters = CacheRequestParamsType(pType)
	}
	return schema
}

// Request is just an http.Request that matches the expected interfaces
type Request http.Request

// ReadRequest assigns the request to the Request struct
func (r *Request) ReadRequest(req *http.Request) error {
	*r = Request(*req)
	return nil
}

// OpenAPIRequestSpec returns an empty RequestSpec
func (r *Request) OpenAPIRequestSpec() RequestSpec {
	return RequestSpec{}
}
