package main

import (
	// add this

	"os"
	"os/signal"
	"syscall"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/scheduler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func main() {
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()
	jobsList := scheduler.NewJobList(int(cfg.MaxJobs))
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM /*syscall.SIGHUP*/)
	exitStatus := scheduler.NewScheduler(s, jobsList).Start(signalCh)
	os.Exit(exitStatus)

}
