package openapi

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/trwk76/openapi/spec"
	"gopkg.in/yaml.v3"
)

type (
	YAMLMediaType struct{}
)

func (YAMLMediaType) Key() string {
	return "yaml"
}

func (YAMLMediaType) ExampleValue(value any) any {
	raw, err := yaml.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("error marshaling %v to yaml: %s", value, err.Error()))
	}

	return string(raw)
}

func (t YAMLMediaType) ReflectField(api *API, fld reflect.StructField, bases *[]spec.SchemaOrRef, req *[]string, props *map[string]spec.SchemaOrRef) {
	tag := fld.Tag.Get("yaml")
	if tag == "-" {
		// safely ignore field
		return
	}

	flds := strings.Split(tag, ",")
	if len(flds) > 1 && slices.Contains(flds[1:], "inline") {
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

var YAML YAMLMediaType

const ContentTypeYAML string = "application/yaml"

var _ MediaType = YAML
