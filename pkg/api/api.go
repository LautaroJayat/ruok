package api

import (
	v1 "github.com/back-end-labs/ruok/pkg/api/v1"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
)

func CreateRouter(apiStorage storage.APIStorage) *gin.Engine {
	r := gin.Default()
	apiV1 := r.Group("/v1")
	{
		apiV1.GET("/status", v1.Status)
		apiV1.GET("/health", v1.Health)
		apiV1.GET("/jobs", v1.ListJobs(apiStorage))
		apiV1.GET("/jobs/:id", v1.ListJobExections(apiStorage))

	}
	return r
}
