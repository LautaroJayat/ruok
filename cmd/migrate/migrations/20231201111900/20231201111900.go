package m20231201111900

import (
	"context"
	"fmt"

	"github.com/back-end-labs/ruok/pkg/storage"
)

func Migrate20231201111900(client storage.Storage) {
	fmt.Println("20231201111900, Generating Function to query own SSL version")
	ctx := context.Background()
	tx, err := client.GetClient().Begin(ctx)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(ctx, `
	CREATE OR REPLACE FUNCTION get_ssl_conn_version(app_name text)
	RETURNS TABLE(ssl_active boolean, ssl_version text)
	AS $$
		SELECT ssl,version FROM pg_stat_ssl
		JOIN pg_stat_activity
		ON pg_stat_ssl.pid = pg_stat_activity.pid
		WHERE application_name = $1
		LIMIT 1
	$$ LANGUAGE SQL;
	`)
	if err != nil {
		fmt.Println("error", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("Migration OK")
}
