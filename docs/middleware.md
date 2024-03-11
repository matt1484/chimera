---
title: Middleware
layout: default
nav_order: 3
parent: Routing
---

# Middleware
`chimera` allows the use of generic middleware of the form:
```golang
func(req *http.Request, ctx chimera.RouteContext, next func(req *http.Request) (ResponseWriter, error)) (ResponseWriter, error)
```
Effectively this allows for:
- requests to be modified before reaching handlers
- responses to be modified before written (which doesn't happen until after all middleware)
- errors in processing to be handled more gracefully
An example of this would be:
```golang
// Use adds middleware to an api
api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
    // ctx contains basic info about the matched route which cant be modified, but can be read
    // middleware can modify the request directly here before calling next
    // next invokes the next middleware or handler function
    resp, err := next(req)
    // resp is an interface technically, so it can't be read directly
    // but you could use ctx.GetResponseHead(resp) to get the headers/status to edit and then 
    // turn into a custom response with the same body like:
    // chimera.NewLazyBodyResponse(head, resp)
    // or you could use ctx.GetResponse(resp) to get the the whole Response in memory.
    // err could also be handled more gracefully here before the actual response gets written
    return resp, err
})
```
It's import to note that middleware can't write responses directly, but they can return a `ResponseWriter` (like `Response`) which will write it later. This is different from how a lot of other frameworks based on the standard lib handle it, but in theory makes it easy to construct responses on the fly. This does come with a few caveats:
1. Middleware can't tell if their response will have an error when writing
2. Standard library based middleware attempting to intercept responses will end up storing the whole response in memory

## Stand lib support
`chimera` has a wrapper function:
```golang
HTTPMiddleware(middleware func(http.Handler) http.Handler) chimera.MiddlewareFunc
```
that attempts to convert standard library based middleware into a lazy-response middleware that is used in `chimera`. In general this will work for most middleware that:
- modify requests
- write small responses
- edit response headers from next

While still being lazy-write like everything else in `chimera`. Any middleware using this that change the `http.ResponseWriter` before calling next will end up with the entire response written in memory. In general this is only an issue for large files or responses that would have historically been streamed to the writer.
