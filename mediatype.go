package openapi

type (
	// MediaTypes maps a content type (ex: application/json) to its media type implementation.
	MediaTypes map[string]MediaType

	MediaType interface {
		Key() string
		ExampleValue(value any) any
	}
)
