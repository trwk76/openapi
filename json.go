package openapi

type (
	JSONMediaType struct{}
)

func (JSONMediaType) Key() string {
	return "json"
}

func (JSONMediaType) ExampleValue(value any) any {
	return value
}

var JSON JSONMediaType

const ContentTypeJSON string = "application/json"

var _ MediaType = JSON
