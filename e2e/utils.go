package e2e

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	v1 "github.com/back-end-labs/ruok/pkg/api/v1"
	"github.com/back-end-labs/ruok/pkg/storage"
)

var jobIdQuery = "jobId"

func MakeTestURL(host string, jobName int, jobId int) string {
	return fmt.Sprintf("%s/test?%s=%d-%d", host, jobIdQuery, jobName, jobId)
}

func MakeAlertUrl(host string, jobName int, jobId int) string {
	return fmt.Sprintf("%s/alert?%s=%d-%d", host, jobIdQuery, jobName, jobId)
}

func CreateJob(t *testing.T, i storage.CreateJobInput) bool {
	bytesBody, err := json.Marshal(i)
	if err != nil {
		return false
	}
	res, err := http.Post("http://localhost:8080/v1/jobs", "", bytes.NewReader(bytesBody))
	if err != nil {
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		response, _ := io.ReadAll(res.Body)
		t.Logf("%d\n%s\n", res.StatusCode, string(response))
	}
	return res.StatusCode == http.StatusCreated
}

func ServerUp(t *testing.T) bool {
	res, err := http.Get("http://localhost:8080/v1/instance")
	if err != nil {
		return false
	}
	if res.StatusCode != 200 {
		return false
	}

	info := &v1.InstanceInfo{}
	err = json.NewDecoder(res.Body).Decode(info)
	defer res.Body.Close()
	if err != nil {
		return false
	}
	t.Log(info)
	return info.DbConnected
}

func ClaimedJobs(t *testing.T) (int, error) {
	res, err := http.Get("http://localhost:8080/v1/instance")
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, errors.New("couldn't get instance info")
	}
	info := &v1.InstanceInfo{}
	err = json.NewDecoder(res.Body).Decode(info)
	defer res.Body.Close()
	if err != nil {
		return 0, err
	}
	t.Log(info.ClaimedJobs, info.MaxJobs)
	return info.ClaimedJobs, nil
}
