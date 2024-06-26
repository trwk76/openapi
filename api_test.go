package openapi_test

import (
	"math"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/trwk76/openapi"
	"github.com/trwk76/openapi/spec"
)

func TestAPI(t *testing.T) {
	e := gin.New()

	mtypes := openapi.MediaTypes{
		openapi.ContentTypeJSON: openapi.JSON,
		openapi.ContentTypeYAML: openapi.YAML,
	}

	api := openapi.New(e, "/api/v1", spec.SecurityRequirements{{apiKey: {}}}, spec.Info{
		Title:       "TestAPI",
		Version:     "1.0.0",
		Summary:     "Test API",
		Description: "API for testing purpose",
	})

	parID := openapi.ParamFor[EntityID](api, "id", "id", spec.ParameterPath, nil, func(s *spec.Parameter) {
		s.Description = "Unique identifier of the entity"
	})

	res400 := openapi.ResponseFor(api, "err400", "Bad request", mtypes, openapi.Examples[Response]{
		"Invalid parameters": {
			Summary:     "Invalid parameters",
			Description: "One or more parameters failed validation.",
			Value: Response{
				CorrID: "00000000000000000000000000000000",
				Status: http.StatusBadRequest,
			},
		},
	}, nil)

	res401 := openapi.ResponseFor(api, "err401", "Unauthorized", mtypes, openapi.Examples[Response]{
		"Invalid credentials": {
			Summary:     "Invalid credentials",
			Description: "Credentials are invalid.",
			Value: Response{
				CorrID: "00000000000000000000000000000000",
				Status: http.StatusUnauthorized,
			},
		},
	}, nil)

	res404 := openapi.ResponseFor(api, "err404", "Not found", mtypes, openapi.Examples[Response]{
		"Entity not found": {
			Summary:     "Entity not found",
			Description: "No entity with this ID could be found.",
			Value: Response{
				CorrID: "00000000000000000000000000000000",
				Status: http.StatusNotFound,
			},
		},
	}, nil)

	tagUser := api.Tag("User", "User related methods")

	api.NamedPath("user", func(p *openapi.Path) {
		p.Tags = append(p.Tags, tagUser)

		p.NamedPath("id", func(p *openapi.Path) {
			p.ParamPath(parID, func(p *openapi.Path) {
				p.GET("Fetch", UserID, func(s *spec.Operation) {
					s.Summary = "Fetch User"
					s.Description = "Fetch a user by its identifier."
					s.Responses = spec.Responses{
						"400": res400,
						"401": res401,
						"404": res404,
					}
				})
			})
		})
	})

	api.Finalize()
}

func UserID(ctx *gin.Context) {

}

type (
	Response struct {
		CorrID CorrID `json:"corrId" yaml:"corrId"`
		Status Status `json:"status" yaml:"status"`
	}

	EntityID int64
	CorrID   string
	Status   uint16
)

func (EntityID) Schema() spec.Schema {
	return spec.Schema{
		Type:    spec.TypeString,
		Minimum: 1,
		Maximum: math.MaxInt64,
	}
}

func (CorrID) Schema() spec.Schema {
	return spec.Schema{
		Type:    spec.TypeString,
		Pattern: "^[0-9A-F]{32}$",
	}
}

func (Status) Schema() spec.Schema {
	return spec.Schema{
		Type:    spec.TypeInteger,
		Minimum: 200,
		Maximum: 599,
	}
}

const (
	apiKey   string = "api"
	basicKey string = "basic"
)
