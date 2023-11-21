package main

import (
	// add this

	"log"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/scheduler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func main() {
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()
	jobsList := scheduler.NewJobList(int(cfg.MaxJobs))
	if err := scheduler.NewScheduler(s, jobsList).Start(); err != nil {
		// Start draining and then gracefully exit
		log.Fatal("we should log what was lost and unregister from everywhere")
	}

}
