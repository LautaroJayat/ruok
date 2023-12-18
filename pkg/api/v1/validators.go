package v1

import (
	"net/url"
	"strings"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/cronParser"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func validHttpMethod(input string) bool {
	switch strings.ToUpper(input) {
	case "GET":
		return true
	case "POST":
		return true
	default:
		return false
	}
}

func validUrl(input string) bool {
	url, err := url.ParseRequestURI(input)
	if err != nil || url.Host == "" || url.Scheme == "" || string(url.Host[0]) == ":" {
		return false
	}
	return true
}

func validateCreateFields(j storage.CreateJobInput) ([]string, bool) {
	hasErrors := false
	errors := []string{}

	if j.Name == "" {
		hasErrors = true
		errors = append(errors, "must provide a name")
	}

	if j.CronExpString == "" {
		hasErrors = true
		errors = append(errors, "cron expression string not found")

	} else if cronParser.IsValidExpression(j.CronExpString) {
		hasErrors = true
		errors = append(errors, "invalid cron expression provided")
	}

	if j.Endpoint == "" {
		hasErrors = true
		errors = append(errors, "endpoint not found")

	} else if !validUrl(j.Endpoint) {
		hasErrors = true
		errors = append(errors, "invalid url provided")
	}

	if j.HttpMethod == "" {
		hasErrors = true
		errors = append(errors, "missing http method")
	} else if !validHttpMethod(j.HttpMethod) {
		hasErrors = true
		errors = append(errors, "invalid http method")
	}

	if len(j.SuccessStatuses) == 0 {
		hasErrors = true
		errors = append(errors, "success statuses not provided")
	}
	if j.AlertStrategy != "" && badAlertStrategy(j.AlertStrategy, config.AlertChannels()) {
		hasErrors = true
		errors = append(errors, "invalid strategy provided")
	}
	if j.AlertEndpoint != "" && !validUrl(j.AlertEndpoint) {
		hasErrors = true
		errors = append(errors, "invalid alert endpoint provided")
	}
	if j.AlertMethod != "" && !validHttpMethod(j.AlertMethod) {
		hasErrors = true
		errors = append(errors, "invalid alert http method provided")
	}
	return errors, hasErrors
}

func badAlertStrategy(ch string, valids []string) bool {
	for _, v := range config.AlertChannels() {
		if ch == v {
			return false
		}
	}
	return true

}

func validateUpdateFields(j storage.UpdateJobInput) ([]string, bool) {
	hasErrors := false
	errors := []string{}

	if j.Id == 0 {
		hasErrors = true
		errors = append(errors, "invalid or missing id")
	}

	if j.Name == "" {
		hasErrors = true
		errors = append(errors, "must provide a name")
	}

	if j.CronExpString == "" {
		hasErrors = true
		errors = append(errors, "cron expression string not found")

	} else if cronParser.IsValidExpression(j.CronExpString) {
		hasErrors = true
		errors = append(errors, "invalid cron expression provided")
	}

	if j.Endpoint == "" {
		hasErrors = true
		errors = append(errors, "endpoint not found")

	} else if !validUrl(j.Endpoint) {
		hasErrors = true
		errors = append(errors, "invalid url provided")
	}

	if j.HttpMethod == "" {
		hasErrors = true
		errors = append(errors, "missing http method")
	} else if !validHttpMethod(j.HttpMethod) {
		hasErrors = true
		errors = append(errors, "invalid http method")
	}

	if len(j.SuccessStatuses) == 0 {
		hasErrors = true
		errors = append(errors, "success statuses not provided")
	}

	if j.AlertStrategy != "" && badAlertStrategy(j.AlertStrategy, config.AlertChannels()) {
		hasErrors = true
		errors = append(errors, "invalid strategy provided")
	}
	if j.AlertEndpoint != "" && !validUrl(j.AlertEndpoint) {
		hasErrors = true
		errors = append(errors, "invalid alert endpoint provided")
	}
	if j.AlertMethod != "" && !validHttpMethod(j.AlertMethod) {
		hasErrors = true
		errors = append(errors, "invalid alert http method provided")
	}

	return errors, hasErrors
}
