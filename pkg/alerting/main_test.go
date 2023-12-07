package alerting

import (
	"errors"
	"testing"

	"github.com/back-end-labs/ruok/pkg/alerting/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateAlertManager(t *testing.T) {

	dummyfn := func(i models.AlertInput) (string, error) {
		_ = i
		return "", nil

	}

	tests := []struct {
		name              string
		availableChannels []string
		plugin            models.AlertPlugin
		expectedCount     int
	}{
		{
			name:              "ValidAlertFn",
			availableChannels: []string{"http"},
			plugin:            func() (string, models.AlertFunc) { return "http", dummyfn },
			expectedCount:     1,
		},
		{
			name:              "InvalidAlertFn",
			availableChannels: []string{"email"},
			plugin:            func() (string, models.AlertFunc) { return "invalid", dummyfn },
			expectedCount:     0,
		},
		{
			name:              "MultipleAlertFns",
			availableChannels: []string{"http", "email"},
			plugin:            func() (string, models.AlertFunc) { return "http", dummyfn },
			expectedCount:     1,
		},
		{
			name:              "NoAlertFns",
			availableChannels: []string{"sms", "email"},
			plugin:            func() (string, models.AlertFunc) { return "http", dummyfn },
			expectedCount:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replace the alertFn slice with the test case's alert function
			registeredFn := models.PluginList{tt.plugin}

			// Call the CreateAlertManager function with availableChannels
			alertManager := CreateAlertManager(tt.availableChannels, registeredFn)

			// Check if the alert manager has the expected count of alert strategies
			assert.Len(t, alertManager.alertStrategies, tt.expectedCount)
		})
	}
}

func TestSendAlert(t *testing.T) {
	validStrategy := "http"
	invalidStrategy := "invalid"
	noErrorResult := "dummyAlertFuncNoErrorResult"
	errorResult := errors.New("err")
	dummyAlertFuncNoError := func(input models.AlertInput) (string, error) {
		return noErrorResult, nil
	}

	dummyAlertFuncWithError := func(input models.AlertInput) (string, error) {
		return errorResult.Error(), errorResult
	}
	alertManager := &AlertManager{
		alertStrategies: map[string]models.AlertFunc{
			validStrategy: nil,
		},
	}

	tests := []struct {
		name           string
		mockFn         models.AlertFunc
		alertInput     models.AlertInput
		expectedResult string
		expectedStatus int
	}{
		{
			name:   "ValidAlertStrategy",
			mockFn: dummyAlertFuncNoError,
			alertInput: models.AlertInput{
				AlertStrategy:  validStrategy,
				Url:            "http://example.com",
				Method:         "GET",
				Payload:        "",
				ExpectedStatus: 200,
				ExpectedMsg:    "OK",
				Headers:        nil,
			},
			expectedResult: noErrorResult,
			expectedStatus: STATUS_OK,
		},
		{
			name:   "InvalidAlertStrategy",
			mockFn: dummyAlertFuncNoError,
			alertInput: models.AlertInput{
				AlertStrategy:  invalidStrategy,
				Url:            "http://example.com",
				Method:         "GET",
				Payload:        "",
				ExpectedStatus: 200,
				ExpectedMsg:    "OK",
				Headers:        nil,
			},
			expectedResult: "",
			expectedStatus: STATUS_FN_NOT_REGISTERED,
		},
		{
			name:   "ErrorWhileSending",
			mockFn: dummyAlertFuncWithError,
			alertInput: models.AlertInput{
				AlertStrategy:  validStrategy,
				Url:            "http://example.com",
				Method:         "GET",
				Payload:        "",
				ExpectedStatus: 200,
				ExpectedMsg:    "OK",
				Headers:        nil,
			},
			expectedResult: errorResult.Error(),
			expectedStatus: STATUS_ERR_WHILE_SENDING,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alertManager.alertStrategies[validStrategy] = tt.mockFn
			result, status := alertManager.SendAlert(tt.alertInput)
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}
