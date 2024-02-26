# Requests
All route handlers use a pointer to a `RequestReader` (i.e. `RequestReaderPtr`) as their input which go through the following process before making it to a handler:
1. The `OpenAPIRequestSpec()` is run when the route is registered
2. `ReadRequest()` is run on it to handle any parsing (i.e. JSON)
3. The resulting value is sent to the handler

## Simple request types
`chimera` provides a few request types that implement `RequestReader` which are:
- `Request` which is just an `http.Request`
- `NoBodyRequest[Params]` which is a request that has no body
- `EmptyRequest` which is a request that has no body or params