package storage

import (
	"context"
	"database/sql"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/rs/zerolog/log"
)

// Writes an execution result in the db
func (sqls *SQLStorage) GetSSLVersion() (bool, string) {
	var sslActive bool
	var sslVersion sql.NullString
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)

	if err != nil {
		log.Error().Err(err).Msg("could not start transaction to get own connection ssl version")
		return sslActive, sslVersion.String
	}

	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, "SELECT ssl_active, ssl_version from get_ssl_conn_version($1)", config.AppName())

	row.Scan(&sslActive, &sslVersion)

	err = tx.Commit(ctx)

	if err != nil {
		log.Error().Err(err).Msg("There was a problem while trying to commit 'get SSL version' transaction")
	}

	return sslActive, sslVersion.String
}
