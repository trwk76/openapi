package openapi

import (
	"reflect"

	"github.com/trwk76/openapi/spec"
)

func HeaderFor[T any](a *API, key string, examples Examples[T], setup func(s *spec.Header)) spec.HeaderOrRef {
	if len(examples) > 0 {
		setup = func(s *spec.Header) {
			s.Examples = examples.spec(nil)

			if setup != nil {
				setup(s)
			}
		}
	}

	return a.HeaderOf(key, reflect.TypeFor[T](), setup)
}

func (a *API) HeaderOf(key string, t reflect.Type, setup func(s *spec.Header)) spec.HeaderOrRef {
	sch := a.SchemaOf(t, nil, nil)
	itm := spec.Header{Schema: &sch}

	if setup != nil {
		setup(&itm)
	}

	if key != "" {
		key = uniqueName(a.s.Components.Headers, key)
		a.s.Components.Headers[key] = spec.HeaderOrRef{Item: itm}
		return spec.HeaderOrRef{Ref: spec.Ref("headers", key)}
	}

	return spec.HeaderOrRef{Item: itm}
}

func (a *API) Header(item spec.HeaderOrRef) spec.Header {
	if item.Ref.Ref != "" {
		return a.s.Components.Headers[refKey(item.Ref.Ref)].Item
	}

	return item.Item
}
