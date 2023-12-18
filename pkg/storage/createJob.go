package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
)

var createJobWithNoAlerts = `
INSERT INTO ruok.jobs (
	job_name,
	cron_exp_string,
	endpoint,
	httpmethod,
	max_retries,
	success_statuses,
	status
) VALUES ($1, $2, $3, $4, $5, $6, $7);
`

var createJobWithAlerts = `
INSERT INTO ruok.jobs (
	job_name,
	cron_exp_string,
	endpoint,
	httpmethod,
	max_retries,
	success_statuses,
	status,
	alert_strategy,
	alert_endpoint,
	alert_method,
	alert_headers_string,
	alert_payload
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
`

type CreateJobInput struct {
	Name            string            `json:"name"`
	CronExpString   string            `json:"cronexp"`
	MaxRetries      int               `json:"maxRetries"`
	Endpoint        string            `json:"endpoint"`
	HttpMethod      string            `json:"httpmethod"`
	Headers         map[string]string `json:"headers"`
	SuccessStatuses []int             `json:"successStatuses"`
	AlertStrategy   string            `json:"alertStrategy"`
	AlertMethod     string            `json:"alertMethod"`
	AlertEndpoint   string            `json:"alertEndpoint"`
	AlertPayload    string            `json:"alertPayload"`
	AlertHeaders    map[string]string `json:"alertHeaders"`
}

func (sqls *SQLStorage) CreateJob(j CreateJobInput) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to create job")
		return errors.New("could not insert into jobs")
	}

	if HasMinAlertFields(j.AlertStrategy, j.AlertEndpoint, j.AlertMethod) {
		var alertPayload sql.NullString
		if j.AlertPayload != "" {
			alertPayload.String = j.AlertPayload
			alertPayload.Valid = true
		}
		var alertHeadersString sql.NullString
		if len(j.AlertHeaders) > 0 {
			headersByte, err := json.Marshal(j.AlertHeaders)
			if err == nil {
				alertHeadersString.String = string(headersByte)
				alertHeadersString.Valid = true
			} else {
				log.Error().Err(err).Msgf("could not convert headers to json string")
			}
		}

		_, err = tx.Exec(ctx, createJobWithAlerts,
			j.Name,
			j.CronExpString,
			j.Endpoint,
			j.HttpMethod,
			j.MaxRetries,
			j.SuccessStatuses,
			"pending to be claimed",
			j.AlertStrategy,
			j.AlertEndpoint,
			j.AlertMethod,
			alertHeadersString,
			alertPayload,
		)
	} else {
		_, err = tx.Exec(ctx, createJobWithNoAlerts,
			j.Name,
			j.CronExpString,
			j.Endpoint,
			j.HttpMethod,
			j.MaxRetries,
			j.SuccessStatuses,
			"pending to be claimed",
		)

	}

	if err != nil {
		log.Error().Err(err).Msg("could not insert into jobs")
		return errors.New("could not insert into job")
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not commit transaction to insert into job")
		return errors.New("could not commit transaction into job")
	}
	return nil
}
