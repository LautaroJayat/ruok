package storage

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/job"
)

// Writes an execution result in the db
func (sqls *SQLStorage) WriteDone(j *job.Job) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to write job execution result")
		return errors.New("could not insert into jobs_results")
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
		claimed_by,
		succeeded
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
	`, j.Id, j.CronExpString, j.Endpoint, j.HttpMethod, j.MaxRetries, j.LastExecution.UnixMicro(),
		j.ShouldExecuteAt.UnixMicro(), j.LastResponseAt.UnixMicro(), j.LastMessage, j.LastStatusCode,
		j.SuccessStatuses, j.Status, j.ClaimedBy, j.Succeeded,
	)

	if err != nil {
		log.Error().Err(err).Msg("could not insert into job_results")
		return errors.New("could not insert into job_results")
	}

	_, err = tx.Exec(ctx, `
	UPDATE jobs SET
		last_execution = $1,
		should_execute_at = $2,
		last_response_at =$3,
		last_message = $4,
		last_status_code = $5,
		succeeded = $6
	WHERE id = $7
	`,
		j.LastExecution.UnixMicro(),
		j.ShouldExecuteAt.UnixMicro(),
		j.LastResponseAt.UnixMicro(),
		j.LastMessage,
		j.LastStatusCode,
		j.Succeeded,
		j.Id)

	if err != nil {
		log.Error().Err(err).Msg("could not update last execution fields for job")
		return errors.New("could not insert into job_results")
	}
	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not commit transaction to insert into job_results")
		return errors.New("could not commit transaction into job_results")
	}
	return nil
}
