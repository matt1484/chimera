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
	routes      []route
	middleware  []MiddlewareFunc
	subAPIs     []*API
	basePath    string
	parent      *API
	staticPaths map[string]string
}

// OpenAPISpec returns the underlying OpenAPI structure for this API
func (a *API) OpenAPISpec() *OpenAPI {
	return &a.openAPISpec
}

// ServeHTTP serves to implement support for the standard library
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	customWriter := httpResponseWriter{
		writer: w,
	}
	a.router.ServeHTTP(&customWriter, req)
	if customWriter.respError != nil {
		if err, ok := customWriter.respError.(APIError); ok {
			for k, vals := range err.Header {
				for _, v := range vals {
					customWriter.Header().Add(k, v)
				}
			}
			if err.StatusCode != 0 {
				customWriter.writer.WriteHeader(err.StatusCode)
			} else {
				customWriter.writer.WriteHeader(500)
			}
			customWriter.Write(err.Body)
		} else {
			customWriter.writer.WriteHeader(500)
			customWriter.Write(default500Error)
		}
	} else {
		// TODO: maybe allow global default response codes for methods?
		if customWriter.response != nil && !reflect.ValueOf(customWriter.response).IsNil() {
			customWriter.response.WriteResponse(&customWriter, customWriter.route.context)
		} else {
			if customWriter.route != nil && customWriter.route.context.responseCode != 0 {
				w.WriteHeader(customWriter.route.context.responseCode)
			} else {
				switch req.Method {
				case http.MethodGet:
				case http.MethodPut:
				case http.MethodPatch:
					w.WriteHeader(200)
				case http.MethodPost:
					w.WriteHeader(201)
				case http.MethodOptions:
				case http.MethodDelete:
					w.WriteHeader(204)
				}
			}
		}
	}
}

// Start uses http.ListenAndServe to start serving requests from addr
func (a *API) Start(addr string) {
	http.ListenAndServe(addr, a)
}

// NewAPI returns an initialized API object
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
func addRoute[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, method, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	reqSchema := ReqPtr(new(Req)).OpenAPIRequestSpec()
	operation := Operation{
		RequestSpec: &reqSchema,
		Responses:   RespPtr(new(Resp)).OpenAPIResponsesSpec(),
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
		} else if len(operation.Responses) == 0 {
			defaultCode = "200"
			responseCode = 200
		}
		pathSchema.Get = &operation
	case http.MethodPost:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["201"] = r
			defaultCode = "201"
			responseCode = 201
			delete(operation.Responses, "")
		} else if len(operation.Responses) == 0 {
			defaultCode = "201"
			responseCode = 201
		}
		pathSchema.Post = &operation
	case http.MethodDelete:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["204"] = r
			defaultCode = "204"
			responseCode = 204
			delete(operation.Responses, "")
		} else if len(operation.Responses) == 0 {
			defaultCode = "204"
			responseCode = 204
		}
		pathSchema.Delete = &operation
	case http.MethodOptions:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["204"] = r
			defaultCode = "204"
			responseCode = 204
			delete(operation.Responses, "")
		} else if len(operation.Responses) == 0 {
			defaultCode = "204"
			responseCode = 204
		}
		pathSchema.Options = &operation
	case http.MethodPatch:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["200"] = r
			defaultCode = "200"
			responseCode = 200
			delete(operation.Responses, "")
		} else if len(operation.Responses) == 0 {
			defaultCode = "200"
			responseCode = 200
		}
		pathSchema.Patch = &operation
	case http.MethodPut:
		if r, ok := operation.Responses[""]; ok {
			operation.Responses["200"] = r
			defaultCode = "200"
			responseCode = 200
			delete(operation.Responses, "")
		} else if len(operation.Responses) == 0 {
			defaultCode = "200"
			responseCode = 200
		}
		pathSchema.Put = &operation
	}
	api.openAPISpec.Paths[api.basePath+path] = pathSchema

	route := route{
		operationSpec: &operation,
		defaultCode:   defaultCode,
		context: &routeContext{
			responseCode: responseCode,
			method:       method,
			path:         path,
		},
	}
	chiHandler := (func(w http.ResponseWriter, r *http.Request) {
		request := ReqPtr(new(Req))
		customWriter := w.(*httpResponseWriter)
		customWriter.route = &route
		customWriter.respError = request.ReadRequest(r, route.context)
		if customWriter.respError != nil {
			return
		}
		customWriter.response, customWriter.respError = handler(request)
	})
	route.handler = chiHandler

	api.routes = append(api.routes, route)
	rebuildAPI(api)
	return &route
}

func rebuildAPI(api *API) {
	a := api
	for ; a.parent != nil; a = api.parent {
	}
	a.rebuildRouter()
}

// Get adds a "GET" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Get[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodGet, path, handler)
}

// Post adds a "POST" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Post[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodPost, path, handler)
}

// Put adds a "PUT" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Put[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodPut, path, handler)
}

// Patch adds a "PATCH" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Patch[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodPatch, path, handler)
}

// Delete adds a "DELETE" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Delete[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodDelete, path, handler)
}

// Options adds a "OPTIONS" route to the API object which will invode the handler function on route match
// it also returns the Route object to allow easy updates of the Operation spec
func Options[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](api *API, path string, handler HandlerFunc[ReqPtr, Req, RespPtr, Resp]) *route {
	return addRoute(api, http.MethodOptions, path, handler)
}

// Idk what trace even does, do people actually use this?
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
	rebuildAPI(a)
}

// Use adds middleware to the API
func (a *API) Use(middleware ...MiddlewareFunc) {
	a.middleware = append(a.middleware, middleware...)
	rebuildAPI(a)
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
	newSub.parent = a
	a.subAPIs = append(a.subAPIs, newSub)
	return newSub
}

// rebuildRouter rebuilds the entire router. This is not particularly efficient
// but at least this allows us to specify middleware/routes/groups in any order
// while still having a guaranteed final order
func (a *API) rebuildRouter() chi.Router {
	var schema []byte
	apiSpec := a.openAPISpec
	router := chi.NewRouter()

	if a.parent == nil {
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

	middlewareChain := make([]MiddlewareFunc, 0)
	for api := a; api != nil; api = api.parent {
		middlewareChain = append(api.middleware, middlewareChain...)
	}
	handler := func(w *httpResponseWriter, r *http.Request) (ResponseWriter, error) {
		w.route.handler(w, r)
		return w.response, w.respError
	}

	for i := len(middlewareChain) - 1; i >= 0; i-- {
		h := handler
		middleware := middlewareChain[i]
		handler = (func(w *httpResponseWriter, r *http.Request) (ResponseWriter, error) {
			wrapped := middlewareWrapper{
				writer:  w,
				handler: h,
			}
			return middleware(r, w.route.context, wrapped.Next)
		})
	}
	for _, route := range a.routes {
		if route.context.path == "" || route.context.path[0] != '/' {
			route.context.path = "/" + route.context.path
		}
		if len(middlewareChain) > 0 {
			router.MethodFunc(route.context.method, route.context.path, func(w http.ResponseWriter, r *http.Request) {
				writer := w.(*httpResponseWriter)
				writer.route = &route
				handler(writer, r)
			})
		} else {
			router.MethodFunc(route.context.method, route.context.path, route.handler)
		}

	}
	a.router = router
	return router
}
