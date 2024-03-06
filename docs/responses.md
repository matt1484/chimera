# Responses
All route handlers use a pointer to a `ResponseWriter` (i.e. `ResponseWriterPtr`) as their output which go through the following process before responding to a request:
1. The `OpenAPIRequestSpec()` is run when the route is registered
2. The response value is passed through middleware (via `next()`)
3. After all middleware is done `WriteHead()` and `WriteResponse()` are used to write the response body/head to the underlying writer

Additionally, handlers can return an `error`, which if non-nil will ignore the response value and instead return a generic `500` or a custom response if it is a `APIError`.
Because response writing is lazy, any middleware can intercept and modify it before the response is actually returned and without having to use callbacks or custom `http.ResponseWriter`s to intercept the response

## Simple response types
`chimera` provides a few response types that implement `ResponseWriter` which are:
- `Response` which is just set of headers, body, response code
- `NoBodyResponse[Params]` which is a response that has no body
- `EmptyResponse` which is a response that has no body or params
- `LazybodyResponse` which is a response with predefined headers/status code and a lazy body (written after middleware)
