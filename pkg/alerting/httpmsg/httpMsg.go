package httpmsg

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	m "github.com/back-end-labs/ruok/pkg/alerting/models"
	"github.com/back-end-labs/ruok/pkg/config"
)

var key = config.ALERT_HTTP

var ErrNoMethod = errors.New("bad alerting attempt, no method provided")

var ErrnoUrl = errors.New("bad alerting attempt, no url provided")

func httpAlert(input m.AlertInput) (string, error) {
	if input.Method == "" {
		return "", ErrNoMethod
	}
	if input.Url == "" {
		return "", ErrnoUrl
	}
	req, err := http.NewRequest(input.Method, input.Url, nil)

	if err != nil {
		return "", err
	}

	if len(input.Headers) > 0 {
		for k, v := range input.Headers {
			req.Header.Set(k, v)
		}
	}

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	stringBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("status code: %d", res.StatusCode), err
	}

	return fmt.Sprintf("status code: %d\nmessage: %s", res.StatusCode, string(stringBody)), err

}

func Plugin() (string, m.AlertFunc) {
	return key, httpAlert
}
