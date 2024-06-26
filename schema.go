package openapi

import (
	"fmt"
	"go/token"
	"math"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/trwk76/openapi/spec"
)

func SchemaOrRefFor[T any](a *API, mediaType MediaType, examples Examples[T]) spec.SchemaOrRef {
	return a.SchemaOrRefOf(reflect.TypeFor[T](), mediaType, examples.spec(mediaType))
}

func Schema[T any](a *API, mediaType MediaType, examples Examples[T]) spec.Schema {
	return a.SchemaOf(reflect.TypeFor[T](), mediaType, examples.spec(mediaType))
}

func (a *API) SchemaOrRefOf(t reflect.Type, mediaType MediaType, examples spec.NamedExampleOrRefs) spec.SchemaOrRef {
	key := ""

	se, sefnd := a.t[t]
	if sefnd {
		if se.simple != "" {
			return spec.SchemaOrRef{Ref: spec.Ref("schemas", se.simple)}
		} else if mediaType != nil {
			if key, ok := se.mtkey[mediaType.Key()]; ok {
				return spec.SchemaOrRef{Ref: spec.Ref("schemas", key)}
			}
		}
	}

	if t.PkgPath() != "" && token.IsIdentifier(t.Name()) {
		key = uniqueName(a.s.Components.Schemas, strcase.ToLowerCamel(t.Name()))

		if se.mtkey == nil {
			se.mtkey = make(map[string]string)
		}

		if mediaType != nil {
			se.mtkey[mediaType.Key()] = key
		} else {
			se.simple = key
		}

		a.t[t] = se
	}

	res := a.SchemaOf(t, mediaType, examples)

	if key != "" {
		a.s.Components.Schemas[key] = res
		return spec.SchemaOrRef{Ref: spec.Ref("schemas", key)}
	}

	return spec.SchemaOrRef{Item: res}
}

func (a *API) SchemaOf(t reflect.Type, mediaType MediaType, examples spec.NamedExampleOrRefs) spec.Schema {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if sch, ok := reflect.Indirect(reflect.New(t)).Interface().(spec.Schemater); ok {
		return sch.Schema()
	}

	res := spec.Schema{}

	switch t.Kind() {
	case reflect.Array:
		items := a.SchemaOrRefOf(t.Elem(), mediaType, nil)

		res = spec.Schema{
			Type:     spec.TypeArray,
			Items:    &items,
			MinItems: uint32(t.Len()),
			MaxItems: uint32(t.Len()),
		}
	case reflect.Bool:
		res = spec.Schema{Type: spec.TypeBoolean}
	case reflect.Float32:
		res = numSchema(spec.TypeNumber, spec.FormatFloat, -math.MaxFloat32, math.MaxFloat32)
	case reflect.Float64:
		res = spec.Schema{Type: spec.TypeNumber, Format: spec.FormatDouble}
	case reflect.Int16:
		res = numSchema(spec.TypeInteger, spec.FormatNone, math.MinInt16, math.MaxInt16)
	case reflect.Int32:
		res = numSchema(spec.TypeInteger, spec.FormatInt32, math.MinInt32, math.MaxInt32)
	case reflect.Int64:
		res = numSchema(spec.TypeInteger, spec.FormatInt64, math.MinInt64, math.MaxInt64)
	case reflect.Int8:
		res = numSchema(spec.TypeInteger, spec.FormatNone, math.MinInt8, math.MaxInt8)
	case reflect.Map:
		key := a.SchemaOf(t.Key(), mediaType, nil)
		if key.Type != spec.TypeString {
			panic(fmt.Errorf("nap key type is not a string"))
		}

		item := a.SchemaOrRefOf(t.Elem(), mediaType, nil)
		res.Type = spec.TypeObject

		if key.Pattern != "" {
			res.PatternProperties = map[string]spec.SchemaOrRef{key.Pattern: item}
		} else {
			res.AdditionalProperties = &item
		}
	case reflect.Slice:
		items := a.SchemaOrRefOf(t.Elem(), mediaType, nil)

		res = spec.Schema{
			Type:  spec.TypeArray,
			Items: &items,
		}
	case reflect.String:
		res = spec.Schema{Type: spec.TypeString}
	case reflect.Struct:
		if mediaType == nil {
			panic(fmt.Errorf("type '%s' requires a media type since it is not a simple type", t.String()))
		}

		bases := make([]spec.SchemaOrRef, 0)
		req := make([]string, 0)
		props := make(map[string]spec.SchemaOrRef)

		for i := 0; i < t.NumField(); i++ {
			mediaType.ReflectField(a, t.Field(i), &bases, &req, &props)
		}

		res.Type = spec.TypeObject
		res.Required = req
		res.Properties = props

		if len(bases) > 0 {
			res = spec.Schema{AllOf: append(bases, spec.SchemaOrRef{Item: res})}
		}
	case reflect.Uint16:
		res = numSchema(spec.TypeInteger, spec.FormatNone, uint16(0), math.MaxUint16)
	case reflect.Uint32:
		res = numSchema(spec.TypeInteger, spec.FormatNone, uint32(0), math.MaxUint32)
	case reflect.Uint64:
		res = numSchema(spec.TypeInteger, spec.FormatNone, uint64(0), math.MaxInt64)
	case reflect.Uint8:
		res = numSchema(spec.TypeInteger, spec.FormatNone, uint8(0), math.MaxInt8)
	default:
		panic(fmt.Errorf("type '%s' cannot be reflected as a json schema", t.String()))
	}

	return res
}

func (a *API) Schema(item spec.SchemaOrRef) spec.Schema {
	if item.Ref.Ref != "" {
		return a.s.Components.Schemas[refKey(item.Ref.Ref)]
	}

	return item.Item
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

func numSchema[T comparable](t spec.Type, f spec.Format, min T, max T) spec.Schema {
	return spec.Schema{
		Type:    t,
		Format:  f,
		Minimum: min,
		Maximum: max,
	}
}
