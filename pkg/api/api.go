package api

import (
	"path"
	"runtime"

	v1 "github.com/back-end-labs/ruok/pkg/api/v1"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
)

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
	_, currentFile, _, _ := runtime.Caller(0)
	r.StaticFile("/", path.Dir(currentFile)+"/static/index.html")
	r.Static("/assets", path.Dir(currentFile)+"/static/assets")
	apiV1 := r.Group("/v1")
	{
		apiV1.Use(CORSMiddleware())
		apiV1.GET("/status", v1.Status)
		apiV1.GET("/health", v1.Health)
		apiV1.GET("/jobs", v1.ListJobs(apiStorage))
		apiV1.GET("/jobs/:id", v1.ListJobExections(apiStorage))
		apiV1.GET("/instance", v1.GetInstanceInfo(apiStorage))
	}
	return r
}
