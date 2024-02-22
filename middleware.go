package chimera

import "net/http"

// NextFunc is the format for allowing middleware to continue to child middleware
type NextFunc func(req *http.Request) (ResponseWriter, error)

// MiddlewareFunc is a function that can be used as middleware
type MiddlewareFunc func(req *http.Request, ctx RouteContext, next NextFunc) (ResponseWriter, error)

// TODO: add func to wrap and convert a http.Handler to a MiddlewareFunc

type middlewareWrapper struct {
	writer  *httpResponseWriter
	handler func(w *httpResponseWriter, r *http.Request) (ResponseWriter, error)
}

func (w *middlewareWrapper) Next(r *http.Request) (ResponseWriter, error) {
	return w.handler(w.writer, r)
}
