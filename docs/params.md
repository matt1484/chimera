---
title: Parameters
layout: default
nav_order: 2
---
# Params
`chimera` attempts to marshal/unmarshal params via the standards set in [OpenAPI](https://spec.openapis.org/oas/v3.1.0#parameter-object).
This is done via 2 methods:
```golang
func UnmarshalParams(request *http.Request, obj any) error 
func MarshalParams(obj any) (http.Header, error) 
```
that each utilize the `param` struct tag of the format:
```golang
type ParamStructTag struct {
	Name            string `structtag:"$name"`
	In              In     `structtag:"in,required"`
	Explode         bool   `structtag:"explode"`
	Style           Style  `structtag:"style"`
	Required        bool   `structtag:"required"`
	Description     string `structtag:"description"`
	Deprecated      bool   `structtag:"deprecated"`
	AllowEmptyValue bool   `structtag:"allowEmptyValue"`
	AllowReserved   bool   `structtag:"allowReserved"`
}
```
The options closely follow the OpenAPI formats but an overview of the options is as follows:
- `$name`: the first value of the struct (similar to json)
- `in`: one of `cookie`, `path`, `query`, or `header` (same as OpenAPI)
- `explode`: controls how `style` works (same as OpenAPI), defaults to false
- `style`: one of `matrix`, `label`, `form`, `simple`, `spaceDelimited`, `pipeDelimited`, or `deepObject` (same as OpenAPI)
- `required`: marks the param as required (validation will fail if the param is missing)
- `description`: describes the param (same as OpenAPI)
- `deprecated`: marks the param as deprecated  (same as OpenAPI)
- `allowEmptyValue`: same as OpenAPI
- `allowReserved`: same as OpenAPI

A complete example of this is:
```golang
type Params struct {
    SomeProp string `param:"propName,in=query,explode,style=form,required,description='this is a param',allowEmptyValue,allowReserved"`
}
```
Each type that supports utilizing param structs would then unmarshal each field using the options provided.
Its important to note that fields that are `struct` types utilize the `prop` struct tag to determine the name of the sub properties of a param but cant provide any additional options for validation.

## Customizing
To support customization of param marshaling/unmarshaling the following functions can be implemented:
- `UnmarshalCookieParam(http.Cookie, ParamStructTag) error`
- `MarshalCookieParam(ParamStructTag) (http.Cookie, error)`
- `UnmarshalHeaderParam([]string, ParamStructTag) error`
- `MarshalHeaderParam(ParamStructTag) (http.Header, error)`
- `UnmarshalPathParam(string, ParamStructTag) error`
- `UnmarshalQueryParam(url.Values, ParamStructTag) error`

## Usage
The included request and response types utilize auto-handling of params like so:
```golang
type RequestParams struct {
	PathParam string `param:"some_param,description:'a parameter'"`
}

type ResponseParams struct {
	HeaderParam string `param:"my-header,description:'a response header'"`
}

chimera.Post(api, "/route/{some_param}", func(req *chimera.NoBodyRequest[RequestParams]) (*chimera.NoBodyResponse[ResponseParams], error) {
    // request contains parsed RequestParams
    return &chimera.NoBodyResponse[ResponseParams]{
		Params: ResponseParams{
			HeaderParam: "some header value", // this will get written to the header "my-header"
		},
	}, nil
})
```