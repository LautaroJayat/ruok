package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

var limitLabel string = "limit"
var offsetLabel string = "offset"
var jobsLabel string = "jobs"
var jobResultsLabel string = "jobResults"
var jobIdLabel string = "jobId"
var claimedJobsLabel string = "claimedJobs"
var errorLabel string = "error"

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
			claimedJobsLabel: config.AppStats.CountClaimedJobs(),
			offsetLabel:      offset,
			limitLabel:       limit,
			jobsLabel:        jobslist,
		})

	}
}

func ListJobExecutions(s storage.APIStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitQ := c.DefaultQuery(limitLabel, "10")
		offsetQ := c.DefaultQuery(offsetLabel, "0")
		jobIdParam := c.Param("id")

		jobId, err := uuid.FromString(jobIdParam)
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

func CreateJob(s storage.APIStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var j storage.CreateJobInput
		if err := c.ShouldBindJSON(&j); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if j.AlertMethod != "" || j.AlertStrategy != "" || j.AlertEndpoint != "" {
			if !storage.HasMinAlertFields(j.AlertStrategy, j.AlertEndpoint, j.AlertMethod) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "if alerting is set, must provide strategy, endpoint and method"})
				return
			}
		}

		errors, hasErrors := validateCreateFields(j)

		if hasErrors {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors,
			})
			return
		}

		j.HttpMethod = strings.ToUpper(j.HttpMethod)

		if j.AlertMethod != "" && validHttpMethod(j.AlertMethod) {
			j.AlertMethod = strings.ToUpper(j.AlertMethod)

		}

		err := s.CreateJob(j)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "an internal error happened while trying to create a new job",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "job created",
		})

	}
}

func UpdateJob(s storage.APIStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var j storage.UpdateJobInput

		id, err := uuid.FromString(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{errorLabel: fmt.Sprintf("bad id provided: %s", c.Param("id"))})
			return
		}

		if err := c.ShouldBindJSON(&j); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{errorLabel: err.Error()})
			return
		}

		j.Id = id

		if j.AlertMethod != "" || j.AlertStrategy != "" || j.AlertEndpoint != "" {
			if !storage.HasMinAlertFields(j.AlertStrategy, j.AlertEndpoint, j.AlertMethod) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "if alerting is set, must provide strategy, endpoint and method"})
				return
			}
		}

		errors, hasErrors := validateUpdateFields(j)

		if hasErrors {
			c.JSON(http.StatusBadRequest, gin.H{
				errorLabel: errors,
			})
			return
		}

		err = s.UpdateJob(j)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				errorLabel: "an internal error happened while trying to create a new job",
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "job created",
		})

	}
}

type InstanceInfo struct {
	AppName     string `json:"appName"`
	DbConnected bool   `json:"dbConnected"`
	DbUrl       string `json:"dbURL"`
	TlsActive   bool   `json:"tlsActive"`
	TlsVersion  string `json:"tlsVersion"`
	ClaimedJobs int    `json:"claimedJobs"`
	StartedAt   int64  `json:"startedAtMicro"`
	UpTimeMicro int64  `json:"upTimeMicro"`
	MaxJobs     int    `json:"maxJobs"`
}

type apiStatsStorage interface {
	GetSSLVersion() (bool, string)
	Connected() bool
}

func GetInstanceInfo(s apiStatsStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.FromEnvs()
		dbConnected := s.Connected()
		tlsActive, tlsVersion := s.GetSSLVersion()
		payload := &InstanceInfo{
			cfg.AppName,
			dbConnected,
			fmt.Sprintf("%s://-:-@%s:%s/%s", cfg.Protocol, cfg.Host, cfg.Port, cfg.Dbname),
			tlsActive,
			tlsVersion,
			config.AppStats.ClaimedJobs,
			config.AppStats.StartedAt,
			config.AppStats.Uptime(),
			config.MaxJobs(),
		}

		c.JSON(200, &payload)
	}
}
