package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateJob(t *testing.T) {
	originalMethods := os.Getenv(config.ALERT_CHANNELS)
	os.Setenv(config.ALERT_CHANNELS, config.ALERT_HTTP)
	defer os.Setenv(config.ALERT_CHANNELS, originalMethods)
	tests := []struct {
		name           string
		input          storage.UpdateJobInput
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "ValidInput",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  false,
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "ValidInputWithAlerts",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   config.ALERT_HTTP,
				AlertMethod:     "GET",
				AlertEndpoint:   "https://something.com",
				AlertHeaders:    map[string]string{"Authorization": "brearer jwt"},
				AlertPayload:    "",
			},
			expectedError:  false,
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "InvalidMissingName",
			input: storage.UpdateJobInput{
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidCronExpression",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "invalid-expression",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "MissingCronExpression",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "MissingEndpoint",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidEndpointURL",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "invalid-url",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "MissingHttpMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GOT",
				SuccessStatuses: []int{200},
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "MissingSuccessStatuses",
			input: storage.UpdateJobInput{
				Name:          "Job 1",
				CronExpString: "*/1 * * * *",
				MaxRetries:    3,
				Endpoint:      "http://example.com",
				HttpMethod:    "GET",
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidStrategy",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   "invalidOne",
				AlertMethod:     "GET",
				AlertEndpoint:   "https://something.com",
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidAlertEndpoint",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   config.ALERT_HTTP,
				AlertMethod:     "GOT",
				AlertEndpoint:   ":8000",
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidAlertMethod",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   config.ALERT_HTTP,
				AlertMethod:     "GOT",
				AlertEndpoint:   "https://something.com",
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "WithAlertButMissingMinFields",
			input: storage.UpdateJobInput{
				Name:            "Job 1",
				CronExpString:   "*/1 * * * *",
				MaxRetries:      3,
				Endpoint:        "http://example.com",
				HttpMethod:      "GET",
				SuccessStatuses: []int{200},
				AlertStrategy:   "invalidOne",
				AlertEndpoint:   "https://something.com",
			},
			expectedError:  true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	cfg := config.FromEnvs()
	defer storage.Drop()
	storage.Drop()
	s, close := storage.NewStorage(&cfg)
	defer close()
	router := CreateRouter(s)

	err := s.CreateJob(storage.CreateJobInput{
		CronExpString:   "*/1 * * * *",
		MaxRetries:      3,
		Endpoint:        "http://example.com",
		HttpMethod:      "GET",
		SuccessStatuses: []int{200},
	})
	assert.NoError(t, err)
	j := s.GetAvailableJobs(100)
	assert.Len(t, j, 1)
	jobId := j[0].Id

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			if err != nil {
				t.Log("there was an error creating request body")
				t.FailNow()
			}
			req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/jobs/%d", jobId), bytes.NewReader(body))
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Result().StatusCode)
		})
	}
}
