package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/job"
)

var updateJobQuery = `
UPDATE ruok.jobs SET 
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
WHERE id = $1;
`

func (sqls *SQLStorage) UpdateJob(j job.Job) error {
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
		j.Id,
		j.CronExpString,
		j.Endpoint,
		j.HttpMethod,
		j.MaxRetries,
		j.SuccessStatuses,
		j.Status,
		j.AlertStrategy,
		j.AlertEndpoint,
		j.AlertMethod,
		alertHeadersString,
		alertPayload,
	)

	if err != nil {
		log.Error().Err(err).Msg("could not update job")
		return errors.New("could not update job")
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not commit transaction to update job job")
		return errors.New("could not commit transaction update job")
	}
	return nil
}
