package chimera

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/matt1484/spectagular"
)

// ResponseStructTag represents the "response" struct tag used by OneOfResponse
// it has a statusCode and a description
type ResponseStructTag struct {
	StatusCode  int    `structtag:"statusCode"`
	Description string `structtag:"description"`
}

var (
	_                   ResponseWriter = new(OneOfResponse[Nil])
	responseTagCache, _                = spectagular.NewFieldTagCache[ResponseStructTag]("response")
	responseWriterType                 = reflect.TypeOf((*ResponseWriter)(nil)).Elem()
)

// OneOfResponse[ResponseType any] is a response that uses the fields of
// ResponseType to determine which response to use as well as ResponseStructTag
// to control the status code, description of the different responses
// All fields must implement ResponseWriter to allow this to work properly
type OneOfResponse[ResponseType any] struct {
	Response ResponseType
}

// WriteResponse writes the response, content-type header, and status code using the first non-nil field
func (r *OneOfResponse[ResponseType]) WriteResponse(w http.ResponseWriter, ctx RouteContext) error {
	body := reflect.ValueOf(r.Response)
	tags, _ := responseTagCache.Get(body.Type())
	for _, tag := range tags {
		field := body.Field(tag.FieldIndex)
		if !field.IsNil() {
			field = fixPointer(field)
			return field.Interface().(ResponseWriter).WriteResponse(w, ctx.WithResponseCode(tag.Value.StatusCode))
		}
	}
	return nil
}

// OpenAPISpec returns the Responses definition of a OneOfResponse using all the OpenAPISpec() functions
// of the fields in ResponseType
func (r *OneOfResponse[ResponseType]) OpenAPISpec() Responses {
	schema := make(Responses)
	body := reflect.ValueOf(*new(ResponseType))
	tags, err := responseTagCache.GetOrAdd(body.Type())
	if err != nil {
		panic("chimera.OneOfResponse[Body]: Had invalid Body type: " + body.Type().Name())
	}
	for _, tag := range tags {
		val := body.Field(tag.FieldIndex)
		if val.Kind() == reflect.Pointer {
			val = reflect.New(val.Type().Elem())
		}
		val = fixPointer(val)
		if !val.Type().Implements(responseWriterType) {
			panic("chimera.OneOfResponse[Body]: Body fields MUST implement chimera.ResponseWriter")
		}
		resp := val.Interface().(ResponseWriter).OpenAPISpec()
		if v, ok := resp[""]; ok {
			v.Description = tag.Value.Description
			resp[fmt.Sprint(tag.Value.StatusCode)] = v
			delete(resp, "")
		}
		schema.Merge(resp)
	}
	return schema
}

// NewOneOfResponse creates a OneOfResponse from a response
func NewOneOfResponse[ResponseType any](response ResponseType) *OneOfResponse[ResponseType] {
	return &OneOfResponse[ResponseType]{
		Response: response,
	}
}
