package jobhandler

import (
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func OnErrorHandler(s storage.SchedulerStorage) func(j *job.Job) {
	// we can hook mor functionalities here if we want
	return func(j *job.Job) {
		s.WriteDone(j)
	}
}
