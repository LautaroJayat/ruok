package api

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"

	v1 "github.com/back-end-labs/ruok/pkg/api/v1"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed static
var staticFiles embed.FS

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding,  Authorization, accept, origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CreateRouter(apiStorage storage.APIStorage) *gin.Engine {

	r := gin.Default()

	// Our API Router
	apiV1 := r.Group("/v1")
	{
		apiV1.Use(CORSMiddleware())
		apiV1.GET("/status", v1.Status)
		apiV1.GET("/health", v1.Health)
		apiV1.GET("/jobs", v1.ListJobs(apiStorage))
		apiV1.GET("/jobs/:id", v1.ListJobExecutions(apiStorage))
		apiV1.POST("/jobs", v1.CreateJob(apiStorage))
		apiV1.PUT("/jobs/:id", v1.UpdateJob(apiStorage))
		apiV1.GET("/instance", v1.GetInstanceInfo(apiStorage))
	}

	// For our SPA
	var indexBytes []byte
	{
		index, err := staticFiles.Open("static/index.html")

		// Ugly but it is a workaround
		if err != nil {
			log.Error().Err(err).Msg("couldn't open index")
			indexBytes = []byte{}
		} else {
			indexBytes, err = io.ReadAll(index)
			if err != nil {
				log.Error().Err(err).Msg("couldn't read index")
				indexBytes = []byte{}
			}
		}
	}

	// Easier as middleware
	r.Use(func(c *gin.Context) {
		// if it starts with "/v1" it is an api request
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/v1") {
			c.Redirect(http.StatusSeeOther, apiV1.BasePath())
			return
		}

		// Only allow harmless requests
		if c.Request.Method != http.MethodGet &&
			c.Request.Method != http.MethodOptions &&
			c.Request.Method != http.MethodHead {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}

		// Serve the index
		if path == "/index.html" || path == "/" || path == "" {
			c.Writer.Header().Add("Cache-Control", "no-cache")
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write(indexBytes)
			return
		}

		// Serve the assets
		f := strings.TrimPrefix(path, "/")
		f = fmt.Sprintf("static/%s", f)
		_, err := staticFiles.Open(f)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		c.Writer.Header().Add("Cache-Control", "max-age=31536000")
		c.FileFromFS(f, http.FS(staticFiles))
	})

	return r
}
