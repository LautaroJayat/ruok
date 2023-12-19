package v1

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidUrl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		pass  bool
	}{
		{
			name:  "ValidURL",
			input: "http://example.com",
			pass:  true,
		},
		{
			name:  "InvalidURL",
			input: "invalid-url",
			pass:  false,
		},
		{
			name:  "EmptyString",
			input: "",
			pass:  false,
		},
		{
			name:  "MalformedURL",
			input: "http://:8080",
			pass:  false,
		},
		{
			name:  "localUrl",
			input: "http://127.0.0.1:43055",
			pass:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validUrl(tt.input)
			assert.Equal(t, tt.pass, result, "expected result does not match")
		})
	}
}

func TestValidateUpdateFields(t *testing.T) {
	id1, _ := uuid.NewV7()

	tests := []struct {
		name          string
		input         storage.UpdateJobInput
		expectedError bool
		expectedList  []string
	}{
		{
			name: "ValidInput",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: false,
			expectedList:  nil,
		},
		{
			name: "MissingName",
			input: storage.UpdateJobInput{
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"must provide a name"},
		},
		{
			name: "MissingID",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid or missing id"},
		},
		{
			name: "InvalidCronExpression",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "invalid-expression",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid cron expression provided"},
		},
		{
			name: "MissingEndpoint",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"endpoint not found"},
		},
		{
			name: "InvalidEndpointURL",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "invalid-url",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid url provided"},
		},
		{
			name: "MissingHttpMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"missing http method"},
		},
		{
			name: "InvalidMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GOT",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid http method"},
		},
		{
			name: "MissingSuccessStatuses",
			input: storage.UpdateJobInput{
				Name:          "Job 1",
				Id:            id1,
				CronExpString: "*/1 * * * *",
				MaxRetries:    3,
				Endpoint:      "http://example.com",
				HttpMethod:    "GET",
			},
			expectedError: true,
			expectedList:  []string{"success statuses not provided"},
		},
		{
			name: "InvalidStrategy",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   "invalidOne",
			},
			expectedError: true,
			expectedList:  []string{"invalid strategy provided"},
		},
		{
			name: "InvalidAlertEndpoint",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertEndpoint:   ":8000",
			},
			expectedError: true,
			expectedList:  []string{"invalid alert endpoint provided"},
		},
		{
			name: "InvalidAlertMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				Id:              id1,
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertMethod:     "GOT",
			},
			expectedError: true,
			expectedList:  []string{"invalid alert http method provided"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, hasErrors := validateUpdateFields(tt.input)
			assert.Equal(t, tt.expectedError, hasErrors, "expected error status does not match")
			if tt.expectedError {
				assert.ElementsMatch(t, tt.expectedList, errors, "expected error list does not match")
			} else {
				assert.Empty(t, errors, "expected error list to be nil")
			}
		})
	}
}

func TestValidateCreateFields(t *testing.T) {
	tests := []struct {
		name          string
		input         storage.CreateJobInput
		expectedError bool
		expectedList  []string
	}{
		{
			name: "ValidInput",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: false,
			expectedList:  nil,
		},
		{
			name: "MissingName",
			input: storage.CreateJobInput{
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"must provide a name"},
		},
		{
			name: "InvalidCronExpression",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "invalid-expression",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid cron expression provided"},
		},
		{
			name: "MissingEndpoint",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"endpoint not found"},
		},
		{
			name: "InvalidEndpointURL",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "invalid-url",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid url provided"},
		},
		{
			name: "MissingHttpMethod",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"missing http method"},
		},
		{
			name: "InvalidMethod",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GOT",
				SuccessStatuses: []int{200},
			},
			expectedError: true,
			expectedList:  []string{"invalid http method"},
		},
		{
			name: "MissingSuccessStatuses",
			input: storage.CreateJobInput{
				Name:          "Job 1",
				CronExpString: "*/1 * * * *",
				MaxRetries:    3,
				Endpoint:      "http://example.com",
				HttpMethod:    "GET",
			},
			expectedError: true,
			expectedList:  []string{"success statuses not provided"},
		},
		{
			name: "InvalidStrategy",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   "invalidOne",
			},
			expectedError: true,
			expectedList:  []string{"invalid strategy provided"},
		},
		{
			name: "InvalidAlertEndpoint",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertEndpoint:   ":8000",
			},
			expectedError: true,
			expectedList:  []string{"invalid alert endpoint provided"},
		},
		{
			name: "InvalidAlertMethod",
			input: storage.CreateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertMethod:     "GOT",
			},
			expectedError: true,
			expectedList:  []string{"invalid alert http method provided"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, hasErrors := validateCreateFields(tt.input)
			assert.Equal(t, tt.expectedError, hasErrors, "expected error status does not match")
			if tt.expectedError {
				assert.ElementsMatch(t, tt.expectedList, errors, "expected error list does not match")
			} else {
				assert.Empty(t, errors, "expected error list to be nil")
			}
		})
	}
}

func TestValidHttpMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "ValidGET",
			input:    "GET",
			expected: true,
		},
		{
			name:     "ValidPOST",
			input:    "POST",
			expected: true,
		},
		{
			name:     "InvalidMethod",
			input:    "PUT",
			expected: false,
		},
		{
			name:     "EmptyString",
			input:    "",
			expected: false,
		},
		{
			name:     "MixedCaseGET",
			input:    "gEt",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validHttpMethod(tt.input)
			assert.Equal(t, tt.expected, result, "expected result does not match")
		})
	}
}
