package jobhandler

import (
	"github.com/back-end-labs/ruok/pkg/alerting"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func OnErrorHandler(s storage.SchedulerStorage, am *alerting.AlertManager) func(j *job.Job) {
	return func(j *job.Job) {
		_, _ = am.SendAlert(j.AlertingInput())
		s.WriteDone(j)
	}
}
