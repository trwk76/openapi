package doc

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Init(g *gin.RouterGroup, apis ...API) {
	handle(g, "/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, g.BasePath()+"/index.html")
	})

	ents, err := docFS.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("error reading swaggerui embedded fs: %s", err.Error()))
	}

	fsrv := http.FS(docFS)

	for _, ent := range ents {
		if ent.Name() == "swagger-initializer.js" {
			raw, err := fs.ReadFile(docFS, ent.Name())
			if err != nil {
				panic(err)
			}

			urlsJSON, err := json.Marshal(apis)
			if err != nil {
				panic(fmt.Errorf("error marshaling api specifications: %s", err.Error()))
			}

			raw = []byte(strings.ReplaceAll(string(raw), "url: \"https://petstore.swagger.io/v2/swagger.json\"", "urls: "+string(urlsJSON)))
			hdl := func(ctx *gin.Context) {
				ctx.Writer.Header().Set("Content-Type", "text/javascript")
				ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(raw)))
				ctx.Writer.WriteHeader(http.StatusOK)

				if ctx.Request.Method == http.MethodGet {
					io.Copy(ctx.Writer, bytes.NewReader(raw))
				}
			}

			handle(g, ent.Name(), hdl)
		} else {
			g.StaticFileFS("/"+ent.Name(), ent.Name(), fsrv)
		}
	}
}

type (
	API struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url"`
	}
)

const Version string = "5.17.14"

func handle(g *gin.RouterGroup, path string, h gin.HandlerFunc) {
	g.GET(path, h)
	g.HEAD(path, h)
}

//go:embed *.css *.html *.js *.map
var docFS embed.FS
