package storage

import (
	"context"
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
)

func TestGetJobsQuery(t *testing.T) {
	Drop()
	Seed()
	t.Run("Test if we are getting the jobs as we expect", func(t *testing.T) {
		claimedStatus := "claimed"
		appName := config.AppName()

		defer Drop()
		cfg := config.FromEnvs()
		s, close := NewStorage(&cfg)
		defer close()
		joblist := s.GetAvailableJobs(100)
		if joblist == nil {
			t.Error("expected non nil job list")
		}
		if len(joblist) != 10 {
			t.Errorf("expected 10 jobs, got %d", len(joblist))
		}
		ctx := context.Background()
		tx, err := s.GetClient().Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			t.Errorf("couldn't start transaction for testing. error=%q", err)
		}
		rows, err := tx.Query(ctx, "select claimed_by, status from jobs")

		if err != nil {
			t.Errorf("couldn't exec transaction for testing. error=%q", err)
		}

		for rows.Next() {
			var claimedBy, claimed string
			rows.Scan(&claimedBy, &claimed)
			if claimedBy != appName {
				t.Errorf("expected claimed_by to be %q, instead got %q", appName, claimedBy)
				break
			}
			if claimed != claimedStatus {
				t.Errorf("expected status to be %q, instead got %q", claimedStatus, claimed)
				break
			}
		}
		tx.Commit(ctx)
	})
}
