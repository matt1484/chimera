---
title: Responses
layout: default
nav_order: 2
parent: Routing
---

# Responses
All route handlers use a pointer to a `ResponseWriter` (i.e. `ResponseWriterPtr`) as part of their output (plus an error) which require the following methods to be defined:
```golang
WriteBody(write func(body []byte) (int, error)) error
WriteHead(head *chimera.ResponseHead) error
OpenAPIResponsesSpec() chimera.Responses
```

All of the provided request types in `chimera` include sensible implementations for these methods so in general you only need to concern yourself with this if you are making your own `ResponseWriter` type.
The internal usage however is defined by:
1. The `OpenAPIRequestSpec()` is run when the route is registered (via `Get`, `Post`, etc.) so any pre-valdation should happen there
2. When handlers return a `ResponseWriter`, the write does NOT happen immediately and can be intercepted by middleware
3. After all middleware is done `WriteHead()` and `WriteResponse()` are used to write the response body/head to the underlying `http.ResponseWriter`
4. `WriteHead()` recieves a `ResponseHead` object with the default status code already set and an empty `http.Header` map
5. `WriteResponse()` and `WriteHead()` should be kept very simple since errors returned by them can not be easily caught. In general the only issues that should occur here are errors involving serialization (i.e. `json.Marshal`) or writing (i.e. `http.ResponseWriter.Write`)
6. handlers can return an `error`, which if non-nil will ignore the response value and instead return a generic `500` or a custom response if it is a `chimera.APIError`.


## Simple response types
`chimera` provides a few response types that implement `ResponseWriter` which are:
- `Response` which is just a set of predefined headers, body, response code
- `NoBodyResponse[Params any]` which is a response that has no body but returns headers
- `EmptyResponse` which is a response that has no body or params
- `LazybodyResponse` which is a response with predefined headers/status code and a lazy body (written after middleware)
