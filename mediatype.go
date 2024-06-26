package openapi

import (
	"reflect"

	"github.com/trwk76/openapi/spec"
)

type (
	// MediaTypes maps a content type (ex: application/json) to its media type implementation.
	MediaTypes map[string]MediaType

	MediaType interface {
		Key() string
		ExampleValue(value any) any
		ReflectField(api *API, fld reflect.StructField, bases *[]spec.SchemaOrRef, req *[]string, props *map[string]spec.SchemaOrRef)
	}
)
