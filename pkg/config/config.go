package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

var defaultMaxJobs int = 10000
var defaultPollInterval time.Duration = time.Minute

type Configs struct {
	Kind         string
	Protocol     string
	Pass         string
	User         string
	Host         string
	Port         string
	Dbname       string
	SSLConfigs   SSLConfig
	AppName      string
	MaxJobs      int
	PollInterval time.Duration
}

var globalConfigs *Configs

func parseMaxJobs(cfg *Configs) {
	maxJobs, err := strconv.ParseInt(os.Getenv("MAX_JOBS"), 10, 64)
	if err != nil {
		fmt.Printf("could not parse MAX_JOBS env defaulting to . error=%q\n", err.Error())
	} else {
		globalConfigs.MaxJobs = int(maxJobs)
	}
}

func ParsePollInterval(cfg *Configs) {
	interval, err := strconv.ParseInt(os.Getenv("POLL_INTERVAL_SECONDS"), 10, 64)
	if err != nil {
		fmt.Printf("could not parse POLLING_INTERVAL_SECONDS env defaulting to 60. error=%q\n", err.Error())
	} else {
		globalConfigs.PollInterval = time.Second * time.Duration(interval)
	}
}

func getEnvOrDefault(env string, defaultValue string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return defaultValue

}

func withinContainer(base string) bool {
	_, err := os.ReadDir(base)
	return err == nil
}

type SSLConfig struct {
	SSLMode     string
	CACertPath  string
	SSLCertPath string
	SSLKeyPath  string
	SSLPassword string
}

func getSSLConfigs(baseDir string) SSLConfig {
	base := baseDir

	tlsConfigs := SSLConfig{
		// disable | require
		SSLMode: getEnvOrDefault("DB_SSLMode", "disable"),
	}
	if tlsConfigs.SSLMode == "disable" {
		return tlsConfigs
	}
	// Assuming there wont be a "/app" folder in "/"
	// Just to be able to develop and test outside docker
	if !withinContainer(base) {
		_, currentFile, _, _ := runtime.Caller(0)
		base = path.Dir(currentFile)
		base = path.Join(base, "..", "..", "ssl")
		base = path.Clean(base)
	}
	tlsConfigs.CACertPath = base + "/ca-cert.pem"
	tlsConfigs.SSLCertPath = base + "/client-cert.pem"
	tlsConfigs.SSLKeyPath = base + "/client-key.pem"
	tlsConfigs.SSLPassword = getEnvOrDefault("DB_SSL_PASS", "clientpass")

	return tlsConfigs
}

func FromEnvs() Configs {
	if globalConfigs == nil {
		baseDir := "/app"
		globalConfigs = &Configs{
			Kind:         getEnvOrDefault("STORAGE_KIND", "postgres"),
			Protocol:     getEnvOrDefault("DB_PROTOCOL", "postgresql"),
			Pass:         getEnvOrDefault("DB_PASS", "password"),
			User:         getEnvOrDefault("DB_USER", "user"),
			Host:         getEnvOrDefault("DB_HOST", "localhost"),
			Port:         getEnvOrDefault("DB_PORT", "5432"),
			Dbname:       getEnvOrDefault("DB_NAME", "db1"),
			SSLConfigs:   getSSLConfigs(baseDir),
			AppName:      getEnvOrDefault("APP_NAME", "application1"),
			MaxJobs:      defaultMaxJobs,
			PollInterval: defaultPollInterval,
		}
	}
	return *globalConfigs
}

func MaxJobs() int {
	if globalConfigs == nil {
		return FromEnvs().MaxJobs
	}
	return globalConfigs.MaxJobs
}

func AppName() string {
	if globalConfigs == nil {
		return FromEnvs().AppName
	}
	return globalConfigs.AppName
}

func PollingInterval() time.Duration {
	if globalConfigs == nil {
		return FromEnvs().PollInterval
	}
	return globalConfigs.PollInterval
}
