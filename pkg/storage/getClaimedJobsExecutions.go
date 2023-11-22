package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
)

// Gets get jobs claimed by this instance
func (sqls *SQLStorage) GetClaimedJobsExecutions(jobId int, limit int, offset int) []*job.JobExecution {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Printf("error=%v\n", err)
		return nil
	}

	rows, err := tx.Query(ctx, `
SELECT 
	id,
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
	success_statuses
 FROM job_results 
 WHERE claimed_by = $1 AND job_id = $2
 LIMIT  $3
 OFFSET $4;
 `, config.AppName(), jobId, limit, offset)

	if err != nil {
		fmt.Println("error", err)
		return nil

	}

	jobResultsList := []*job.JobExecution{}

	for rows.Next() {
		var Id int
		var JobId int
		var CronExpString string
		var Endpoint string
		var HttpMethod string
		var MaxRetries int
		var LastExecution sql.NullInt64
		var ShouldExecuteAt sql.NullInt64
		var LastResponseAt sql.NullInt64
		var LastMessage sql.NullString
		var LastStatusCode sql.NullInt32
		var SuccessStatuses []int

		err = rows.Scan(
			&Id,
			&JobId,
			&CronExpString,
			&Endpoint,
			&HttpMethod,
			&MaxRetries,
			&LastExecution,
			&ShouldExecuteAt,
			&LastResponseAt,
			&LastMessage,
			&LastStatusCode,
			&SuccessStatuses,
		)
		if err != nil {
			fmt.Println("error while scanning", err.Error())
		}

		j := &job.JobExecution{
			Id:              Id,
			JobId:           JobId,
			CronExpString:   CronExpString,
			Endpoint:        Endpoint,
			HttpMethod:      HttpMethod,
			LastExecution:   time.UnixMicro(LastExecution.Int64),
			ShouldExecuteAt: time.UnixMicro(ShouldExecuteAt.Int64),
			LastResponseAt:  time.UnixMicro(LastResponseAt.Int64),
			LastMessage:     LastMessage.String,
			LastStatusCode:  int(LastStatusCode.Int32),
			SuccessStatuses: SuccessStatuses,
			ClaimedBy:       config.AppName(),
		}

		jobResultsList = append(jobResultsList, j)
	}

	rows.Close()

	err = tx.Commit(ctx)

	if err != nil {
		log.Printf("couldn't commit transaction. error=%q\n", err)
		return nil
	}

	return jobResultsList
}
