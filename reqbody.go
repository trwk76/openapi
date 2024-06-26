package openapi

import (
	"github.com/trwk76/openapi/spec"
)

func RequestBodyFor[T any](a *API, key string, mediaTypes MediaTypes, examples Examples[T], setup func(s *spec.RequestBody)) spec.RequestBodyOrRef {
	res := spec.RequestBody{Content: ContentFor(a, mediaTypes, examples)}

	if setup != nil {
		setup(&res)
	}

	if key != "" {
		key = uniqueName(a.s.Components.RequestBodies, key)
		a.s.Components.RequestBodies[key] = spec.RequestBodyOrRef{Item: res}
		return spec.RequestBodyOrRef{Ref: spec.Ref("requestBodies", key)}
	}

	return spec.RequestBodyOrRef{Item: res}
}

func (a *API) RequestBody(item spec.RequestBodyOrRef) spec.RequestBody {
	if item.Ref.Ref != "" {
		return a.s.Components.RequestBodies[refKey(item.Ref.Ref)].Item
	}

	return item.Item
}
