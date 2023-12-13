package jobhandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/back-end-labs/ruok/pkg/job"
)

func TestHTTPExecutor_SuccessfulRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, "OK")
	}))
	defer server.Close()

	// Create a Job with the mock server URL
	testJob := &job.Job{
		HttpMethod: "GET",
		Endpoint:   server.URL,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}

	// Call the HTTPExecutor function
	result := HTTPExecutor(testJob)

	// Assert
	if result.Status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.Status)
	}

	if result.Message != "OK" {
		t.Errorf("Expected message 'OK', got '%s'", result.Message)
	}

	if result.SchedulerError != "" {
		t.Errorf("Expected no scheduler error, got '%s'", result.SchedulerError)
	}
}

func TestHTTPExecutor_ErrorCreatingRequest(t *testing.T) {
	// Create a Job with an invalid URL to simulate an error in creating the request
	testJob := &job.Job{
		HttpMethod: "GET",
		Endpoint:   "invalid-url",
	}

	// Call the HTTPExecutor function
	result := HTTPExecutor(testJob)

	// Assert
	if result.Status != 0 {
		t.Errorf("Expected status code 0, got %d", result.Status)
	}

	if result.Message != "" {
		t.Errorf("Expected empty message, got '%s'", result.Message)
	}

	if result.SchedulerError == "" {
		t.Error("Expected non-empty scheduler error, got empty")
	}
}

func TestHTTPExecutor_ErrorSendingRequest(t *testing.T) {
	// Create a Job with a server that always returns an error
	testJob := &job.Job{
		HttpMethod: "GET",
		Endpoint:   "https://bad.unreachable.example.com", // This URL will always return an error
	}

	// Call the HTTPExecutor function
	result := HTTPExecutor(testJob)

	// Assert
	if result.Status != 0 {
		t.Errorf("Expected status code 0, got %d", result.Status)
	}

	if result.Message != "" {
		t.Errorf("Expected empty message, got '%s'", result.Message)
	}

	if result.SchedulerError == "" {
		t.Error("Expected non-empty scheduler error, got empty")
	}
}
