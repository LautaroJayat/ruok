package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
)

func TestReleaseAll(t *testing.T) {
	Drop()
	Seed()
	defer Drop()
	t.Run("Test if release process is working as intended", func(t *testing.T) {
		pendingStatus := "pending to be claimed"
		defer Drop()
		cfg := config.FromEnvs()
		s, close := NewStorage(&cfg)
		defer close()
		joblist := s.GetAvailableJobs(100)
		if joblist == nil {
			t.Error("expected non nil job list")
		}
		if len(joblist) != 10 {
			t.Errorf("expected only 10 jobs, got %d", len(joblist))
		}

		err := s.ReleaseAll(joblist)
		if err != nil {
			t.Error("release process shouldn't produce an error")
		}

		ctx := context.Background()
		tx, err := s.GetClient().Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			t.Errorf("couldn't start transaction for testing. error=%q", err)
		}
		rows, err := tx.Query(ctx, "select id, claimed_by, status from jobs")

		if err != nil {
			t.Errorf("couldn't exec transaction for testing. error=%q", err)
		}

		counter := 0

		for rows.Next() {
			counter++
			var id int64
			var claimedBy sql.NullString
			var status string
			rows.Scan(&id, &claimedBy, &status)
			if claimedBy.String != "" {
				t.Errorf("expected claimed_by to be null, instead got %q", claimedBy.String)

			}
			if status != pendingStatus {
				t.Errorf("expected status to be %q, instead got %q", pendingStatus, status)

			}
		}
		tx.Commit(ctx)

		if counter != len(joblist) {
			t.Errorf("expected released number of jobs equal to previously claimed jobs. Released=%d, previously claimed=%d", counter, len(joblist))
		}
	})
}
