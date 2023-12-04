package scheduler

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
	"github.com/spf13/cobra"

	"github.com/back-end-labs/ruok/pkg/api"
	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/scheduler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func start() {

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
	log.Info().Msgf("scheduler returned with exit code of %d", exitStatus)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info().Msg("Shutting down server")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Server Shutdown error")
		exitStatus = 1
	}
	<-ctx.Done()
	log.Info().Msg("timeout of 5 seconds.")
	log.Info().Msgf("Server exiting with status %d", exitStatus)
	os.Exit(exitStatus)
}

var StartScheduler = &cobra.Command{
	Use:   "start",
	Short: "Starts the scheduler main process",
	Long: `Starts the scheduler main process.
  * It will get the configurations,
  * connect to the database,
  * mount http endpoints and schedule the maximum amount of jobs it can.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}
