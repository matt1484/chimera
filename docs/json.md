---
title: JSON
layout: default
nav_order: 3
---

# JSON
`chimera` supports requests and responses with JSON bodies via these classes:
- `JSONRequest[Body, Params any]`: JSON request with `Body` type being parsed via `encoding/json` and `Params` type being parsed via `chimera.UnmarshalParams`
- `JSONResponse[Body, Params any]`: JSON response with `Body` type being marshaled via `encoding/json` and `Params` type being marshaled via `chimera.MarshalParams`
- `JSON[Body, Params any]`: represents both `JSONRequest[Body, Params any]` and `JSONResponse[Body, Params any]`

Essentially, you just need to provide any valid JSON type for `Body` and any struct with `param` tags for `Params` and the object is structured like so:
```golang
type JSON[Body, Params any] struct {
	Body    Body
	Params  Params
}
```
so the parsed body ends up in `Body` and parsed params end up in `Params`.

## Usage
An example of how to use JSON in chimera is:
```golang
type RequestBody struct {
    S string `json:"s" jsonschema:"description='some property'"`
}

type RequestParams struct {
    P string `param:"p,in=path"
}

type ResponseBody struct {
    I int `json:"i"`
}

type ResponseParams struct {
    H string `param:"h,in=header" 
}

chimera.Post(api, "/route/{path}", func(req *chimera.JSON[RequestBody, RequestParams]) (*chimera.JSON[ResponseBody, ResponseParams], error) {
    // req contains an already parsed JSON body
    // any returned response will be marshaled as JSON before writing the body
    return nil, nil
})
```