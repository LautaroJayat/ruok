package m20231117172601

import (
	"context"
	"fmt"

	"github.com/back-end-labs/ruok/pkg/storage"
)

func Migrate20231117172601(client storage.Storage) {
	fmt.Println("20231117172601, Generating Jobs Results Table")
	ctx := context.Background()
	tx, err := client.GetClient().Begin(ctx)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS public.job_results (
		id bigserial,
		job_id bigint,
		cron_exp_string varchar,
		endpoint varchar,
		httpmethod varchar,
		max_retries smallint,
		execution_time bigint,
		should_execute_at bigint,
		last_response_at bigint,
		last_message varchar,
		last_status_code int,
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
