package main

import (
	// add this

	"os"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/scheduler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func main() {
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()
	jobsList := scheduler.NewJobList(int(cfg.MaxJobs))
	exitStatus := scheduler.NewScheduler(s, jobsList).Start()
	os.Exit(exitStatus)

}
