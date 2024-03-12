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

An example of how some of this might look in practice is:

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
        }, nil
    })
    api.Start(":8000")
}

```
By using struct tags, generics, reflection, and carefully designed interfaces `chimera` can infer and automate a lot of the typical tasks that go developers often manually do.

This library relies heavily on `chi` to support the routing/grouping/middleware internally. While `chi` is all about being inline with the standard library, this library is more opinionated and instead goes for a more standard layout of development interfaces as seen in other languages/frameworks to support quicker development. A lof of its design was actually heavily inspired by `fastapi` and you can even see the parallels in this example here:

```python
from typing import Annotated
from fastapi import FastAPI, Header
from pydantic import BaseModel

class TestBody(BaseModel):
    prop: str

api = fastapi.FastAPI()

@api.middleware("http")
def add_process_time_header(request: Request, next):
    response = next(request)
    return response

@api.get("/test/{path}")
def test(path: str, header: Annotated[str, Header()] = None, body: TestBody) -> TestBody
    return body
```

## Docs
There are docs on github pages hosted [here](https://matt1484.github.io/chimera/)

## TODO
- Proper XML support
- Multipart form support