package chimera

import (
	"fmt"
	"net/http"
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
}

// Route contains basic info about an API route and allows for inline editing of itself
type Route interface {
	// OpenAPISpec returns the Operation spec for this route
	OpenAPISpec() *Operation
	// WithResponseCode sets the default response code for this route
	// NOTE: the first time this is called, the presumption is that default code has been set based on http method
	WithResponseCode(code int) Route
	// WithResponses performs a merge on the operation's responses for this route
	WithResponses(resp Responses) Route
	// WithRequest performs a merge on the operation's request spec for this route
	WithRequest(req RequestSpec) Route
	// WithOperation performs a merge on the operation's spec for this route
	WithOperation(op Operation) Route
	// UsingResponses replaces the operation's responses for this route
	UsingResponses(resp Responses) Route
	// UsingRequest replaces the operation's request spec for this route
	UsingRequest(req RequestSpec) Route
	// UsingOperation replaces the operation's spec for this route
	UsingOperation(op Operation) Route
}

// OpenAPISpec returns the Operation spec for this route
func (r *route) OpenAPISpec() *Operation {
	return r.operationSpec
}

// WithResponseCode sets the default response code for this route
// NOTE: the first time this is called, the presumption is that default code has been set based on http method
func (r *route) WithResponseCode(code int) Route {
	r.context.responseCode = code
	if r.operationSpec == nil {
		return r
	}
	if _, ok := r.operationSpec.Responses[r.defaultCode]; ok {
		r.operationSpec.Responses[fmt.Sprint(code)] = r.operationSpec.Responses[r.defaultCode]
		delete(r.operationSpec.Responses, r.defaultCode)
		r.defaultCode = fmt.Sprint(code)
	}
	return r
}

// WithResponses performs a merge on the operation's responses for this route
func (r *route) WithResponses(resp Responses) Route {
	r.operationSpec.Responses.Merge(resp)
	return r
}

// WithRequest performs a merge on the operation's request spec for this route
func (r *route) WithRequest(req RequestSpec) Route {
	r.operationSpec.RequestSpec.Merge(req)
	return r
}

// WithOperation performs a merge on the operation's spec for this route
func (r *route) WithOperation(op Operation) Route {
	r.operationSpec.Merge(op)
	return r
}

// UsingResponses replaces the operation's responses for this route
func (r *route) UsingResponses(resp Responses) Route {
	r.operationSpec.Responses = resp
	return r
}

// UsingRequest replaces the operation's request spec for this route
func (r *route) UsingRequest(req RequestSpec) Route {
	r.operationSpec.RequestSpec = &req
	return r
}

// UsingOperation replaces the operation's spec for this route
func (r *route) UsingOperation(op Operation) Route {
	r.operationSpec = &op
	return r
}
