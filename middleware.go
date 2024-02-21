package chimera

import "net/http"

// NextFunc is the format for allowing middleware to continue to child middlewares
type NextFunc func(req *http.Request) (ResponseWriter, error)

// MiddlewareFunc is a function that can be used as middleware
type MiddlewareFunc func(req *http.Request, next NextFunc) (ResponseWriter, error)

// TODO: add func to wrap and convert a http.Handler to a MiddlewareFunc
// TODO: allow middleware to read status codes correctly? the problem with lazy writing
// is that middleware cant get headers/body/response code easily
