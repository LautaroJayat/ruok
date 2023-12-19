package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type UpdateJobInput struct {
	Id              uuid.UUID         `json:"id"`
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

var updateJobQuery = `
UPDATE ruok.jobs SET 
	job_name = $1,
	cron_exp_string = $2,
	endpoint = $3,
	httpmethod = $4,
	max_retries = $5,
	success_statuses = $6,
	status = $7,
	alert_strategy = $8,
	alert_endpoint = $9,
	alert_method = $10,
	alert_headers_string = $11,
	alert_payload = $12,
	updated_at = ruok.micro_unix_now()
WHERE id = $13;
`

func (sqls *SQLStorage) UpdateJob(j UpdateJobInput) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to update job")
		return errors.New("could not update job")
	}

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

	_, err = tx.Exec(ctx, updateJobQuery,
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
		j.Id,
	)

	if err != nil {
		log.Error().Err(err).Msg("could not update job")
		return errors.New("could not update job")
	}

	_, err = tx.Exec(ctx, "select pg_notify($1, $2)", config.AppName(), j.Id.String())

	if err != nil {
		log.Error().Err(err).Msg("could not notify updated job")
		return errors.New("could not notify updated job")
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not commit transaction to update job job")
		return errors.New("could not commit transaction update job")
	}
	return nil
}
