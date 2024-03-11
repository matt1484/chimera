---
title: OpenAPI
layout: default
nav_order: 4
---
# OpenAPI
`chimera` has built in support for automatically generating OpenAPI 3.1 documentation. (the current latest version)
It does this by providing structs in [`openapi.go`](../openapi.go) that cover almost the entirety of the OpenAPI 3.1 [spec](https://spec.openapis.org/oas/v3.1.0)

When starting a server, the OpenAPI docs are available at `/openapi.json` with a Swagger UI at `/docs`.

## API Routes
The `API` struct contains top-level OpenAPI docs that can be retrieved and edited using `API.OpenAPISpec()`.
Technically sub-apis (a la `API.Group()`) also have an `OpenAPI` object but they get merged in to the main (i.e. top-most parent) `API` object. 
Each call to `Get`, `Post`, `Patch`, `Put`, `Delete`, `Options` adds an `Operation` to the `OpenAPI` object by using the provided path, associated HTTP method, and the `OpenAPIRequestSpec()` and `OpenAPIResponsesSpec()` functions on the handler's `RequestReader` and `ResponseWriter` respectively. 
The exception is routes created created with `Static()` which are hidden from the spec. Similarly you can run `Internalize()` on any route to hide it from the api spec.
Default status code is determined using the following logic:
- `Get`, `Patch`, `Put` assumes the default response code is `200` (unless this is explicitly changed on the `Route` object)
- `Post` assumes the default response code is `201` (unless this is explicitly changed on the `Route` object)
- `Delete`, `Options` assumes the default response code is `204` (unless this is explicitly changed on the `Route` object)
- The provided response types (i.e. `JSONResponse`) set the response code to "", so custom types that use that in their `Responses` object will have that code overwritten
- If path/method combinations are overwritten, the most recent one will take precendence in the spec
- Routes can have the default status code changed using `WithResponseCode()`

## JSONSchema
Since OpenAPI 3.1 supports JSONSchema, `chimera` uses [`invopop/jsonschema`](https://github.com/invopop/jsonschema) to generate schemas from request/response bodies. This relies heavily on the `jsonschema` struct tag and other relevant tags based on type (i.e. `json`, `form`, `param`, `prop`)

## Parameters
OpenAPI parameters use the `param` struct tag (`ParamStructTag`) to define how parameters are defined in the spec. The parameters section of the docs covers this further.