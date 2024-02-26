package chimera

import (
	"fmt"
	"net/http"
)

var (
	_ RouteContext = new(routeContext)
)

// routeContext contains basic info about a matched Route
type routeContext struct {
	path         string
	method       string
	responseCode int
}

// RouteContext contains basic info about a matched Route
type RouteContext interface {
	// Path returns the path that the route was setup with (i.e. /route/{var})
	Path() string
	// Method returns the method intended to be used by the route
	Method() string
	// DefaultResponseCode returns the default response code for this route
	DefaultResponseCode() int
	// WithResponseCode replaces the default status code for this route
	WithResponseCode(int) RouteContext
}

// Path returns the path that the route was setup with (i.e. /route/{var})
func (r *routeContext) Path() string {
	return r.path
}

// Method returns the method intended to be used by the route
func (r *routeContext) Method() string {
	return r.method
}

// DefaultResponseCode returns the default response code for this route
func (r *routeContext) DefaultResponseCode() int {
	return r.responseCode
}

func (r *routeContext) WithResponseCode(code int) RouteContext {
	return &routeContext{
		method:       r.method,
		path:         r.path,
		responseCode: code,
	}
}

// route contains basic info about an API route
type route struct {
	handler       http.HandlerFunc
	operationSpec *Operation
	context       *routeContext
	defaultCode   string
	hidden        bool
	api           *API
}

// Route contains basic info about an API route and allows for inline editing of itself
type Route struct {
	route *route
}

// OpenAPIOperationSpec returns the Operation spec for this route
func (r Route) OpenAPIOperationSpec() *Operation {
	return r.route.operationSpec
}

// WithResponseCode sets the default response code for this route
// NOTE: the first time this is called, the presumption is that default code has been set based on http method
func (r Route) WithResponseCode(code int) Route {
	r.route.context.responseCode = code
	if r.route.operationSpec == nil {
		return r
	}
	if _, ok := r.route.operationSpec.Responses[r.route.defaultCode]; ok {
		r.route.operationSpec.Responses[fmt.Sprint(code)] = r.route.operationSpec.Responses[r.route.defaultCode]
		delete(r.route.operationSpec.Responses, r.route.defaultCode)
		r.route.defaultCode = fmt.Sprint(code)
	}
	return r
}

// WithResponses performs a merge on the operation's responses for this route
func (r Route) WithResponses(resp Responses) Route {
	r.route.operationSpec.Responses.Merge(resp)
	return r
}

// WithRequest performs a merge on the operation's request spec for this route
func (r Route) WithRequest(req RequestSpec) Route {
	r.route.operationSpec.RequestSpec.Merge(req)
	return r
}

// WithOperation performs a merge on the operation's spec for this route
func (r Route) WithOperation(op Operation) Route {
	r.route.operationSpec.Merge(op)
	return r
}

// UsingResponses replaces the operation's responses for this route
func (r Route) UsingResponses(resp Responses) Route {
	r.route.operationSpec.Responses = resp
	return r
}

// UsingRequest replaces the operation's request spec for this route
func (r Route) UsingRequest(req RequestSpec) Route {
	r.route.operationSpec.RequestSpec = &req
	return r
}

// UsingOperation replaces the operation's spec for this route
func (r Route) UsingOperation(op Operation) Route {
	r.route.operationSpec = &op
	return r
}

// Internalize hides the route from the api spec
func (r Route) Internalize() Route {
	r.route.hidden = true
	rebuildAPI(r.route.api)
	return r
}

// HandlerFunc is a handler function. The generic signature may look odd but its effectively:
// func(req *RequestReader) (*ResponseWriter, error)
type HandlerFunc[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any] func(ReqPtr) (RespPtr, error)

// HTTPHandler is a function that converts a standard http.HandlerFunc into one that works with chimera
func HTTPHandler(handler http.HandlerFunc) HandlerFunc[*Request, Request, *Response, Response] {
	return func(req *Request) (*Response, error) {
		response := Response{}
		handler(&response, (*http.Request)(req))
		return &response, nil
	}
}
