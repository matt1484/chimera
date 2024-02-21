package chimera

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/invopop/jsonschema"
	"github.com/swaggest/swgui/v5emb"
)

var (
	default500Error = []byte("Unknown error occurred")
)

// APIError is an error that can be converted to a response
type APIError struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Error returns the string representation of the error
func (a APIError) Error() string {
	if a.StatusCode < 1 {
		a.StatusCode = 500
	}
	return fmt.Sprintf("%v error: %s", a.StatusCode, a.Body)
}

// Nil is an empty struct that is designed to represent "nil"
// and is typically used to denote that a request/response
// has no body or parameters depending on context
type Nil struct{}

// API is a collection of routes and middleware with an associated OpenAPI spec
type API struct {
	openAPISpec OpenAPI
	router      *chi.Mux
	routes      []Route
	middlewares []MiddlewareFunc
	subAPIs     []*API
	basePath    string
	prime       *API
	staticPaths map[string]string
}

// OpenAPISpec returns the underlying OpenAPI structure for this API
func (api *API) OpenAPISpec() *OpenAPI {
	return &api.openAPISpec
}

// ServeHTTP serves to implement support for the standard library
func (api *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	customWriter := httpResponseWriter{
		writer: w,
	}
	api.router.ServeHTTP(&customWriter, req)
	if customWriter.respError != nil {
		if err, ok := customWriter.respError.(APIError); ok {
			if err.StatusCode != 0 {
				customWriter.writer.WriteHeader(err.StatusCode)
			} else {
				customWriter.writer.WriteHeader(500)
			}
			customWriter.Write(err.Body)
			for k, vals := range err.Header {
				for _, v := range vals {
					customWriter.Header().Add(k, v)
				}
			}
		} else {
			fmt.Println(customWriter.respError)
			customWriter.writer.WriteHeader(500)
			customWriter.Write(default500Error)
		}
	} else {
		// TODO: maybe allow global default response codes for methods?
		if customWriter.route != nil && customWriter.route.responseCode != 0 {
			customWriter.respCode = customWriter.route.responseCode
		} else {
			switch req.Method {
			case http.MethodGet:
			case http.MethodPut:
			case http.MethodPatch:
				customWriter.respCode = 200
			case http.MethodPost:
				customWriter.respCode = 201
			case http.MethodOptions:
			case http.MethodDelete:
				customWriter.respCode = 204
			}
		}
		if customWriter.response != nil && !reflect.ValueOf(customWriter.response).IsNil() {
			customWriter.response.WriteResponse(&customWriter)
		}
		if customWriter.respCode != 0 {
			customWriter.WriteHeader(customWriter.respCode)
		}
	}
}

// Start uses http.ListenAndServe to start serving requests from addr
func (api *API) Start(addr string) {
	http.ListenAndServe(addr, api)
}

// NewAPI return an initialized API object
func NewAPI() *API {
	api := API{
		router: chi.NewRouter(),
		openAPISpec: OpenAPI{
			OpenAPI: "3.1.0",
			Paths:   make(map[string]Path),
			Info: Info{
				Version: "v0.0.0",
				Title:   "API",
			},
			Servers: make([]Server, 0),
			Components: &Components{
				Schemas: make(map[string]jsonschema.Schema),
			},
		},
	}
	return &api
}

// addRoute creates a route based on method, path, handler, etc.
func addRoute[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, method, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	reqSchema := ReqPtr(new(Req)).OpenAPISpec()
	operation := Operation{
		RequestSpec: &reqSchema,
		Responses:   RespPtr(new(Resp)).OpenAPISpec(),
	}

	if reqSchema.RequestBody != nil {
		for k, v := range reqSchema.RequestBody.Content {
			if v.Schema != nil {
				// TODO: maybe implement a global schema resolver?
				// otherwise some classes with the same name may clobber
				// eachother in the final spec
				standardizedSchemas(v.Schema, api.openAPISpec.Components.Schemas)
			}
			reqSchema.RequestBody.Content[k] = v
		}
	}
	if operation.Responses != nil {
		for c, r := range operation.Responses {
			for k, v := range r.Content {
				if v.Schema != nil {
					standardizedSchemas(v.Schema, api.openAPISpec.Components.Schemas)
				}
				r.Content[k] = v
			}
			for k, v := range r.Headers {
				if v.Schema != nil {
					standardizedSchemas(v.Schema, api.openAPISpec.Components.Schemas)
				}
				r.Headers[k] = v
			}
			operation.Responses[c] = r
		}
	}
	if reqSchema.Parameters != nil {
		for i, p := range reqSchema.Parameters {
			if p.Schema != nil {
				standardizedSchemas(p.Schema, api.openAPISpec.Components.Schemas)
			}
			reqSchema.Parameters[i] = p
		}
	}
	pathSchema := Path{}
	if p, ok := api.openAPISpec.Paths[path]; ok {
		pathSchema = p
	}
	defaultCode := ""
	responseCode := 0
	switch method {
	case http.MethodGet:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["200"] = r
			defaultCode = "200"
			responseCode = 200
			delete(operation.Responses, "")
		}
		pathSchema.Get = &operation
	case http.MethodPost:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["201"] = r
			defaultCode = "201"
			responseCode = 201
			delete(operation.Responses, "")
		}
		pathSchema.Post = &operation
	case http.MethodDelete:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["204"] = r
			defaultCode = "204"
			responseCode = 204
			delete(operation.Responses, "")
		}
		pathSchema.Delete = &operation
	case http.MethodOptions:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["204"] = r
			defaultCode = "204"
			responseCode = 204
			delete(operation.Responses, "")
		}
		pathSchema.Options = &operation
	case http.MethodPatch:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["200"] = r
			defaultCode = "200"
			responseCode = 200
			delete(operation.Responses, "")
		}
		pathSchema.Patch = &operation
	case http.MethodPut:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["200"] = r
			defaultCode = "200"
			responseCode = 200
			delete(operation.Responses, "")
		}
		pathSchema.Put = &operation
	}
	api.openAPISpec.Paths[api.basePath+path] = pathSchema

	route := Route{
		method:        method,
		path:          path,
		operationSpec: &operation,
		defaultCode:   defaultCode,
		responseCode:  responseCode,
	}
	chiHandler := (func(w http.ResponseWriter, r *http.Request) {
		request := ReqPtr(new(Req))
		customWriter := w.(*httpResponseWriter)
		customWriter.route = &route
		customWriter.respError = request.ReadRequest(r)
		if customWriter.respError != nil {
			return
		}
		customWriter.response, customWriter.respError = handler(request)
	})
	route.handler = chiHandler

	api.routes = append(api.routes, route)
	if api.prime == nil {
		api.rebuildRouter()
	} else {
		api.prime.rebuildRouter()
	}
	return &route
}

