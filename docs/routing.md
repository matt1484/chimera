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