package storage

import (
	"context"
	"log"

	"github.com/back-end-labs/ruok/pkg/config"
)

var seedQuery string = `
INSERT INTO jobs (
    cron_exp_string,
    endpoint,
    httpmethod,
    max_retries,
    success_statuses,
    status
) VALUES
    ('*/5 * * * * * *', 'https://www.google.com', 'GET', 1, '{200}',  'pending to be claimed')
`

func Seed() {
	cfg := config.FromEnvs()
	s, close := NewStorage(&cfg)
	defer close()
	ctx := context.Background()
	tx, err := s.GetClient().Begin(ctx)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}

	_, err = tx.Exec(ctx, seedQuery)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("couldn't seed. error=%q", err)
	}
}

var dropJobsQuery string = "delete from jobs"
var dropJobResultsQuery string = "delete from job_results"

func Drop() {
	cfg := config.FromEnvs()
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
