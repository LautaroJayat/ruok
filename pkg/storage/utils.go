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
    ('*/5 * * * *', 'https://localhost/api/job1', 'GET', 1, '{200}',  'pending to be claimed'),
    ('0 1 * * *', 'https://localhost/api/job2', 'GET', 1, '{200}',  'pending to be claimed'),
    ('30 2 * * *', 'https://localhost/api/job3', 'GET', 1, '{200}',  'pending to be claimed'),
    ('15 3 * * *', 'https://localhost/api/job4', 'GET', 1, '{200}',  'pending to be claimed'),
    ('0 */6 * * *', 'https://localhost/api/job5', 'GET', 1, '{200}',  'pending to be claimed'),
    ('*/10 * * * *', 'https://localhost/api/job6', 'GET', 1, '{200}',  'pending to be claimed'),
    ('45 4 * * *', 'https://localhost/api/job7', 'GET', 1, '{200}',  'pending to be claimed'),
    ('0 5 * * *', 'https://localhost/api/job8', 'GET', 1, '{200}',  'pending to be claimed'),
    ('*/15 * * * *', 'https://localhost/api/job9', 'GET', 1, '{200}',  'pending to be claimed'),
    ('0 6 * * *', 'https://localhost/api/job10', 'GET', 1, '{200}',  'pending to be claimed');
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
