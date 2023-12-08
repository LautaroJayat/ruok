package storage

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/stretchr/testify/assert"
)

func TestCreateJob(t *testing.T) {
	defer Drop()
	tests := []struct {
		name       string
		job        job.Job
		expectErr  bool
		assertFunc func(*testing.T, *job.Job)
	}{
		{
			name: "CreateJobWithNoAlerts",
			job: job.Job{
				CronExpString:   "*/1 * * * *",
				Endpoint:        "/test",
				HttpMethod:      "GET",
				MaxRetries:      3,
				Status:          "pending to be claimed",
				SuccessStatuses: []int{200},
			},
			expectErr: false,
			assertFunc: func(t *testing.T, j *job.Job) {
				assert.Equal(t, "*/1 * * * *", j.CronExpString)
				assert.Equal(t, "/test", j.Endpoint)
				assert.Equal(t, "GET", j.HttpMethod)
				assert.Equal(t, 3, j.MaxRetries)
				assert.Equal(t, "claimed", j.Status)
				assert.ElementsMatch(t, []int{200}, j.SuccessStatuses)
			},
		},
		{
			name: "CreateJobWithAlerts",
			job: job.Job{
				CronExpString:   "*/1 * * * *",
				Endpoint:        "/test",
				HttpMethod:      "GET",
				MaxRetries:      3,
				SuccessStatuses: []int{200},
				Status:          "pending to be claimed",
				AlertStrategy:   "email",
				AlertEndpoint:   "test@example.com",
				AlertMethod:     "POST",
				AlertPayload:    "alert payload",
				AlertHeaders:    map[string]string{"Content-Type": "application/json"},
			},
			expectErr: false,
			assertFunc: func(t *testing.T, j *job.Job) {
				assert.Equal(t, "*/1 * * * *", j.CronExpString)
				assert.Equal(t, "/test", j.Endpoint)
				assert.Equal(t, "GET", j.HttpMethod)
				assert.Equal(t, 3, j.MaxRetries)
				assert.ElementsMatch(t, []int{200}, j.SuccessStatuses)
				assert.Equal(t, "claimed", j.Status)
				assert.Equal(t, "email", j.AlertStrategy)
				assert.Equal(t, "test@example.com", j.AlertEndpoint)
				assert.Equal(t, "POST", j.AlertMethod)
				assert.Equal(t, "alert payload", j.AlertPayload)
				assert.Equal(t, map[string]string{"Content-Type": "application/json"}, j.AlertHeaders)
			},
		},
	}

	cfg := config.FromEnvs()
	s, close := NewStorage(&cfg)
	defer close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Drop()

			err := s.CreateJob(tt.job)

			if tt.expectErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}

			createdJobs := s.GetAvailableJobs(100)

			assert.Len(t, createdJobs, 1)

			if tt.assertFunc != nil {
				tt.assertFunc(t, createdJobs[0])
			}
		})
	}
}
