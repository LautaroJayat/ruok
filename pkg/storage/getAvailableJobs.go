package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
)

// Gets pending to be claimed jobs from the db and returns a list of all jobs that could be claimed
func (sqls *SQLStorage) GetAvailableJobs(limit int) []*job.Job {
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
	tls_client_cert,
	created_at
 FROM jobs 
 WHERE status = 'pending to be claimed' 
 LIMIT  $1;`, limit)

	if err != nil {
		fmt.Println("error", err)
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
		var TLSClientCert sql.NullString
		var CreatedAt int

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
			&TLSClientCert,
			&CreatedAt,
		)
		if err != nil {
			fmt.Println("error while scanning", err.Error())
		}

		Headers := []job.Header{}

		if HeadersString.Valid && HeadersString.String != "" {

			if err := json.Unmarshal([]byte(HeadersString.String), &Headers); err != nil {

				fmt.Printf("couldt unmarshal headers. error=%q\n", err.Error())

				jobsList = append(jobsList, &job.Job{
					Status: "bad headers",
					Id:     Id,
				})

				continue
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
			TLSClientCert:   TLSClientCert.String,
			ClaimedBy:       config.AppName(),
			Status:          "claimed",
			Handlers:        job.Handlers{},
			CreatedAt:       CreatedAt,
		}

		jobsList = append(jobsList, j)
	}

	rows.Close()

	for i := 0; i < len(jobsList); i++ {

		if jobsList[i].Status == "claimed" {
			_, err = tx.Exec(
				ctx,
				"UPDATE jobs SET claimed_by = $1, status = 'claimed' WHERE id = $2",
				jobsList[i].ClaimedBy,
				jobsList[i].Id,
			)
		} else {
			_, err = tx.Exec(
				ctx,
				"UPDATE jobs SET claimed_by = NULL, status = $1 WHERE id = $2",
				jobsList[i].Status,
				jobsList[i].Id,
			)

		}

		if err != nil {
			fmt.Println("error after exec: ", err.Error())
			return nil
		}

	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Printf("couldn't commit transaction. error=%q\n", err)
		return nil
	}

	return jobsList
}
