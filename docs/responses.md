# Responses
All route handlers use a pointer to a `ResponseWriter` (i.e. `ResponseWriterPtr`) as their output which go through the following process before making it to a handler:
1. The `OpenAPIRequestSpec()` is run when the route is registered
2. The value is passed through to middleware
3. `WriteResponse()` is run on it to write the response

Additionally, handlers can return an `error`, which if non-nil will ignore the response value and instead return a generic `500` or a custom response if it is a `APIError`.
Because response writing is lazy, any middleware can intercept and modify it before the response is actually returned, and writing of responses requires a `ResponseContext` since the intent is that `ResponseWriter` implementations are route-agnostic (i.e. are not bound to response code or route path). Middleware can call `ResponseHead()` on `ResponseWriter` objects to get the headers/status code without needing to write body

## Simple response types
`chimera` provides a few response types that implement `ResponseWriter` which are:
- `Response` which is just set of headers, body, response code
- `NoBodyResponse[Params]` which is a response that has no body
- `EmptyResponse` which is a response that has no body or params