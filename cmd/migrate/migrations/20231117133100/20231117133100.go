package m20231117133100

import (
	"context"
	"fmt"

	"github.com/back-end-labs/ruok/pkg/storage"
)

func Migrate20231117133100(client storage.Storage) {
	fmt.Println("20231117133100, Creating micro_unix_now() function")
	ctx := context.Background()
	tx, err := client.GetClient().Begin(ctx)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(ctx, `
	CREATE OR REPLACE FUNCTION micro_unix_now() RETURNS BIGINT AS $$
		BEGIN 
			RETURN (SELECT (EXTRACT(epoch FROM now()) * 1000)::bigint);
		END;
 	$$ LANGUAGE plpgsql;
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
