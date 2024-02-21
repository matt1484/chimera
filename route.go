package chimera

import (
	"fmt"
	"net/http"
)

// Route contains basic info about an API route
type Route struct {
	handler       http.HandlerFunc
	operationSpec *Operation
	path          string
	method        string
	defaultCode   string
	responseCode  int
}

// OpenAPISpec returns the Operation spec for this route
func (r *Route) OpenAPISpec() *Operation {
	return r.operationSpec
}

// WithResponseCode sets the default response code for this route
// NOTE: the first time this is called, the presumption is that default code has been set based on http method
func (r *Route) WithResponseCode(code int) *Route {
	r.responseCode = code
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
func (r *Route) WithResponses(resp Responses) *Route {
	r.operationSpec.Responses.Merge(resp)
	return r
}

// WithRequest performs a merge on the operation's request spec for this route
func (r *Route) WithRequest(req RequestSpec) *Route {
	r.operationSpec.RequestSpec.Merge(req)
	return r
}

// WithOperation performs a merge on the operation's spec for this route
func (r *Route) WithOperation(op Operation) *Route {
	r.operationSpec.Merge(op)
	return r
}

// UsingResponses replaces the operation's responses for this route
func (r *Route) UsingResponses(resp Responses) *Route {
	r.operationSpec.Responses = resp
	return r
}

// UsingRequest replaces the operation's request spec for this route
func (r *Route) UsingRequest(req RequestSpec) *Route {
	r.operationSpec.RequestSpec = &req
	return r
}

// UsingOperation replaces the operation's spec for this route
func (r *Route) UsingOperation(op Operation) *Route {
	r.operationSpec = &op
	return r
}
