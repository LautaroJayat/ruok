package storage

import (
	"context"
	"fmt"
	"log"

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
	GetAvailableJobs(limit int) []*job.Job
	WriteDone(*job.Job) error
	RegisterSelf()
	GetClient() *pgxpool.Pool
	ReleaseAll(j []*job.Job) error
}

type APIStorage interface {
	GetClaimedJobs(limit int, offset int) []*job.Job
	GetClaimedJobsExecutions(jobId int, limit int, offset int) []*job.JobExecution
}

type SQLStorage struct {
	Db *pgxpool.Pool
}

type Closer func()

// Returns the raw client
func (sqls *SQLStorage) GetClient() *pgxpool.Pool {
	return sqls.Db
}

// Should register the url, name of the application and so on in the db
func (sqls *SQLStorage) RegisterSelf() {}

// It connects to a db
func NewStorage(cfg *config.Configs) (Storage, Closer) {
	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Kind,
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.SSLConfigs.SSLMode,
	)
	if cfg.SSLConfigs.SSLMode != "disable" {
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
		db, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			log.Fatal(err)
		}
		s := &SQLStorage{Db: db}
		return s, db.Close

	default:
		log.Fatalf("error=unrecognized storage %q. Must use one of [ postgres ]", cfg.Kind)
	}
	return nil, nil
}
