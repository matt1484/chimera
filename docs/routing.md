# Routing
Internally `chimera` uses [`go-chi/chi`](https://github.com/go-chi/chi) to setup and handle routing. This means that all pathing rules in `chi` apply here with some slight variations:
- Middleware and routes can be defined in any order with middleware always taking precedence
- Paths must be proper OpenAPI paths (i.e. `/path/{param}`) so regex paths arent fully supported yet
- Group routes can be defined indepently (instead of all in one function)

There are also some additional things of note:
- `Static()` binds a local directory for content serving (this is NOT shown in the OpenAPI spec)
- `/openapi.json` and `/docs*` are used for OpenAPI docs but can be overwritten
- Any errors returned during middleware/handlers are returned as generic 500s unless they are of type `APIError` which allows you to customize the response
- Writing of responses is lazy which means it wont happen until all middleware are ran

## Operation spec
All routes have an associated `Operation` spec (as per OpenAPI) which can be edited inline via a few helper functions:
- `WithResponseCode(code int)` which changes the default response code
- `WithResponses(resp Responses)` which can update the `Responses` object
- `WithRequest(req RequestSpec)` which can update the `Request` object
- `WithOperation(op Operation)` which can update the entire `Operation` object
- `Internalize()` which hides the route from the OpenAPI spec
- `UsingResponses(resp Responses)` which will replace the `Responses` object
- `UsingRequest(req RequestSpec)` which will replace the `Request` object
- `UsingOperation(op Operation)` which will replace the entire `Operation` object
An example of how to use these funciton would be something like:
```golang
chimera.Post(api, "/post", func (req *chimera.EmptyRequest) (*chimera.EmptyResponse, error) {
    ...
}).WithResponseCode(418).UsingOperation(chimera.Operation{})
```