---
title: Overview
layout: default
nav_order: 1
---

# chimera
Chi-based Module for Easy REST APIs

## Overview
`chimera` is designed for fast/easy API development based on OpenAPI with the following core features:
- Automatic OpenAPI (3.1) docs from structs (no comments or files needed)
- Automatic parsing of JSON/text/binary/form requests
- Automatic serialization of JSON/text/binary responses
- Automatic handling of request/response parameters (cookies/query/headers/path)
- Middleware (with easy error handling)
- Route groups (with isolated middleware)
- Error handling as responses
- Static file serving from directories

Gettings started is as easy as:

```go
package main

import "github.com/matt1484/chimera"

type TestBody struct {
    Property string `json:"prop"`
}

type TestParams struct {
    Path string   `param:"path,in=path"`
    Header string `param:"header,in=header"`
}

func main() {
    api := chimera.NewAPI()
    api.Use(func(req *http.Request, ctx chimera.RouteContext, next chimera.NextFunc) (chimera.ResponseWriter, error) {
        resp, err := next(req)
        return resp, err
    })
    chimera.Get(api, "/test/{path}", func(req *chimera.JSON[TestBody, TestParams]) (*chimera.JSON[TestBody, chimera.Nil], error) {
        return &chimera.JSON[TestBody, TestParams]{
            Body: req.Body,
        }
    })
    api.Start(":8000")
}

```
