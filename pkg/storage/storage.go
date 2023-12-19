package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/jackc/pgx/v5/pgxpool"
)

// An interface that should help to describe what the Storage is doing
type Storage interface {
	SchedulerStorage
	APIStorage
}

type SchedulerStorage interface {
	ListenForChanges(ch chan uuid.UUID, ctx context.Context)
	StopListeningForChanges() error
	GetJobUpdates(jobId uuid.UUID) *JobUpdates
	GetAvailableJobs(limit int) []*job.Job
	WriteDone(*job.Job) error
	RegisterSelf()
	GetClient() *pgxpool.Pool
	ReleaseAll(j []*job.Job) error
}

type APIStorage interface {
	GetClaimedJobs(limit int, offset int) []*job.Job
	GetClaimedJobsExecutions(jobId uuid.UUID, limit int, offset int) []*job.JobExecution
	Connected() bool
	GetSSLVersion() (bool, string)
	CreateJob(j CreateJobInput) error
	UpdateJob(j UpdateJobInput) error
}

type SQLStorage struct {
	Db *pgxpool.Pool
}

type Closer func()

// Returns the raw client
func (sqls *SQLStorage) GetClient() *pgxpool.Pool {
	return sqls.Db
}

func (sqls *SQLStorage) Connected() bool {
	return sqls.Db.Ping(context.Background()) == nil
}

// Should register the url, name of the application and so on in the db
func (sqls *SQLStorage) RegisterSelf() {}

// It connects to a db
func NewStorage(cfg *config.Configs) (Storage, Closer) {
	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s&application_name=%s",
		cfg.Kind,
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.SSLConfigs.SSLMode,
		cfg.AppName,
	)
	if cfg.SSLConfigs.SSLMode != config.DISABLE_SSL {
		connStr = fmt.Sprintf("%s&sslcert=%s&sslkey=%s&sslrootcert=%s&sslpassword=%s",
			connStr,
			cfg.SSLConfigs.SSLCertPath,
			cfg.SSLConfigs.SSLKeyPath,
			cfg.SSLConfigs.CACertPath,
			cfg.SSLConfigs.SSLPassword)
	}
	// Connect to database
	switch cfg.Kind {

	case "postgres":
		dbconfig, err := pgxpool.ParseConfig(connStr)

		if err != nil {
			log.Fatal().Err(err).Msg("could not parse url to configs")
		}

		dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			pgxuuid.Register(conn.TypeMap())
			return nil
		}

		db, err := pgxpool.NewWithConfig(context.Background(), dbconfig)

		if err != nil {
			log.Fatal().Err(err).Msg("could no stablish a connection with the database, aborting.")
		}
		s := &SQLStorage{Db: db}
		return s, db.Close

	default:
		log.Fatal().Err(errors.New("unrecognized storage")).Msg("Kind field must use one of [ postgres ]")
	}
	return nil, nil
}
