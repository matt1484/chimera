# Middleware
`chimera` allows the use of generic middleware of the form:
```golang
func(req *http.Request, ctx chimera.RouteContext, next func(req *http.Request) (ResponseWriter, error)) (ResponseWriter, error)
```
Effectively this allows handlers to pass lazy-responses (i.e. not yet written) and errors to the middleware to allow easier error handling while still allowing requests and context to be modified before reaching the handler. An example of this would be:
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
    // or you could use ctx.GetResponseresp) to get the the whole Response in memory.
    // err could also be handled more gracefully here before making it to the error handler
    return resp, err
})
```
It's import to note that middleware can't write responses directly, but they can return a `ResponseWriter` (like `Response`) which will write it later. This is different from how a lot of other frameworks based on the standard lib handle it, but in theory makes it easy to construct responses on the fly. This does come with a few caveats:
1. Middleware can't tell if their response will have an error when writing
2. Standard library based middleware have to store the whole resp object in memory
