package openapi

import (
	"reflect"

	"github.com/trwk76/openapi/spec"
)

func SchemaOrRefFor[T any](a *API, mediaType MediaType, examples Examples[T]) spec.SchemaOrRef {
	return a.SchemaOrRefOf(reflect.TypeFor[T](), mediaType, examples.spec(mediaType))
}

func Schema[T any](a *API, mediaType MediaType, examples Examples[T]) spec.Schema {
	return a.SchemaOf(reflect.TypeFor[T](), mediaType, examples.spec(mediaType))
}

func (a *API) SchemaOrRefOf(t reflect.Type, mediaType MediaType, examples spec.NamedExampleOrRefs) spec.SchemaOrRef {

}

func (a *API) SchemaOf(t reflect.Type, mediaType MediaType, examples spec.NamedExampleOrRefs) spec.Schema {

}

type (
	Examples[T any] map[string]Example[T]

	Example[T any] struct {
		Summary     string
		Description string
		Value       T
	}

	schemaEntries map[reflect.Type]schemaEntry

	schemaEntry struct {
		simple string
		mtkey  map[string]string
	}
)

func (e Examples[T]) spec(mtype MediaType) spec.NamedExampleOrRefs {
	res := make(spec.NamedExampleOrRefs)

	for key, itm := range e {
		res[key] = spec.ExampleOrRef{Item: itm.spec(mtype)}
	}

	return res
}

func (e Example[T]) spec(mtype MediaType) spec.Example {
	return spec.Example{
		Summary:     e.Summary,
		Description: e.Description,
		Value:       mtype.ExampleValue(e.Value),
	}
}

func newSchemaEntries() schemaEntries {
	return make(schemaEntries)
}
