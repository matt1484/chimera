package chimera

import "net/http"

// NextFunc is the format for allowing middleware to continue to child middleware
type NextFunc func(req *http.Request) (ResponseWriter, error)

// MiddlewareFunc is a function that can be used as middleware
type MiddlewareFunc func(req *http.Request, ctx RouteContext, next NextFunc) (ResponseWriter, error)

type middlewareWrapper struct {
	writer  *httpResponseWriter
	handler func(w *httpResponseWriter, r *http.Request) (ResponseWriter, error)
}

func (w *middlewareWrapper) Next(r *http.Request) (ResponseWriter, error) {
	return w.handler(w.writer, r)
}

// TODO: add func to wrap and convert a http.Handler to a MiddlewareFunc
type httpMiddlewareWrapper struct {
	writer http.ResponseWriter
	dirty  bool
}

// Header returns the response headers
func (w *httpMiddlewareWrapper) Header() http.Header {
	return w.writer.Header()
}

// Write writes to the response body
func (w *httpMiddlewareWrapper) Write(b []byte) (int, error) {
	w.dirty = true
	return w.writer.Write(b)
}

// WriteHeader sets the status code
func (w *httpMiddlewareWrapper) WriteHeader(s int) {
	w.writer.WriteHeader(s)
}
