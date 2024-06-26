package openapi

import (
	"reflect"
	"slices"
	"strings"

	"github.com/trwk76/openapi/spec"
)

type (
	JSONMediaType struct{}
)

func (JSONMediaType) Key() string {
	return "json"
}

func (JSONMediaType) ExampleValue(value any) any {
	return value
}

func (t JSONMediaType) ReflectField(api *API, fld reflect.StructField, bases *[]spec.SchemaOrRef, req *[]string, props *map[string]spec.SchemaOrRef) {
	tag := fld.Tag.Get("json")
	if tag == "-" {
		// safely ignore field
		return
	}

	flds := strings.Split(tag, ",")
	if flds[0] == "" && fld.Anonymous {
		// add base reference
		*bases = append(*bases, api.SchemaOrRefOf(fld.Type, t, nil))
		return
	}

	if flds[0] == "" {
		flds[0] = fld.Name
	}

	var (
		sch spec.SchemaOrRef
		opt bool = false
	)

	if len(flds) > 1 && slices.Contains(flds[1:], "omitempty") {
		opt = true
	}

	if len(flds) > 1 && slices.Contains(flds[1:], "string") {
		sch = spec.SchemaOrRef{Item: spec.Schema{Type: spec.TypeString}}
	} else {
		sch = api.SchemaOrRefOf(fld.Type, t, nil)
	}

	(*props)[flds[0]] = sch

	if !opt {
		*req = append(*req, flds[0])
	}
}

var JSON JSONMediaType

const ContentTypeJSON string = "application/json"

var _ MediaType = JSON
