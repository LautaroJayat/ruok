package m20231117135301

import (
	"context"
	"fmt"

	"github.com/back-end-labs/ruok/pkg/storage"
)

func Migrate20231117133101(s storage.Storage) {
	fmt.Println("20231117133101, Generating Jobs Registration Table DB")
	ctx := context.Background()
	tx, err := s.GetClient().Begin(ctx)
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS jobs (
		id bigserial PRIMARY KEY,
		cron_exp_string varchar,
		endpoint varchar,
		httpmethod varchar,
		max_retries smallint,
		last_execution bigint,
		should_execute_at bigint,
		last_response_at bigint,
		last_message varchar,
		last_status_code int,
		headers_string varchar,
		success_statuses int[],
		tls_client_cert varchar,
		status varchar,
		claimed_by varchar
	  );`)
	if err != nil {
		fmt.Println("error", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("Migration OK")
}
