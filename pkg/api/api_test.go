package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/back-end-labs/ruok/pkg/api/v1"
	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestHealthRoute(t *testing.T) {
	router := CreateRouter(nil)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/health", nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestStatusRoute(t *testing.T) {
	router := CreateRouter(nil)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/status", nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestClaimedJobs_BadParams(t *testing.T) {
	router := CreateRouter(nil)

	queries := []string{
		"offset=a1",
		"limit=a1",
		"limit=0&offset=a1",
		"limit=a1&offset=0",
		"limit=a1&offset=a1",
	}

	for _, query := range queries {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/jobs?"+query, nil)
		router.ServeHTTP(rr, req)
		assert.Equal(t, 400, rr.Code)
	}
}

func TestClaimedJobs_OKQueries(t *testing.T) {
	storage.Drop()
	storage.Seed()
	defer storage.Drop()
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()
	jobs := s.GetAvailableJobs(10)
	jobIds := []int{}
	for _, j := range jobs {
		jobIds = append(jobIds, j.Id)
	}
	router := CreateRouter(s)

	tests := []struct {
		query        string
		expectedJobs int
	}{
		{"offset=0", 10},
		{"limit=10", 10},
		{"limit=0&offset=10", 0},
		{"limit=10&offset=5", 5},
		{"limit=5&offset=5", 5},
		{"limit=10&offset=10", 0},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/jobs?"+test.query, nil)
		router.ServeHTTP(rr, req)
		assert.Equal(t, 200, rr.Code)
		body := &struct {
			Jobs []*job.Job `json:"jobs"`
		}{}
		err := json.Unmarshal(rr.Body.Bytes(), body)
		//t.Logf("query is %q, expected is %d, jobLen is %d", test.query, test.expectedJobs, len(body.Jobs))
		assert.Nil(t, err)
		assert.Equal(t, test.expectedJobs, len(body.Jobs))
		for _, jobFromBody := range body.Jobs {
			assert.Contains(t, jobIds, jobFromBody.Id)
		}
	}
}

func TestClaimedJobExecutions_BadParams(t *testing.T) {
	router := CreateRouter(nil)
	queries := []string{
		"offset=a1",
		"limit=a1",
		"limit=0&offset=a1",
		"limit=a1&offset=0",
		"limit=a1&offset=a1",
	}

	for _, query := range queries {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/jobs/1?"+query, nil)
		router.ServeHTTP(rr, req)
		assert.Equal(t, 400, rr.Code)
	}
}

func TestGetInstanceInfo(t *testing.T) {
	// Create a mock storage instance

	cfg := config.FromEnvs()

	expectedTLSActive := cfg.SSLConfigs.SSLMode == config.REQUIRE_SSL

	var expectedTLSVersion []string

	if expectedTLSActive {
		expectedTLSVersion = []string{"tlsv1.2", "tlsv1.2", "tlsv1.1"}
	} else {
		expectedTLSVersion = []string{""}

	}

	s, closeDb := storage.NewStorage(&cfg)
	defer closeDb()

	currentTime := time.Now()
	config.AppStats = &config.Stats{
		ClaimedJobs: 10,
		StartedAt:   currentTime.UnixMicro(),
	}
	router := CreateRouter(s)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/instance", nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	instanceInfo := &v1.InstanceInfo{}
	err := json.Unmarshal(rr.Body.Bytes(), &instanceInfo)
	t.Log(s.GetClient().Stat().TotalConns() > 0)
	assert.NoError(t, err, "Error unmarshaling JSON")
	assert.Equal(t, "application1", instanceInfo.AppName)
	assert.Equal(t, true, instanceInfo.DbConnected)
	assert.Equal(t, "postgresql://-:-@localhost:5432/db1", instanceInfo.DbUrl)
	assert.Equal(t, expectedTLSActive, instanceInfo.TlsActive)
	assert.Contains(t, expectedTLSVersion, instanceInfo.TlsVersion)
	assert.Equal(t, 10, instanceInfo.ClaimedJobs)
	assert.Equal(t, currentTime.UnixMicro(), instanceInfo.StartedAt)
	assert.True(t, instanceInfo.UpTimeMicro > 0) // Just ensure it's positive
	assert.Equal(t, 10000, instanceInfo.MaxJobs)
}
