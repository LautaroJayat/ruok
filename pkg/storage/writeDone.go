package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/back-end-labs/ruok/pkg/job"
)

// Writes an execution result in the db
func (sqls *SQLStorage) WriteDone(j *job.Job) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Printf("error=%q\n", err)
		return fmt.Errorf("could not insert into jobs_results. error=%q", err)
	}
	_, err = tx.Exec(ctx, `
	INSERT INTO job_results (
		job_id,
		cron_exp_string,
		endpoint,
		httpmethod,
		max_retries,
		execution_time,
		should_execute_at,
		last_response_at,
		last_message,
		last_status_code,
		success_statuses,
		status,
		claimed_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
	`, j.Id, j.CronExpString, j.Endpoint, j.HttpMethod, j.MaxRetries, j.LastExecution.UnixMicro(),
		j.ShouldExecuteAt.UnixMicro(), j.LastResponseAt.UnixMicro(), j.LastMessage, j.LastStatusCode,
		j.SuccessStatuses, j.Status, j.ClaimedBy,
	)
	if err != nil {
		fmt.Printf("There was a problem while trying to rollback query execution into job_results. error=%q", err)
		return fmt.Errorf("could not insert into job_results. error=%q", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		fmt.Printf("There was a problem while trying to rollback transaction into jobs_result. error=%q", err)
		return fmt.Errorf("could not commit transaction into job_results. error=%q", err)
	}
	return nil
}
