package storage

import (
	"context"
	"testing"
	"time"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestListenForChanges(t *testing.T) {
	cfg := config.FromEnvs()
	s, closeDbCon := NewStorage(&cfg)
	defer closeDbCon()
	ch := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	s.ListenForChanges(ch, ctx)
	signals := []string{"0", "1", "2", "3", "4", "5"}
	for i, sig := range signals {
		ctx := context.Background()
		_, err := s.GetClient().Exec(ctx, "select pg_notify($1, $2)", config.AppName(), sig)
		if err != nil {
			t.Errorf("could not send test message: %q", err.Error())
		}
		v := <-ch
		assert.Equal(t, i, v)
	}
	cancel()
	s.StopListeningForChanges()
}

func TestGetJobUpdates(t *testing.T) {
	var seedOneJobQuery = `
INSERT INTO jobs (
	id,
	cron_exp_string,
	endpoint,
	httpmethod,
	max_retries,
	success_statuses,
	status,
	claimed_by
) VALUES (1,'* * * * *', '/', 'GET', 1, '{200}',  'claimed','application1')
`
	cfg := config.FromEnvs()
	s, closeDbCon := NewStorage(&cfg)
	defer closeDbCon()
	_, err := s.GetClient().Exec(context.Background(), seedOneJobQuery)
	if err != nil {
		t.Errorf("couldn't seed one job for the test, %q", err.Error())
		t.FailNow()
	}

	new_cron_exp_string := "0 * * * * *"
	new_endpoint := "/slash"
	new_httpmethod := "POST"
	new_max_retries := 3
	new_headers_string := "{}"
	new_success_statuses := []int{200, 201}
	new_tls_client_cert := "a cert"
	new_updated_at := time.Now().UnixMicro()

	_, err = s.GetClient().Exec(context.Background(), `
		UPDATE jobs SET 
			cron_exp_string = $1,
			endpoint = $2,
			httpmethod = $3,
			max_retries = $4,
			headers_string = $5,
			success_statuses = $6,
			tls_client_cert = $7,
			updated_at = $8
		WHERE id = $9`,
		new_cron_exp_string,
		new_endpoint,
		new_httpmethod,
		new_max_retries,
		new_headers_string,
		new_success_statuses,
		new_tls_client_cert,
		new_updated_at,
		1,
	)
	if err != nil {
		t.Errorf("couldn't update one job for the test, %q", err.Error())
		t.FailNow()
	}
	j := s.GetJobUpdates(1)

	assert.NotNil(t, j, "GetJobUpdates should return a non-nil JobUpdates instance")
	assert.Equal(t, new_cron_exp_string, j.Cron_exp_string, "Unexpected cron_exp_string")
	assert.Equal(t, new_endpoint, j.Endpoint, "Unexpected endpoint")
	assert.Equal(t, new_httpmethod, j.Httpmethod, "Unexpected httpmethod")
	assert.Equal(t, new_max_retries, j.Max_retries, "Unexpected max_retries")
	assert.Equal(t, new_headers_string, j.Headers_string.String, "Unexpected headers_string")
	assert.Equal(t, new_success_statuses, j.Success_statuses, "Unexpected success_statuses")
	assert.Equal(t, new_tls_client_cert, j.Tls_client_cert.String, "Unexpected tls_client_cert")
	assert.Equal(t, new_updated_at, j.Updated_at, "Unexpected updated_at")
}

func TestStopListening(t *testing.T) {
	cfg := config.FromEnvs()
	s, closeDbCon := NewStorage(&cfg)
	err := s.StopListeningForChanges()
	assert.NoError(t, err)
	closeDbCon()
	assert.Error(t, s.StopListeningForChanges())
}
