package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/gofrs/uuid"
)

func seedQuery() string {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	id3, _ := uuid.NewV7()
	id4, _ := uuid.NewV7()
	id5, _ := uuid.NewV7()
	id6, _ := uuid.NewV7()
	id7, _ := uuid.NewV7()
	id8, _ := uuid.NewV7()
	id9, _ := uuid.NewV7()
	id10, _ := uuid.NewV7()
	return fmt.Sprintf(`
	INSERT INTO ruok.jobs (
		id,
		job_name,
		cron_exp_string,
		endpoint,
		httpmethod,
		max_retries,
		success_statuses,
		status
		) VALUES
		('%s','job 1', '*/5 * * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 2', '0 1 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 3', '30 2 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 4', '15 3 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 5', '0 */6 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 6', '*/10 * * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 7', '45 4 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 8', '0 5 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 9', '*/15 * * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed'),
		('%s','job 10', '0 6 * * *', 'http://localhost:8080/v1/status', 'GET', 1, '{200}',  'pending to be claimed');
		`, id1, id2, id3, id4, id5, id6, id7, id8, id9, id10)
}

func Seed() {
	cfg := config.FromEnvs()
	// use a testing role with all privileges
	cfg.User = "testing_user"
	s, close := NewStorage(&cfg)
	defer close()
	ctx := context.Background()
	tx, err := s.GetClient().Begin(ctx)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}

	_, err = tx.Exec(ctx, seedQuery())
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}
}

var dropJobsQuery string = "delete from ruok.jobs"
var dropJobResultsQuery string = "delete from ruok.job_results"

func Drop() {
	cfg := config.FromEnvs()
	// use a testing role with all privileges
	cfg.User = "testing_user"
	s, close := NewStorage(&cfg)
	defer close()
	ctx := context.Background()
	tx, err := s.GetClient().Begin(ctx)

	if err != nil {
		log.Fatalf("couldn't init transaction. error=%q", err)
	}

	_, err = tx.Exec(ctx, dropJobsQuery)
	if err != nil {
		log.Fatalf("couldn't delete jobs. error=%q", err)
	}

	_, err = tx.Exec(ctx, dropJobResultsQuery)
	if err != nil {
		log.Fatalf("couldn't delete job results. error=%q", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}
}

func HasMinAlertFields(strategy string, endpoint string, method string) bool {
	if strategy == "" {
		return false
	}
	if endpoint == "" {
		return false
	}
	if method == "" {
		return false
	}
	return true
}

func seedOneJobQuery(id uuid.UUID) string {
	return fmt.Sprintf(`
	INSERT INTO ruok.jobs (
		id,
		job_name,
		cron_exp_string,
		endpoint,
		httpmethod,
		max_retries,
		success_statuses,
		status,
		claimed_by
	) VALUES ('%s', 'testing job', '* * * * *', '/', 'GET', 1, '{200}',  'claimed','application1')
	`, id.String())
}
