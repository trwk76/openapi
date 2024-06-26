package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/trwk76/openapi/spec"
)

func New(e *gin.Engine, path string, security spec.SecurityRequirements, info spec.Info) *API {
	return &API{
		paths: newPaths(),
		g:     e.Group(path),
		s: spec.OpenAPI{
			OpenAPI: spec.Version,
			Info:    info,
			Servers: []spec.Server{{URL: path, Description: "Current server"}},
			Paths:   make(spec.Paths),
			Components: &spec.Components{
				Schemas:         make(spec.NamedSchemas),
				Responses:       make(spec.NamedResponseOrRefs),
				Parameters:      make(spec.NamedParameterOrRefs),
				RequestBodies:   make(spec.NamedRequestBodyOrRefs),
				Headers:         make(spec.NamedHeaderOrRefs),
				Examples:        make(spec.NamedExampleOrRefs),
				SecuritySchemes: make(spec.NamedSecuritySchemeOrRefs),
			},
			Security: security,
		},
		t: newSchemaEntries(),
	}
}

type (
	API struct {
		paths
		g *gin.RouterGroup
		s spec.OpenAPI
		t schemaEntries
	}
)
