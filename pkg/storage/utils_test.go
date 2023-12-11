package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasMinAlertFields(t *testing.T) {
	tests := []struct {
		name     string
		job      CreateJobInput
		expected bool
	}{
		{
			name: "MinimumFieldsSet",
			job: CreateJobInput{
				AlertStrategy: "email",
				AlertEndpoint: "test@example.com",
				AlertMethod:   "POST",
			},
			expected: true,
		},
		{
			name: "MissingAlertStrategy",
			job: CreateJobInput{
				AlertEndpoint: "test@example.com",
				AlertMethod:   "POST",
			},
			expected: false,
		},
		{
			name: "MissingAlertEndpoint",
			job: CreateJobInput{
				AlertStrategy: "email",
				AlertMethod:   "POST",
			},
			expected: false,
		},
		{
			name: "MissingAlertMethod",
			job: CreateJobInput{
				AlertStrategy: "email",
				AlertEndpoint: "test@example.com",
			},
			expected: false,
		},
		{
			name:     "AllFieldsMissing",
			job:      CreateJobInput{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasMinAlertFields(tt.job.AlertStrategy, tt.job.AlertEndpoint, tt.job.AlertMethod)
			assert.Equal(t, tt.expected, result)
		})
	}
}
