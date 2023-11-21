package config

import (
	"fmt"
	"os"
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

func parsePollInterval(cfg *Configs) {
	interval, err := strconv.ParseInt(os.Getenv("POLL_INTERVAL_SECONDS"), 10, 64)
	if err != nil {
		fmt.Printf("could not parse POLLING_INTERVAL_SECONDS env defaulting to 60. error=%q\n", err.Error())
	} else {
		globalConfigs.PollInterval = time.Second * time.Duration(interval)
	}
}

func FromEnvs() Configs {
	if globalConfigs == nil {
		globalConfigs = &Configs{
			Kind:         os.Getenv("STORAGE_KIND"),
			Protocol:     os.Getenv("DB_PROTOCOL"),
			Pass:         os.Getenv("DB_PASS"),
			User:         os.Getenv("DB_USER"),
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			Dbname:       os.Getenv("DB_NAME"),
			AppName:      os.Getenv("APP_NAME"),
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
