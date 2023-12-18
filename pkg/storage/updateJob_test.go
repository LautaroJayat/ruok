package storage

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/stretchr/testify/assert"
)

func TestUpdateJob(t *testing.T) {
	defer Drop()
	tests := []struct {
		name       string
		job        UpdateJobInput
		expectErr  bool
		assertFunc func(*testing.T, *job.Job)
	}{
		{
			name: "UpdateJobWithNoAlerts",
			job: UpdateJobInput{
				Name:            "New Name",
				Id:              1,
				CronExpString:   "*/5 * * * *",
				Endpoint:        "/updated-test",
				HttpMethod:      "PUT",
				MaxRetries:      5,
				SuccessStatuses: []int{201},
			},
			expectErr: false,
			assertFunc: func(t *testing.T, j *job.Job) {
				assert.Equal(t, "New Name", j.Name)
				assert.Equal(t, "*/5 * * * *", j.CronExpString)
				assert.Equal(t, "/updated-test", j.Endpoint)
				assert.Equal(t, "PUT", j.HttpMethod)
				assert.Equal(t, 5, j.MaxRetries)
				assert.ElementsMatch(t, []int{201}, j.SuccessStatuses)
				assert.Equal(t, "claimed", j.Status)
			},
		},
		{
			name: "UpdateJobWithAlerts",
			job: UpdateJobInput{
				Id:              1,
				Name:            "New Name",
				CronExpString:   "*/10 * * * *",
				Endpoint:        "/updated-test",
				HttpMethod:      "PUT",
				MaxRetries:      3,
				SuccessStatuses: []int{200},
				AlertStrategy:   "sms",
				AlertEndpoint:   "123456789",
				AlertMethod:     "GET",
				AlertPayload:    "alert payload updated",
				AlertHeaders:    map[string]string{"Content-Type": "application/json"},
			},
			expectErr: false,
			assertFunc: func(t *testing.T, j *job.Job) {
				assert.Equal(t, "New Name", j.Name)
				assert.Equal(t, "*/10 * * * *", j.CronExpString)
				assert.Equal(t, "/updated-test", j.Endpoint)
				assert.Equal(t, "PUT", j.HttpMethod)
				assert.Equal(t, 3, j.MaxRetries)
				assert.ElementsMatch(t, []int{200}, j.SuccessStatuses)
				assert.Equal(t, "claimed", j.Status)
				assert.Equal(t, "sms", j.AlertStrategy)
				assert.Equal(t, "123456789", j.AlertEndpoint)
				assert.Equal(t, "GET", j.AlertMethod)
				assert.Equal(t, "alert payload updated", j.AlertPayload)
				assert.Equal(t, map[string]string{"Content-Type": "application/json"}, j.AlertHeaders)
			},
		},
	}

	cfg := config.FromEnvs()
	s, close := NewStorage(&cfg)
	defer close()

	// Create a job to update
	initialJob := CreateJobInput{
		CronExpString:   "*/2 * * * *",
		Endpoint:        "/initial-test",
		HttpMethod:      "POST",
		MaxRetries:      2,
		SuccessStatuses: []int{200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Drop()

			err := s.CreateJob(initialJob)
			assert.NoError(t, err, "failed to create initial job")

			jobs := s.GetAvailableJobs(1)
			assert.Len(t, jobs, 1)

			tt.job.Id = jobs[0].Id
			err = s.UpdateJob(tt.job)

			if tt.expectErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}

			updatedJobs := s.GetAvailableJobs(1)

			assert.Len(t, updatedJobs, 1)

			if tt.assertFunc != nil {
				tt.assertFunc(t, updatedJobs[0])
			}
		})
	}
}
