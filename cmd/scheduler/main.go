package main

import (
	// add this

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/back-end-labs/ruok/pkg/api"
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
	srv := &http.Server{
		Addr:    ":8080",
		Handler: api.CreateRouter(s),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error while api was listening: %q\n", err.Error())
		}
	}()

	exitStatus := scheduler.NewScheduler(s, jobsList).Start(signalCh)
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown error:", err.Error())
		exitStatus = 1
	}
	<-ctx.Done()
	fmt.Println("timeout of 5 seconds.")
	fmt.Println("Server exiting")
	os.Exit(exitStatus)
}
