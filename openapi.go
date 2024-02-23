package chimera

import (
	"github.com/invopop/jsonschema"
)

// OpenAPI is used to store an entire openapi spec
type OpenAPI struct {
	OpenAPI           string                `json:"openapi,omitempty"`
	Info              Info                  `json:"info,omitempty"`
	JSONSchemaDialect string                `json:"jsonSchemaDialect,omitempty"`
	Servers           []Server              `json:"servers,omitempty"`
	Paths             map[string]Path       `json:"paths,omitempty"`
	Webhooks          map[string]Path       `json:"webhooks,omitempty"`
	Components        *Components           `json:"components,omitempty"`
	Security          []map[string][]string `json:"security,omitempty"`
	Tags              []Tag                 `json:"tags,omitempty"`
	ExternalDocs      *ExternalDocs         `json:"externalDocs,omitempty"`
}

// Merge attempts to fold a spec into the current spec
// In general, the provided spec takes precedence over the current one
func (o *OpenAPI) Merge(spec OpenAPI) {
	o.Servers = append(spec.Servers, o.Servers...)
	for path, obj := range spec.Paths {
		pathObj, ok := spec.Paths[path]
		if !ok {
			pathObj = obj
		} else {
			// only replace if not found seems appropriate
			// but maybe later add a merge function?
			if pathObj.Delete == nil {
				pathObj.Delete = obj.Delete
			}
			if pathObj.Get == nil {
				pathObj.Get = obj.Get
			}
			if pathObj.Patch == nil {
				pathObj.Patch = obj.Patch
			}
			if pathObj.Post == nil {
				pathObj.Post = obj.Post
			}
			if pathObj.Put == nil {
				pathObj.Put = obj.Put
			}
			if pathObj.Options == nil {
				pathObj.Options = obj.Options
			}
			if pathObj.Trace == nil {
				pathObj.Trace = obj.Trace
			}
			// TODO: handle dupes?
			pathObj.Parameters = append(pathObj.Parameters, obj.Parameters...)
			if pathObj.Description == "" {
				pathObj.Description = obj.Description
			} else {
				pathObj.Description += "\n" + obj.Description
			}
			pathObj.Servers = append(spec.Servers, o.Servers...)
			if pathObj.Summary == "" {
				pathObj.Summary = obj.Summary
			} else {
				pathObj.Summary += "\n" + obj.Summary
			}
		}
		o.Paths[path] = pathObj
	}
	if o.Components == nil {
		o.Components = spec.Components
	} else {
		for k, v := range spec.Components.Parameters {
			if _, ok := o.Components.Parameters[k]; !ok {
				o.Components.Parameters[k] = v
			}
		}
		for k, v := range spec.Components.Callbacks {
			if _, ok := o.Components.Callbacks[k]; !ok {
				o.Components.Callbacks[k] = v
			}
		}
		for k, v := range spec.Components.Headers {
			if _, ok := o.Components.Headers[k]; !ok {
				o.Components.Headers[k] = v
			}
		}
		for k, v := range spec.Components.Links {
			if _, ok := o.Components.Links[k]; !ok {
				o.Components.Links[k] = v
			}
		}
		for k, v := range spec.Components.PathItems {
			if _, ok := o.Components.PathItems[k]; !ok {
				o.Components.PathItems[k] = v
			}
		}
		for k, v := range spec.Components.RequestBodies {
			if _, ok := o.Components.RequestBodies[k]; !ok {
				o.Components.RequestBodies[k] = v
			}
		}
		for k, v := range spec.Components.Responses {
			if _, ok := o.Components.Responses[k]; !ok {
				o.Components.Responses[k] = v
			}
		}
		for k, v := range spec.Components.Schemas {
			if _, ok := o.Components.Schemas[k]; !ok {
				o.Components.Schemas[k] = v
			}
		}
		for k, v := range spec.Components.SecuritySchemes {
			if _, ok := o.Components.SecuritySchemes[k]; !ok {
				o.Components.SecuritySchemes[k] = v
			}
		}
	}
	if o.ExternalDocs == nil {
		o.ExternalDocs = spec.ExternalDocs
	}
	// TODO: maybe actually merge the sub dicts?
	o.Security = append(spec.Security, o.Security...)
	// TODO: dedupe
	o.Tags = append(spec.Tags, o.Tags...)
	if len(o.Webhooks) == 0 {
		o.Webhooks = spec.Webhooks
	} else {
		for k, v := range spec.Webhooks {
			if _, ok := o.Webhooks[k]; !ok {
				o.Webhooks[k] = v
			}
		}
	}
}

