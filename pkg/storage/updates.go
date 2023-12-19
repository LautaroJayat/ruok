package storage

import (
	"context"
	"database/sql"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

func (s *SQLStorage) StopListeningForChanges() error {
	ctx := context.Background()
	_, err := s.Db.Exec(ctx, "unlisten "+config.AppName())
	if err != nil {
		log.Error().Err(err).Msg("could not send unlisten command to db")
		return err
	}
	return nil
}

// Creates a gorutine that waits for messages in a loop and sends them over "jobIDUpdatedCh".
//
// It will block until the waiting loop starts
func (s *SQLStorage) ListenForChanges(jobIDUpdatedCh chan uuid.UUID, ctx context.Context) {
	ready := make(chan struct{})

	go func(jobIDUpdatedCh chan uuid.UUID, ctx context.Context) {
		ownChannel := config.AppName()
		conn, err := s.Db.Acquire(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("could not acquire connection to listen for notifications")
			return
		}
		defer conn.Release()

		_, err = conn.Exec(context.Background(), "listen "+config.AppName())

		if err != nil {
			log.Error().Err(err).Msgf("could not listen to %q channel", ownChannel)
			return
		}

		ready <- struct{}{}
		for {
			notification, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Info().Msgf("done listening for notifications. msg: %q", ctx.Err().Error())
					close(jobIDUpdatedCh)
					break
				}
				log.Error().Err(err).Msgf("an error occurred while listening into %q channel", ownChannel)
				close(jobIDUpdatedCh)
				break
			}
			if notification.Payload == "release" {
				log.Info().Msg("Received release, done!")
				return
			}
			if notification == nil {
				log.Info().Msg("empty notification")
				continue
			}
			id, err := uuid.FromString(notification.Payload)
			if err != nil {
				log.Error().Err(err).Msgf("an error occurred while parsing payload into job id %q for job update", notification.Payload)
				continue
			}
			jobIDUpdatedCh <- id
		}
		log.Debug().Msg("exiting from listening updates gorutine")
	}(jobIDUpdatedCh, ctx)

	<-ready
	close(ready)

}

const getJobsUpdatesQuery = `
SELECT
	job_name,
	cron_exp_string,
	endpoint,
	httpmethod,
	max_retries,
	headers_string,
	success_statuses,
	tls_client_cert,
	alert_strategy,
	alert_endpoint,
	alert_method,
	updated_at
FROM ruok.jobs
WHERE id = $1
`

type JobUpdates struct {
	Job_name         string
	Cron_exp_string  string
	Endpoint         string
	Httpmethod       string
	Max_retries      int
	Headers_string   string
	Success_statuses []int
	Tls_client_cert  string
	Alert_strategy   string
	Alert_endpoint   string
	Alert_method     string
	Updated_at       int64
}

func (s *SQLStorage) GetJobUpdates(jobId uuid.UUID) *JobUpdates {
	ctx := context.Background()
	tx, err := s.Db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("could not begin transaction to get updates for job %d", jobId)
		return nil
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, getJobsUpdatesQuery, jobId)
	var job_name string
	var cron_exp_string string
	var endpoint string
	var httpmethod string
	var max_retries int
	var headers_string sql.NullString
	var success_statuses []int
	var tls_client_cert sql.NullString
	var updated_at sql.NullInt64
	var alert_strategy sql.NullString
	var alert_endpoint sql.NullString
	var alert_method sql.NullString

	err = row.Scan(
		&job_name,
		&cron_exp_string,
		&endpoint,
		&httpmethod,
		&max_retries,
		&headers_string,
		&success_statuses,
		&tls_client_cert,
		&alert_strategy,
		&alert_endpoint,
		&alert_method,
		&updated_at,
	)

	if err != nil {
		log.Error().Err(err).Msgf("could not scan row to get updates for job %d", jobId)
		return nil
	}
	return &JobUpdates{
		job_name,
		cron_exp_string,
		endpoint,
		httpmethod,
		max_retries,
		headers_string.String,
		success_statuses,
		tls_client_cert.String,
		alert_strategy.String,
		alert_endpoint.String,
		alert_method.String,
		updated_at.Int64,
	}
}
