package storage

import (
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimedJobExecutions(t *testing.T) {
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

		for _, j := range joblist {
			_ = s.WriteDone(j)
			_ = s.WriteDone(j)
			j.ClaimedBy = "not this app"
			_ = s.WriteDone(j)
		}
		jobExecutions := []*job.JobExecution{}
		for _, j := range joblist {
			jel := s.GetClaimedJobsExecutions(j.Id, 100, 0)
			jobExecutions = append(jobExecutions, jel...)
		}
		assert.Equal(t, len(jobExecutions), len(joblist)*2)
	})
}
