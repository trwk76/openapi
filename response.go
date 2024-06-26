package openapi

import (
	"github.com/trwk76/openapi/spec"
)

func ResponseFor[T any](a *API, key string, desc string, mediaTypes MediaTypes, examples Examples[T], setup func(s *spec.Response)) spec.ResponseOrRef {
	res := spec.Response{Description: desc, Content: ContentFor(a, mediaTypes, examples)}

	if setup != nil {
		setup(&res)
	}

	if key != "" {
		key = uniqueName(a.s.Components.Responses, key)
		a.s.Components.Responses[key] = spec.ResponseOrRef{Item: res}
		return spec.ResponseOrRef{Ref: spec.Ref("responses", key)}
	}

	return spec.ResponseOrRef{Item: res}
}

func (a *API) Response(item spec.ResponseOrRef) spec.Response {
	if item.Ref.Ref != "" {
		return a.s.Components.Responses[refKey(item.Ref.Ref)].Item
	}

	return item.Item
}
