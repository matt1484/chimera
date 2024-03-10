package chimera

import (
	"net/http"
)

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
type httpMiddlewareWriter struct {
	parts []ResponseBodyWriter
	statusCode int
	header http.Header
	savedHeader http.Header
}

type respBody []byte

func (r respBody) WriteBody(write BodyWriteFunc) error {
	_, err := write(r)
	return err
}

// Header returns the response headers
func (w *httpMiddlewareWriter) Header() http.Header {
	if w.statusCode > 0 {
		return w.savedHeader
	}
	return w.header
}

// Write writes to the response body
func (w *httpMiddlewareWriter) Write(b []byte) (int, error) {
	if w.statusCode < 1 {
		w.WriteHeader(200)
	}
	w.parts = append(w.parts, respBody(b))
	return len(b), nil
}

// WriteHeader sets the status code
func (w *httpMiddlewareWriter) WriteHeader(s int) {
	if w.statusCode > 0 {
		return
	}
	w.statusCode = s
	w.savedHeader = make(http.Header)
	for k, v := range w.header {
		w.savedHeader[k] = v
	}
}

// WriteHead returns the status code and header for this response object
func (w *httpMiddlewareWriter) WriteHead(head *ResponseHead) error {
	if w.statusCode > 0 {
		head.StatusCode = w.statusCode
	}
	for k, v := range w.Header() {
		head.Headers[k] = v
	}
	return nil
}

func (w *httpMiddlewareWriter) WriteBody(write BodyWriteFunc) error {
	var err error
	for _, p := range w.parts {
		err = p.WriteBody(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *httpMiddlewareWriter) OpenAPIResponsesSpec() Responses {
	return Responses{}
}


func HTTPMiddleware(middleware func(http.Handler) http.Handler) MiddlewareFunc {
	return func(req *http.Request, ctx RouteContext, next NextFunc) (ResponseWriter, error) {
		writer := httpMiddlewareWriter{
			header: make(http.Header),
		}
		nextWrapper := http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
			unchanged := w.(*httpMiddlewareWriter) == &writer
			resp, err := next(req)
			if err != nil {
				if unchanged && writer.statusCode < 1 {
					writeError(err, w)
				}
				return
			}
			if resp == nil {
				return
			}
			head := ResponseHead{
				Headers: w.Header(),
				StatusCode: ctx.DefaultResponseCode(),
			}
			err = resp.WriteHead(&head)
			if err != nil {
				if unchanged && writer.statusCode < 1 {
					writeError(err, w)
				}
				return
			}
			if head.StatusCode > 0 {
				w.WriteHeader(head.StatusCode)
			} else {
				w.WriteHeader(ctx.DefaultResponseCode())
			}
			if resp == nil {
				return
			}
			if unchanged {
				writer.parts = append(writer.parts, resp)
			} else {
				resp.WriteBody(w.Write)
			}
		})
		middleware(nextWrapper).ServeHTTP(&writer, req)
		return &writer, nil
	}
}