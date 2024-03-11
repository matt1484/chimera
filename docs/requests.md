---
title: Requests
layout: default
nav_order: 1
parent: Routing
---

# Requests
All route handlers use a pointer to a `RequestReader` (i.e. `RequestReaderPtr`) as their input which require the following methods to be defined:
```golang
// responsible for unmarshaling the request object
ReadRequest(*http.Request) error
// describes the request object (body, params, etc.)
OpenAPIRequestSpec() chimera.RequestSpec
```

All of the provided request types in `chimera` include sensible implementations for these methods so in general you only need to concern yourself with this if you are making your own `RequestReader` type.
The internal usage however is defined by:
1. `OpenAPIRequestSpec()` is run when the route is registered (via `Get`, `Post`, etc.) so any pre-valdation should happen there
2. `ReadRequest()` is run just before being passed to the handler (i.e. JSON parsing)
3. If any errors occur in `ReadRequest()` then the handler will NOT be executed

## Basic request types
`chimera` provides a few request types that implement `RequestReader` which are:
- `Request` which is just an alias for `http.Request`
- `NoBodyRequest[Params any]` which is a request with customizable params and no body (useful for GET requests)
- `EmptyRequest` which is a request that has no body or params (useful for GET requests)