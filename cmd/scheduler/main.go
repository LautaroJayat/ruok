package main

import (
	// add this

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/back-end-labs/ruok/pkg/api"
	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/scheduler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func main() {

	cfg := config.FromEnvs()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Info().Msg("Starting RUOK scheduler :)")

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
			log.Error().Err(err).Msg("error while HTTP server was listening")
		}
	}()

	exitStatus := scheduler.NewScheduler(s, jobsList).Start(signalCh)
	log.Info().Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown error")
		exitStatus = 1
	}
	<-ctx.Done()
	log.Info().Msg("timeout of 5 seconds.")
	log.Info().Msg("Server exiting")
	os.Exit(exitStatus)
}
