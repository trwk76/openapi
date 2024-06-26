package openapi

import (
	"github.com/trwk76/openapi/spec"
)

func ContentFor[T any](a *API, mediaTypes MediaTypes, examples Examples[T]) spec.MediaTypes {
	res := make(spec.MediaTypes)

	for ctype, mtype := range mediaTypes {
		sch := SchemaOrRefFor[T](a, mtype, nil)

		res[ctype] = spec.MediaType{
			Schema:   &sch,
			Examples: examples.spec(mtype),
		}
	}

	return res
}
