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

An example of how to use JSON in chimera is:
```golang
chimera.Post(api, "/route", func(req *chimera.JSON[Body, Params]) (*chimera.JSON[Body, Params], error) {
    // req contains an already parsed JSON body
    // any returned response will marshal as JSON before writing the body
    return nil, nil
})
```