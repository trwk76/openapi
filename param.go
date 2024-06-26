package openapi

import (
	"reflect"

	"github.com/trwk76/openapi/spec"
)

func ParamFor[T any](a *API, key string, name string, in spec.ParameterIn, examples Examples[T], setup func(s *spec.Parameter)) spec.ParameterOrRef {
	if len(examples) > 0 {
		setup = func(s *spec.Parameter) {
			s.Examples = examples.spec(nil)

			if setup != nil {
				setup(s)
			}
		}
	}

	return a.ParamOf(key, name, in, reflect.TypeFor[T](), setup)
}

func (a *API) ParamOf(key string, name string, in spec.ParameterIn, t reflect.Type, setup func(s *spec.Parameter)) spec.ParameterOrRef {
	sch := a.SchemaOf(t, nil, nil)
	itm := spec.Parameter{Name: name, In: in, Schema: &sch}

	if setup != nil {
		setup(&itm)
	}

	if key != "" {
		key = uniqueName(a.s.Components.Parameters, key)
		a.s.Components.Parameters[key] = spec.ParameterOrRef{Item: itm}
		return spec.ParameterOrRef{Ref: spec.Ref("parameters", key)}
	}

	return spec.ParameterOrRef{Item: itm}
}

func (a *API) Param(item spec.ParameterOrRef) spec.Parameter {
	if item.Ref.Ref != "" {
		return a.s.Components.Parameters[refKey(item.Ref.Ref)].Item
	}

	return item.Item
}
