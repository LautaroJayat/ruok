package storage

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimedJobs(t *testing.T) {
	Drop()
	Seed()
	defer Drop()
	t.Run("Test if we are getting the all the claimed jobs as we expect", func(t *testing.T) {
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
		claimedJobs := s.GetClaimedJobs(len(joblist), 0)
		assert.Equal(t, len(claimedJobs), len(joblist))
		expectedIds := []uuid.UUID{}
		for _, j := range joblist {
			expectedIds = append(expectedIds, j.Id)
		}
		for _, j := range claimedJobs {
			assert.Contains(t, expectedIds, j.Id)
		}

	})
}
