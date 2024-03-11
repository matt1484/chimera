---
title: Routing
layout: default
nav_order: 2
has_children: true
---

# Routing
Internally `chimera` uses [`go-chi/chi`](https://github.com/go-chi/chi) to setup and handle routing. This means that all pathing rules in `chi` apply here with some slight variations:
- Middleware and routes can be defined in any order with middleware always taking precedence
- Paths must be proper OpenAPI paths (i.e. `/path/{param}`) so regex paths arent fully supported yet
- Group routes can be defined indepently (instead of all in one function)

## APIs
`chimera` uses an `API` object to serve, group, and manage routes/middleware. Effectively a user must first create an `API` object like so:
```golang
api := chimera.NewAPI()
```
Afterwards ths `API` object is used to add handlers, middleware or sub-`API`s (via `Group()` and `Mount()`). 
Sub-`API`s allow routes to be isolated into groups and maintain separate middleware with the following rules:
- the parent `API` has its middleware evaluated first (including grandparents and so on)
- middleware in the sub-`API` wont be called for routes not directly attached to the sub-`API` object
- calls to `Group()` with the same base path will return the same sub-`API`
- calls to `Mount()` with the same base path will overwrite the existing sub-`API`
- if a base path doesn't match an immediate child `API` base path it will lead to route collision

An example of how to use sub-`API`s is:
```golang
api := chimera.NewAPI()

foo := chimera.NewAPI()
chimera.Get(foo, "/bar", func(req *chimera.Request) (*chimera.Response, error) {})
api.Mount("/foo", foo) // now /bar is actually /foo/bar

baz := foo.Group("/baz") // technically /foo/baz
// this route is technically /foo/baz/qux
chimera.Get(test, "/qux", func(req *chimera.Request) (*chimera.Response, error) {})
```

## Handlers
`API` handlers are affectively functions of the form 
```golang
func[ReqPtr RequestReaderPtr[Req], Req any, RespPtr ResponseWriterPtr[Resp], Resp any](ReqPtr) (RespPtr, error)
```
While this may seem confusing at first glance this is effectively just a generic function that:
- accepts a pointer to a type that implements `RequestReader`
- returns a pointer to a type that implements `ResponseWriter` (and error)

The simplest example being:
```golang
func(req *chimera.Request) (*chimera.Response, error) {}
```

## Routes
Handler functions are then bound to `HTTP` methods and paths using the following route functions:
- `Get(api *API, path string, handler HandlerFunc) Route`
- `Post(api *API, path string, handler HandlerFunc) Route`
- `Put(api *API, path string, handler HandlerFunc) Route`
- `Patch(api *API, path string, handler HandlerFunc) Route`
- `Delete(api *API, path string, handler HandlerFunc) Route`
- `Options(api *API, path string, handler HandlerFunc) Route`

An example of adding a route is:
```golang
api := chimera.NewAPI()
chimera.Post(api, "/test/{path}", func(req *chimera.Request) (*chimera.Response, error) {
    return &chimera.Response{
        Body: []byte("hello world"),
        StatusCode: 200,
    }, nil
})
```


All route functions in `chimera` return a `Route` object which contains its `OpenAPI` operation spec which is automatically created based on the parameters but can be edited/managed using the following helper methods:
- `OpenAPIOperationSpec() *Operation`: returns the raw `Operation` to be edited directly
- `WithResponseCode(code int) Route`: updates the default response code inline and returns the route
- `WithResponses(resp Responses) Route`: merges the `Responses` object into the existing one and returns the route
- `WithRequest(req RequestSpec) Route`: merges the `Request` object into the existing one and returns the route
- `WithOperation(op Operation) Route`: merges the `Operation` object into the existing one and returns the route
- `UsingResponses(resp Responses) Route`: replaces the existing `Responses` object with this one one and returns the route
- `UsingRequest(req RequestSpec) Route`: replaces the existing `Request` object with this one one and returns the route
- `UsingOperation(op Operation) Route`: replaces the existing `Operation`  object with this one one and returns the route
- `Internalize() Route`: marks the `Route` as "internal", meaning that it wont show up in the `OpenAPI` spec

Most of these functions were designed to be daisy-chained on route creation like so:
```golang
chimera.Post(api, "/test/{path}", func(req *chimera.Request) (*chimera.Response, error) {
}).WithResponseCode(
    418, // now this route will return a 418 by default. having these defaults is important for enforcing the spec
).WithOperation(chimera.Operation{ // now the docs will be more verbose
    Tags: []string{"tag1"},
    Summary: "this is a route",
})
```

## Standard lib support
`chimera` has a function that attempts to wrap/convert a generic standard library handler:
```golang
func HTTPHandler(handler http.HandlerFunc) chimera.HandlerFunc[*chimera.Request, chimera.Request, *chimera.Response, chimera.Response]
```
effectively the output of this function can be passed directly to a routing function but with the following caveats:
1. The OpenAPI spec for the route is very empty by default
2. The response body is passed to middleware in-memory so this should not be used for large response bodies
3. The response is still lazy so there are no errors ever from write which may be misleading