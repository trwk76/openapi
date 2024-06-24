package openapi

import (
	"fmt"

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

var YAML YAMLMediaType

const ContentTypeYAML string = "application/yaml"

var _ MediaType = YAML
