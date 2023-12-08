package storage

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/stretchr/testify/assert"
)

func TestHasMinAlertFields(t *testing.T) {
	tests := []struct {
		name     string
		job      job.Job
		expected bool
	}{
		{
			name: "MinimumFieldsSet",
			job: job.Job{
				AlertStrategy: "email",
				AlertEndpoint: "test@example.com",
				AlertMethod:   "POST",
			},
			expected: true,
		},
		{
			name: "MissingAlertStrategy",
			job: job.Job{
				AlertEndpoint: "test@example.com",
				AlertMethod:   "POST",
			},
			expected: false,
		},
		{
			name: "MissingAlertEndpoint",
			job: job.Job{
				AlertStrategy: "email",
				AlertMethod:   "POST",
			},
			expected: false,
		},
		{
			name: "MissingAlertMethod",
			job: job.Job{
				AlertStrategy: "email",
				AlertEndpoint: "test@example.com",
			},
			expected: false,
		},
		{
			name:     "AllFieldsMissing",
			job:      job.Job{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasMinAlertFields(tt.job)
			assert.Equal(t, tt.expected, result)
		})
	}
}
