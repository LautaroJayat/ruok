package httpmsg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/back-end-labs/ruok/pkg/alerting/models"
	"github.com/stretchr/testify/assert"
)

func TestHTTPAlert(t *testing.T) {
	tests := []struct {
		name           string
		useOkUrl       bool
		alertInput     m.AlertInput
		expectedResult string
		expectedError  bool
	}{
		{
			name:     "Valid request",
			useOkUrl: true,
			alertInput: m.AlertInput{
				Url:     "",
				Method:  "GET",
				Headers: map[string]string{"Content-Type": "application/json"},
			},

			expectedResult: "status code: 200\nmessage: response body",
			expectedError:  false,
		},
		{
			name:           "Invalid request - Missing Method",
			useOkUrl:       true,
			alertInput:     m.AlertInput{},
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:     "Invalid request - Missing Url",
			useOkUrl: false,
			alertInput: m.AlertInput{
				Url:    "", // Invalid URL to trigger an error
				Method: "GET",
			},
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:     "Invalid request - Bad url",
			useOkUrl: false,
			alertInput: m.AlertInput{
				Url:    "invalid-url", // Invalid URL to trigger an error
				Method: "GET",
			},
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:     "Invalid Request - Unreachable Url",
			useOkUrl: false,
			alertInput: m.AlertInput{
				Url:    "http://example.com",
				Method: "GET",
			},
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "response body")
			}))
			defer server.Close()

			if tt.useOkUrl {
				tt.alertInput.Url = server.URL
			}

			result, err := httpAlert(tt.alertInput)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
