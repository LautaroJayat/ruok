package v1

import (
	"fmt"
	"strconv"

	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
)

var limitLabel string = "limit"
var offsetLabel string = "offset"
var jobsLabel string = "jobs"
var jobResultsLabel string = "jobResults"
var jobIdLabel string = "jobId"

func Status(c *gin.Context) {
	c.String(200, "OK")
}

func Health(c *gin.Context) {
	c.String(200, "OK")
}

func BadQueryError(query string, value string) string {
	return fmt.Sprintf("Bad request. %q needs to be an integer, instead got %q\n", query, value)
}

func ListJobs(s storage.APIStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitQ := c.DefaultQuery(limitLabel, "10")
		offsetQ := c.DefaultQuery(offsetLabel, "0")
		limit, err := strconv.Atoi(limitQ)
		if err != nil {
			c.JSON(400, gin.H{
				"message": BadQueryError(limitLabel, limitQ),
			})
			return
		}

		offset, err := strconv.Atoi(offsetQ)
		if err != nil {
			c.JSON(400, gin.H{
				"message": BadQueryError(offsetLabel, offsetQ),
			})
			return
		}

		jobslist := s.GetClaimedJobs(limit, offset)

		c.JSON(200, gin.H{
			limitLabel:  limit,
			offsetLabel: offset,
			jobsLabel:   jobslist,
		})

	}
}

func ListJobExections(s storage.APIStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitQ := c.DefaultQuery(limitLabel, "10")
		offsetQ := c.DefaultQuery(offsetLabel, "0")
		jobIdParam := c.Param("id")

		jobId, err := strconv.Atoi(jobIdParam)
		if err != nil {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Could not found jobs with id %v", jobId),
			})
			return
		}

		limit, err := strconv.Atoi(limitQ)
		if err != nil {
			c.JSON(400, gin.H{
				"message": BadQueryError(limitLabel, limitQ),
			})
			return
		}

		offset, err := strconv.Atoi(offsetQ)
		if err != nil {
			c.JSON(400, gin.H{
				"message": BadQueryError(offsetLabel, offsetQ),
			})
			return
		}

		jobExecutionList := s.GetClaimedJobsExecutions(jobId, limit, offset)

		if len(jobExecutionList) == 0 {
			c.JSON(404, gin.H{
				limitLabel:      limit,
				offsetLabel:     offset,
				jobIdLabel:      jobIdParam,
				jobResultsLabel: jobExecutionList,
			})
			return
		}

		c.JSON(200, gin.H{
			limitLabel:      limit,
			offsetLabel:     offset,
			jobIdParam:      jobIdParam,
			jobResultsLabel: jobExecutionList,
		})

	}
}
