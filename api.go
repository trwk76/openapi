package openapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trwk76/openapi/spec"
	"gopkg.in/yaml.v3"
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

func (a *API) Finalize() {
	renderSpec(a.g, "json", a.s, json.Marshal)
	renderSpec(a.g, "yaml", a.s, yaml.Marshal)
}

func renderSpec(g *gin.RouterGroup, ext string, s spec.OpenAPI, m func(val any) ([]byte, error)) {
	raw, _ := m(s)

	hdl := func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/"+ext)
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(raw)))
		ctx.Writer.WriteHeader(http.StatusOK)

		if ctx.Request.Method == http.MethodGet {
			io.Copy(ctx.Writer, bytes.NewReader(raw))
		}
	}

	g.HEAD("/openapi."+ext, hdl)
	g.GET("/openapi."+ext, hdl)
}
