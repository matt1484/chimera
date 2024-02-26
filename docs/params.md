# Params
`chimera` attempts to marshal/unmarshal params via the standards set in [OpenAPI](https://spec.openapis.org/oas/v3.1.0#parameter-object).
This is done via 2 methods (`UnmarshalParams()` and `MarshalParams()`) that each utilize the `param` struct tags of the format:
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