package storage

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/job"
)

// Writes an execution result in the db
func (sqls *SQLStorage) ReleaseAll(j []*job.Job) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to release all jobs")
		return errors.New("could not insert into jobs_results")
	}
	defer tx.Rollback(ctx)

	for i := 0; i < len(j); i++ {
		_, err := tx.Exec(ctx, "UPDATE ruok.jobs SET claimed_by = NULL, status = $1 WHERE id = $2", "pending to be claimed", j[i].Id)
		if err != nil {
			log.Error().Err(err).Msgf("There was a problem while trying to exec release of job with id of %v", j[i].Id)
			return errors.New("could not update jobs")
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("There was a problem while trying to commit 'release all jobs' transaction")
		return errors.New("could not commit transaction")
	}
	return nil
}
