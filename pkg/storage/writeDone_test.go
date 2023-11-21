package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
)

func TestWriteDone(t *testing.T) {
	Drop()
	t.Run("Done jobs are written as they should", func(t *testing.T) {
		now := time.Now()
		j := job.Job{
			Id:              1,
			CronExpString:   "* * * * *",
			LastExecution:   now,
			ShouldExecuteAt: now,
			LastResponseAt:  now,
			Status:          "claimed",
			ClaimedBy:       config.AppName(),
			LastStatusCode:  200,
			HttpMethod:      "POST",
			Endpoint:        "/",
			SuccessStatuses: []int{200},
			Headers:         []job.Header{},
			MaxRetries:      1,
		}
		cfg := config.FromEnvs()
		s, close := NewStorage(&cfg)
		defer close()

		err := s.WriteDone(&j)
		if err != nil {
			t.Errorf("wrriting a job result shouldn't error. error=%q\n", err.Error())
		}
		ctx := context.Background()

		var (
			id              int64
			jobID           int64
			cronExpString   string
			endpoint        string
			httpMethod      string
			maxRetries      int
			executionTime   int64
			shouldExecuteAt int64
			lastResponseAt  int64
			lastMessage     sql.NullString
			lastStatusCode  int
			successStatuses []int
			tlsClientCert   sql.NullString
			status          string
			claimedBy       string
		)

		// Assuming jobId is the parameter to retrieve a specific job result.
		row := s.GetClient().QueryRow(ctx, "SELECT * FROM public.job_results WHERE job_id = $1", j.Id)
		err = row.Scan(
			&id,
			&jobID,
			&cronExpString,
			&endpoint,
			&httpMethod,
			&maxRetries,
			&executionTime,
			&shouldExecuteAt,
			&lastResponseAt,
			&lastMessage,
			&lastStatusCode,
			&successStatuses,
			&tlsClientCert,
			&status,
			&claimedBy,
		)
		if err != nil {
			t.Errorf("querying a done job should not produce an error. error=%q\n", err.Error())
		}

		if jobID != int64(j.Id) {
			t.Errorf("Expected JobID: %d, Got: %d", j.Id, jobID)
		}

		if cronExpString != j.CronExpString {
			t.Errorf("Expected CronExpString: %s, Got: %s", j.CronExpString, cronExpString)
		}

		if endpoint != j.Endpoint {
			t.Errorf("Expected Endpoint: %s, Got: %s", j.Endpoint, endpoint)
		}

		if httpMethod != j.HttpMethod {
			t.Errorf("Expected HttpMethod: %s, Got: %s", j.HttpMethod, httpMethod)
		}

		if maxRetries != j.MaxRetries {
			t.Errorf("Expected MaxRetries: %d, Got: %d", j.MaxRetries, maxRetries)
		}

		if executionTime != j.LastExecution.UnixMicro() {
			t.Errorf("Expected ExecutionTime: %d, Got: %d", j.LastExecution.UnixMicro(), executionTime)
		}

		if shouldExecuteAt != j.ShouldExecuteAt.UnixMicro() {
			t.Errorf("Expected ShouldExecuteAt: %d, Got: %d", j.ShouldExecuteAt.UnixMicro(), shouldExecuteAt)
		}

		if lastResponseAt != j.LastResponseAt.UnixMicro() {
			t.Errorf("Expected LastResponseAt: %d, Got: %d", j.LastResponseAt.UnixMicro(), lastResponseAt)
		}

		if lastMessage.String != "" && lastMessage.String != j.LastMessage {
			t.Errorf("Expected LastMessage: %s, Got: %s", j.LastMessage, lastMessage.String)
		}

		if lastStatusCode != j.LastStatusCode {
			t.Errorf("Expected LastStatusCode: %d, Got: %d", j.LastStatusCode, lastStatusCode)
		}

		if tlsClientCert.String != j.TLSClientCert {
			t.Errorf("Expected TlsClientCert: %s, Got: %s", j.TLSClientCert, tlsClientCert.String)
		}

		if status != j.Status {
			t.Errorf("Expected Status: %s, Got: %s", j.Status, status)
		}

		if claimedBy != j.ClaimedBy {
			t.Errorf("Expected ClaimedBy: %s, Got: %s", j.ClaimedBy, claimedBy)
		}
	})
	Drop()
}