// Info holds info about an API
type Info struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

// License describes the license of an API
type License struct {
	Name       string `json:"name,omitempty"`
	URL        string `json:"url,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// Contact stores basic contanct info
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// Server describes a server
type Server struct {
	URL         string                    `json:"url,omitempty"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable is a variable used in servers
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Path stores all operations allowed on a particular path
type Path struct {
	// Ref         string      `json:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty"`
	Trace       *Operation  `json:"trace,omitempty"`
	Servers     []Server    `json:"servers,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// RequestSpec is the description of an openapi request used in an Operation
type RequestSpec struct {
	Parameters  []Parameter  `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"requestBody,omitempty"`
}

// Merge folds a Request object into the current Request object
// In general, the provided request takes precedence over the current one
func (r *RequestSpec) Merge(other RequestSpec) {
	// TODO ensure the parameters can be otherwritten?
	r.Parameters = append(r.Parameters, other.Parameters...)
	if r.RequestBody == nil {
		r.RequestBody = other.RequestBody
	} else if other.RequestBody != nil {
		if len(other.RequestBody.Description) != 0 {
			r.RequestBody.Description = other.RequestBody.Description
		}
		r.RequestBody.Required = other.RequestBody.Required || r.RequestBody.Required
		if r.RequestBody.Content == nil {
			r.RequestBody.Content = other.RequestBody.Content
		} else {
			for k, v := range other.RequestBody.Content {
				r.RequestBody.Content[k] = v
			}
		}
	}
}

// Operation describes an openapi Operation
type Operation struct {
	*RequestSpec
	Tags         []string                   `json:"tags,omitempty"`
	Summary      string                     `json:"summary,omitempty"`
	Description  string                     `json:"description,omitempty"`
	ExternalDocs *ExternalDocs              `json:"externalDocs,omitempty"`
	OperationID  string                     `json:"operationId,omitempty"`
	Callbacks    map[string]map[string]Path `json:"callbacks,omitempty"`
	Deprecated   bool                       `json:"deprecated,omitempty"`
	Security     []map[string][]string      `json:"security,omitempty"`
	Servers      []Server                   `json:"servers,omitempty"`
	Responses    Responses                  `json:"responses,omitempty"`
}

// Merge folds an Operation object into the current Operation object
// In general, the provided operation takes precedence over the current one
func (o *Operation) Merge(other Operation) {
	if other.RequestSpec != nil {
		if o.RequestSpec != nil {
			o.RequestSpec.Merge(*other.RequestSpec)
		} else {
			o.RequestSpec = other.RequestSpec
		}
	}
	if o.Tags == nil || len(o.Tags) == 0 {
		o.Tags = other.Tags
	} else if len(other.Tags) != 0 {
		tagMap := make(map[string]struct{})
		for _, tag := range append(o.Tags, other.Tags...) {
			tagMap[tag] = struct{}{}
		}
		o.Tags = make([]string, len(tagMap))
		i := 0
		for t := range tagMap {
			o.Tags[i] = t
			i++
		}
	}
	if len(other.Summary) != 0 {
		o.Summary = other.Summary
	}
	if len(other.Description) != 0 {
		o.Description = other.Description
	}
	if len(other.OperationID) != 0 {
		o.OperationID = other.OperationID
	}
	o.Deprecated = o.Deprecated || other.Deprecated
	if other.ExternalDocs != nil {
		o.ExternalDocs = other.ExternalDocs
	}
	o.Servers = append(o.Servers, other.Servers...)
	o.Responses.Merge(other.Responses)
	for k, cb := range other.Callbacks {
		if v, ok := o.Callbacks[k]; ok {
			for p, path := range cb {
				v[p] = path
			}
		} else {
			o.Callbacks[k] = cb
		}
	}
	o.Security = append(o.Security, other.Security...)
}

// ResponseSpec is an openapi Response description
type ResponseSpec struct {
	Description string               `json:"description"`
	Headers     map[string]Parameter `json:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Links       map[string]Link      `json:"links,omitempty"`
}

// Link descript a link to parts of a spec
type Link struct {
	OperationRef string         `json:"operationRef,omitempty"`
	OperationId  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
}

// RequestBody is the spec of a request body
type RequestBody struct {
	Description string               `json:"description"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType describes a media type in openapi
type MediaType struct {
	Schema   *jsonschema.Schema  `json:"schema,omitempty"`
	Encoding map[string]Encoding `json:"encoding,omitempty"`
	Example  any                 `json:"example,omitempty"`
	Examples *map[string]Example `json:"examples,omitempty"`
}

// Encoding is used to describe a content encoding in an API
type Encoding struct {
	ContentType   string               `json:"contentType,omitempty"`
	Headers       map[string]Parameter `json:"headers,omitempty"`
	Style         string               `json:"style,omitempty"`
	Explode       bool                 `json:"explode,omitempty"`
	AllowReserved bool                 `json:"allowReserved,omitempty"`
}

// Parameter describes a paramater used in requests/responses
type Parameter struct {
	Name            string              `json:"name,omitempty"`
	In              string              `json:"in,omitempty"`
	Description     string              `json:"description,omitempty"`
	Required        bool                `json:"required,omitempty"`
	Deprecated      bool                `json:"deprecated,omitempty"`
	AllowEmptyValue bool                `json:"allowEmptyValue,omitempty"`
	Style           string              `json:"style,omitempty"`
	Explode         bool                `json:"explode,omitempty"`
	AllowReserved   bool                `json:"allowReserved,omitempty"`
	Schema          *jsonschema.Schema  `json:"schema,omitempty"`
	Example         any                 `json:"example,omitempty"`
	Examples        *map[string]Example `json:"examples,omitempty"`
}

// Example is an example of any type
type Example struct {
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
	Example       any    `json:"example,omitempty"`
}

// Responses is a map of status code string to ResponseSpec
type Responses map[string]ResponseSpec

// Merge combines 2 responses objects into a single map
func (r *Responses) Merge(other Responses) {
	for k, v := range other {
		(*r)[k] = v
	}
}

// Components describes all openapi components
type Components struct {
	Schemas         map[string]jsonschema.Schema    `json:"schemas,omitempty"`
	Responses       Responses                       `json:"responses,omitempty"`
	Parameters      map[string]Parameter            `json:"parameters,omitempty"`
	Examples        map[string]Example              `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBody          `json:"requestBodies,omitempty"`
	Headers         map[string]map[string]Parameter `json:"headers,omitempty"`
	SecuritySchemes map[string]map[string][]string  `json:"securitySchemes,omitempty"`
	Links           map[string]Link                 `json:"links,omitempty"`
	Callbacks       map[string]map[string]Path      `json:"callbacks,omitempty"`
	PathItems       map[string]Path                 `json:"pathItems,omitempty"`
}

// Tag is used to tag parts of an API
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ExternalDocs is a link to external API documentation
type ExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// type Reference struct {
// 	Ref         string `json:"$ref"`
// 	Summary     string `json:"summary,omitempty"`
// 	Description string `json:"description,omitempty"`
// }
