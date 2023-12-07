package config

import (
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// SSL File Names
var CA_CERT_FILE string = "/ca-cert.pem"
var CLIENT_CERT_FILE string = "/client-cert.pem"
var CLIENT_KEY_FILE string = "/client-key.pem"

// SSL Modes
var DISABLE_SSL = "disable"
var REQUIRE_SSL = "require"

// PROD_ENVIRONMENT
var ProdRuokEnvironment = "production"

// ALERT CHANNELS
var ALERT_HTTP string = "http"
var AVAILABLE_CHANNELS = []string{ALERT_HTTP}

//var ALERT_SMS string = "sms"
//var ALERT_MAIL string = "mail"

// EnvNames
var DB_SSLMode string = "DB_SSLMode"
var DB_SSL_PASS string = "DB_SSL_PASS"
var STORAGE_KIND string = "STORAGE_KIND"
var DB_PROTOCOL string = "DB_PROTOCOL"
var DB_PASS string = "DB_PASS"
var DB_USER string = "DB_USER"
var DB_HOST string = "DB_HOST"
var DB_PORT string = "DB_PORT"
var DB_NAME string = "DB_NAME"
var APP_NAME string = "APP_NAME"
var POLL_INTERVAL_SECONDS string = "POLL_INTERVAL_SECONDS"
var MAX_JOBS string = "MAX_JOBS"
var RUOK_ENVIRONMENT = "RUOK_ENVIRONMENT"
var ALERT_CHANNELS = "ALERT_CHANNELS"

// Defaults
var defaultMaxJobs int = 10000
var defaultPollInterval time.Duration = time.Minute
var defaultKind string = "postgres"
var defaultProtocol string = "postgresql"
var defaultPass string = "password"
var defaultUser string = "testing_user"
var defaultHost string = "localhost"
var defaultPort string = "5432"
var defaultDbname string = "db1"
var defaultAppName string = "application1"
var defaultBaseDir string = "/app"
var defaultSSLMode string = DISABLE_SSL
var defaultSSLPass string = "clientpass"
var defaultRuokEnvironment string = "development"
var defaultAlertChannels = []string{ALERT_HTTP}

type Stats struct {
	ClaimedJobs int
	StartedAt   int64
}

func (s *Stats) Uptime() int64 {
	return time.Now().UnixMicro() - s.StartedAt
}

func (s *Stats) CountClaimedJobs() int {
	return s.ClaimedJobs
}

var AppStats *Stats = &Stats{
	StartedAt: time.Now().UnixMicro(),
}

type Configs struct {
	Kind          string
	Protocol      string
	Pass          string
	User          string
	Host          string
	Port          string
	Dbname        string
	SSLConfigs    SSLConfig
	AppName       string
	MaxJobs       int
	PollInterval  time.Duration
	StartedAt     int64
	AlertChannels []string
}

var globalConfigs *Configs = nil

func parseAlertChannels() []string {
	chanString := getEnvOrDefault(ALERT_CHANNELS, ALERT_HTTP)
	inputChannels := strings.Split(chanString, ",")

	for i, ch := range inputChannels {
		inputChannels[i] = strings.Trim(ch, " ")
	}

	channels := []string{}

	for _, ch := range inputChannels {
		for _, availableCh := range AVAILABLE_CHANNELS {
			if ch == availableCh {
				channels = append(channels, ch)
			}
		}
	}

	if len(channels) == 0 {
		return defaultAlertChannels
	}

	return channels

}

func parseMaxJobs(cfg *Configs) {
	maxJobs, err := strconv.ParseInt(os.Getenv(MAX_JOBS), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("could not parse MAX_JOBS env defaulting to %d", defaultMaxJobs)
		globalConfigs.MaxJobs = defaultMaxJobs
	} else {
		globalConfigs.MaxJobs = int(maxJobs)
	}
}

func ParsePollInterval(cfg *Configs) {
	interval, err := strconv.ParseInt(os.Getenv(POLL_INTERVAL_SECONDS), 10, 64)
	if err != nil {
		log.Error().Err(err).Msgf("could not parse POLLING_INTERVAL_SECONDS env defaulting to %f seconds", defaultPollInterval.Seconds())
		globalConfigs.PollInterval = defaultPollInterval
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

// This function assumes "/app" as the dir where the application files will be within the container.
// When developing in local, we usually do not have an "/app" folder. If it doesn't exist, that means
// we are working in the host machine and we should be using this same folder structure.
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

func generateLocalBasePath() string {
	base := ""
	_, currentFile, _, _ := runtime.Caller(0)
	base = path.Dir(currentFile)
	base = path.Join(base, "..", "..", "ssl")
	base = path.Clean(base)
	return base
}

func getSSLConfigs() SSLConfig {
	base := defaultBaseDir

	tlsConfigs := SSLConfig{
		// disable | require
		SSLMode: getEnvOrDefault(DB_SSLMode, defaultSSLMode),
	}
	if tlsConfigs.SSLMode == DISABLE_SSL {
		return tlsConfigs
	}
	// Assuming there wont be a "/app" folder in "/"
	// Just to be able to develop and test outside docker
	if !withinContainer(base) {
		base = generateLocalBasePath()
	}
	log.Debug().Msg(base)
	tlsConfigs.CACertPath = base + CA_CERT_FILE
	tlsConfigs.SSLCertPath = base + CLIENT_CERT_FILE
	tlsConfigs.SSLKeyPath = base + CLIENT_KEY_FILE
	tlsConfigs.SSLPassword = getEnvOrDefault(DB_SSL_PASS, defaultSSLPass)

	return tlsConfigs
}

// Returns TRUE if name is invalid, FALSE if valid.
// Because only letters, numbers and low dashes are allowed
func isInvalidAppName(s string) bool {
	if s == "" {
		return true
	}
	invalid, err := regexp.MatchString(`[^\w|_]`, s)
	if err != nil {
		log.Error().Err(err).Msgf("could not validate app name %q", s)
		return true
	}
	return invalid
}

func validateAppNameOrFail() string {
	appName := getEnvOrDefault(APP_NAME, defaultAppName)
	if isInvalidAppName(appName) {
		log.Fatal().Msgf(
			"Cant continue. Invalid app name. Only letters, numbers and '_' are allowed. Submitted application name was: %q",
			appName,
		)
	}
	return appName
}

func FromEnvs() Configs {
	if globalConfigs == nil {
		globalConfigs = &Configs{
			Kind:          getEnvOrDefault(STORAGE_KIND, defaultKind),
			Protocol:      getEnvOrDefault(DB_PROTOCOL, defaultProtocol),
			Pass:          getEnvOrDefault(DB_PASS, defaultPass),
			User:          getEnvOrDefault(DB_USER, defaultUser),
			Host:          getEnvOrDefault(DB_HOST, defaultHost),
			Port:          getEnvOrDefault(DB_PORT, defaultPort),
			Dbname:        getEnvOrDefault(DB_NAME, defaultDbname),
			AppName:       validateAppNameOrFail(),
			SSLConfigs:    getSSLConfigs(),
			MaxJobs:       defaultMaxJobs,
			PollInterval:  defaultPollInterval,
			AlertChannels: parseAlertChannels(),
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

func AlertChannels() []string {
	if globalConfigs == nil {
		return FromEnvs().AlertChannels
	}
	return globalConfigs.AlertChannels
}
