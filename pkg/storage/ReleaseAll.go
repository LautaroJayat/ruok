package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/back-end-labs/ruok/pkg/job"
)

// Writes an execution result in the db
func (sqls *SQLStorage) ReleaseAll(j []*job.Job) error {
	ctx := context.Background()
	tx, err := sqls.Db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		log.Printf("error=%q\n", err)
		return fmt.Errorf("could not insert into jobs_results. error=%q", err)
	}

	for i := 0; i < len(j); i++ {
		_, err := tx.Exec(ctx, "UPDATE jobs SET claimed_by = NULL, status = $1 WHERE id = $2", "pending to be claimed", j[i].Id)
		if err != nil {
			fmt.Printf("There was a problem while trying to exec release of job . error=%q", err)
			return fmt.Errorf("could not update jobs. error=%q", err)
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		fmt.Printf("There was a problem while trying to commit transaction into jobs. error=%q", err)
		return fmt.Errorf("could not commit transaction into job. error=%q", err)
	}
	return nil
}
