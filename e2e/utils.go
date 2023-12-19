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

func MakeTestURL(host string, jobName string) string {
	return fmt.Sprintf("%s/test?%s=%s", host, jobIdQuery, jobName)
}

func MakeAlertUrl(host string, jobName string) string {
	return fmt.Sprintf("%s/alert?%s=%s", host, jobIdQuery, jobName)
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

func UpdateJob(t *testing.T, id string, i storage.UpdateJobInput) bool {
	bytesBody, err := json.Marshal(i)
	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/jobs/"+id, bytes.NewReader(bytesBody))

	if err != nil {
		return false
	}

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return false
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		response, _ := io.ReadAll(res.Body)
		t.Logf("%d\n%s\n", res.StatusCode, string(response))
	}

	return res.StatusCode == http.StatusAccepted
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