// Get adds a "GET" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Get[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodGet, path, handler)
}

// Post adds a "POST" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Post[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodPost, path, handler)
}

// Put adds a "PUT" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Put[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodPut, path, handler)
}

// Patch adds a "PATCH" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Patch[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodPatch, path, handler)
}

// Delete adds a "DELETE" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Delete[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodDelete, path, handler)
}

// Options adds a "OPTIONS" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Options[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
	return addRoute(api, http.MethodOptions, path, handler)
}

// func Trace[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler func(ReqPtr) (RespPtr, error)) *Route {
// 	return addRoute(api, http.MethodTrace, path, handler)
// }

// Static adds support for serving static content from a directory, this route is hidden from the OpenAPI spec
func (a *API) Static(apiPath, filesPath string) {
	if len(apiPath) == 0 {
		apiPath = "/"
	}
	if apiPath[0] != '/' {
		apiPath = "/" + apiPath
	}
	if apiPath[len(apiPath)-1] != '/' {
		apiPath += "/"
	}
	a.staticPaths[apiPath+"*"] = filesPath
	if a.prime == nil {
		a.rebuildRouter()
	} else {
		a.prime.rebuildRouter()
	}
}

// Use adds middleware to the API
func (a *API) Use(middlewares ...MiddlewareFunc) {
	a.middlewares = append(a.middlewares, middlewares...)
	if a.prime == nil {
		a.rebuildRouter()
	} else {
		a.prime.rebuildRouter()
	}
}

// Group creates a sub-API with seperate middleware and routes using a base path.
// The middleware of the parent API is always evaluated first and any route collisions
// are handled by chi directly
func (a *API) Group(basePath string) *API {
	basePath = a.basePath + basePath
	for _, sub := range a.subAPIs {
		if sub.basePath == basePath {
			return a
		}
	}
	newSub := NewAPI()
	newSub.basePath = basePath
	if a.prime == nil {
		newSub.prime = a
	} else {
		newSub.prime = a.prime
	}
	a.subAPIs = append(a.subAPIs, newSub)
	return newSub
}

// rebuildRouter rebuilds the entire router. This is not particularly efficient
// but at least this allows us to specify middleware/routes/groups in any order
// while still having a guaranteed final order
func (a *API) rebuildRouter() chi.Router {
	apiSpec := a.openAPISpec
	router := chi.NewRouter()
	for _, middleware := range a.middlewares {
		router.Use(
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					customWriter := w.(*httpResponseWriter)
					customNext := (func(r *http.Request) (ResponseWriter, error) {
						next.ServeHTTP(w, r)
						return customWriter.response, customWriter.respError
					})
					customWriter.response, customWriter.respError = middleware(r, customNext)
				})
			},
		)
	}
	if a.prime == nil {
		var schema []byte
		router.MethodFunc(http.MethodGet, "/openapi.json", func(w http.ResponseWriter, r *http.Request) {
			if schema == nil {
				schema, _ = json.Marshal(a.openAPISpec)
			}
			w.Write(schema)
		})
		router.Handle("/docs*",
			v5emb.New(
				a.openAPISpec.Info.Title,
				"/openapi.json",
				"/docs",
			),
		)
	}
	for _, sub := range a.subAPIs {
		if sub.basePath == "" || sub.basePath[0] != '/' {
			sub.basePath = "/" + sub.basePath
		}
		router.Mount(sub.basePath, sub.rebuildRouter())
		apiSpec.Merge(sub.openAPISpec)
	}
	for apiPath, filesPath := range a.staticPaths {
		fileServer := http.FileServer(http.Dir(filesPath))
		router.Get(apiPath+"*", http.StripPrefix(apiPath, fileServer).ServeHTTP)
	}
	for _, route := range a.routes {
		if route.path == "" || route.path[0] != '/' {
			route.path = "/" + route.path
		}
		router.MethodFunc(route.method, route.path, route.handler)
	}
	a.router = router
	return router
}
