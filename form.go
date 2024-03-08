package chimera

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/invopop/jsonschema"
)

var (
	formBodyDecoder = form.NewDecoder()
)

// FormRequest[Body, Params any] is a request type that decodes request bodies to a
// user-defined struct for the Body and Params
type FormRequest[Body, Params any] struct {
	request *http.Request
	Body    Body
	Params  Params
}

// Context returns the context that was part of the original http.Request
func (r *FormRequest[Body, Params]) Context() context.Context {
	if r.request != nil {
		return r.request.Context()
	}
	return nil
}

// ReadRequest reads the body of an http request and assigns it to the Body field using
// http.Request.ParseForm and the "go-playground/form" package.
// This function also reads the parameters using UnmarshalParams and assigns it to the Params field.
// NOTE: the body of the request is closed after this function is run.
func (r *FormRequest[Body, Params]) ReadRequest(req *http.Request) error {
	defer req.Body.Close()
	err := req.ParseForm()
	if err != nil {
		return err
	}

	r.Body = *new(Body)
	if _, ok := any(r.Body).(Nil); !ok {
		err = formBodyDecoder.Decode(&r.Body, req.PostForm)
		if err != nil {
			return err
		}
	}

	r.Params = *new(Params)
	if _, ok := any(r.Params).(Nil); !ok {
		err = UnmarshalParams(req, &r.Params)
		if err != nil {
			return err
		}
	}

	r.request = req
	return nil
}

// flattenFormSchemas is kind of a jank way to convert jsonschema.Schema objects to be of a "form"
// style using patternProperties to represent arrays/object paths
func flattenFormSchemas(schema *jsonschema.Schema, properties map[string]*jsonschema.Schema, refs jsonschema.Definitions, prefix string) {
	if schema.Ref != "" && len(refs) > 0 {
		name := strings.Split(schema.Ref, "/")
		schema = refs[name[len(name)-1]]
		if schema != nil {
			delete(refs, name[len(name)-1])
		} else {
			properties["^"+prefix+".*$"] = &jsonschema.Schema{}
			return
		}
	}

	switch schema.Type {
	case "object":
		if prefix != "" {
			prefix += "."
		}
		for p := schema.Properties.Oldest(); p != nil; p = p.Next() {
			flattenFormSchemas(p.Value, properties, refs, prefix+p.Key)
		}
	case "array":
		flattenFormSchemas(schema.Items, properties, refs, prefix+"\\[\\d+\\]")
	default:
		properties["^"+prefix+"$"] = schema
	}
}

// OpenAPIRequestSpec returns the Request definition of a FormRequest
// It attempts to utilize patternProperties to try to define the body schema
// i.e. objects/arrays use dotted/bracketed paths X.Y.Z[i]
func (r *FormRequest[Body, Params]) OpenAPIRequestSpec() RequestSpec {
	bType := reflect.TypeOf(new(Body))
	for ; bType.Kind() == reflect.Pointer; bType = bType.Elem() {
	}

	schema := RequestSpec{}
	if bType != reflect.TypeOf(Nil{}) {
		s := (&jsonschema.Reflector{FieldNameTag: "form"}).Reflect(new(Body))
		// s.ID = jsonschema.ID(bType.PkgPath() + "_" + bType.Name())
		if s.PatternProperties == nil {
			s.PatternProperties = make(map[string]*jsonschema.Schema)
		}
		sType := s.Type
		if s.Ref != "" && len(s.Definitions) > 0 {
			name := strings.Split(s.Ref, "/")
			sType = s.Definitions[name[len(name)-1]].Type
		}
		flattenFormSchemas(s, s.PatternProperties, s.Definitions, "")
		s.Type = sType
		s.Ref = ""

		schema.RequestBody = &RequestBody{
			Content: map[string]MediaType{
				"application/x-www-form-urlencoded ": {
					Schema: s,
				},
			},
			Required: reflect.TypeOf(*new(Body)).Kind() != reflect.Pointer,
		}
	}

	pType := reflect.TypeOf(new(Params))
	for ; pType.Kind() == reflect.Pointer; pType = pType.Elem() {
	}
	if pType != reflect.TypeOf(Nil{}) {
		schema.Parameters = CacheRequestParamsType(pType)
	}
	// flatten form schemas
	return schema
}
