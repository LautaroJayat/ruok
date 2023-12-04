package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
)

// Gets get jobs claimed by this instance
func (sqls *SQLStorage) GetClaimedJobs(limit int, offset int) []*job.Job {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to get claimed jobs")
		return nil
	}

	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
SELECT 
	id,
	cron_exp_string,
	endpoint,
	httpmethod,
	max_retries,
	last_execution,
	should_execute_at,
	last_response_at,
	last_message,
	last_status_code,
	headers_string,
	success_statuses,
	created_at,
	succeeded
 FROM ruok.jobs 
 WHERE claimed_by = $1 
 ORDER BY id ASC 
 LIMIT  $2
 OFFSET $3;
 `, config.AppName(), limit, offset)

	if err != nil {
		log.Error().Err(err).Msg("could not query for claimed jobs")
		return nil

	}

	jobsList := []*job.Job{}

	for rows.Next() {
		var Id int
		var CronExpString string
		var Endpoint string
		var HttpMethod string
		var MaxRetries int
		var LastExecution sql.NullInt64
		var ShouldExecuteAt sql.NullInt64
		var LastResponseAt sql.NullInt64
		var LastMessage sql.NullString
		var LastStatusCode sql.NullInt32
		var HeadersString sql.NullString
		var SuccessStatuses []int
		var CreatedAt int
		var Succeeded sql.NullString

		err = rows.Scan(
			&Id,
			&CronExpString,
			&Endpoint,
			&HttpMethod,
			&MaxRetries,
			&LastExecution,
			&ShouldExecuteAt,
			&LastResponseAt,
			&LastMessage,
			&LastStatusCode,
			&HeadersString,
			&SuccessStatuses,
			&CreatedAt,
			&Succeeded,
		)
		if err != nil {
			log.Error().Err(err).Msg("could not scan claimed jobs row")
			continue
		}

		Headers := []job.Header{}

		if HeadersString.Valid && HeadersString.String != "" {
			if err := json.Unmarshal([]byte(HeadersString.String), &Headers); err != nil {
				log.Error().Err(err).Msg("could not unmarshal headers of claimed job")
			}
		}

		j := &job.Job{
			Id:              Id,
			CronExpString:   CronExpString,
			Endpoint:        Endpoint,
			HttpMethod:      HttpMethod,
			MaxRetries:      MaxRetries,
			LastExecution:   time.UnixMicro(LastExecution.Int64),
			ShouldExecuteAt: time.UnixMicro(ShouldExecuteAt.Int64),
			LastResponseAt:  time.UnixMicro(LastResponseAt.Int64),
			LastMessage:     LastMessage.String,
			Headers:         Headers,
			LastStatusCode:  int(LastStatusCode.Int32),
			SuccessStatuses: SuccessStatuses,
			ClaimedBy:       config.AppName(),
			Handlers:        job.Handlers{},
			CreatedAt:       CreatedAt,
			Succeeded:       Succeeded.String,
		}

		jobsList = append(jobsList, j)
	}

	rows.Close()

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not commit 'get claimed jobs' transaction")
		return nil
	}

	return jobsList
}
