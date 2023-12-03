package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromEnvsReturnsACopy(t *testing.T) {
	globalConfigs = nil
	globals := &globalConfigs
	generated := FromEnvs()
	assert.Equal(t, **globals, generated)
	generated.AppName = "x"
	generated.Dbname = "x"
	generated.Host = "x"
	generated.Port = "x"
	generated.Pass = "x"
	generated.User = "x"
	assert.NotEqual(t, **globals, generated)
	notAffectedCOnfigs := FromEnvs()
	assert.Equal(t, **globals, notAffectedCOnfigs)
	assert.NotEqual(t, notAffectedCOnfigs, generated)
}

func TestGetEnvOrDefault(t *testing.T) {
	TESTING_ENV := "TESTING_ENV"
	testingValue := "testing value"
	defaultValue := "defaultValue"
	os.Setenv(TESTING_ENV, testingValue)
	assert.Equal(t, getEnvOrDefault(TESTING_ENV, defaultValue), testingValue)
	os.Unsetenv(TESTING_ENV)
	assert.Equal(t, getEnvOrDefault(TESTING_ENV, defaultValue), defaultValue)
}

func TestWithinContainer(t *testing.T) {
	// because it should exist
	base := os.TempDir()
	assert.True(t, withinContainer(base))
	assert.False(t, withinContainer(base+"PROBABLY WONT EXIST FOLDER"))
}

func TestGetSSLConfigs_NO_SSL(t *testing.T) {
	SSL_FROM_ENV := os.Getenv(DB_SSLMode)
	os.Unsetenv(DB_SSLMode)
	os.Setenv(DB_SSLMode, DISABLE_SSL)
	defer func() {
		os.Unsetenv(DB_SSLMode)
		os.Setenv(DB_SSLMode, SSL_FROM_ENV)
	}()
	tlsConfigs := getSSLConfigs()
	assert.Equal(t, tlsConfigs.SSLMode, DISABLE_SSL)
	assert.Equal(t, tlsConfigs.SSLCertPath, "")
	assert.Equal(t, tlsConfigs.CACertPath, "")
	assert.Equal(t, tlsConfigs.SSLPassword, "")
	assert.Equal(t, tlsConfigs.SSLKeyPath, "")

}

func TestGetSSLConfigs_W_SSL(t *testing.T) {
	SSL_FROM_ENV := os.Getenv(DB_SSLMode)
	os.Unsetenv(DB_SSLMode)
	os.Setenv(DB_SSLMode, REQUIRE_SSL)

	SSL_PASSWORD_FROM_ENV := os.Getenv(DB_SSL_PASS)
	os.Unsetenv(DB_SSL_PASS)
	os.Setenv(DB_SSL_PASS, defaultPass)

	defer func() {
		os.Unsetenv(DB_SSLMode)
		os.Setenv(DB_SSLMode, SSL_FROM_ENV)
		os.Unsetenv(DB_SSL_PASS)
		os.Setenv(DB_SSL_PASS, SSL_PASSWORD_FROM_ENV)
	}()

	basePath := generateLocalBasePath()
	tlsConfigs := getSSLConfigs()

	assert.Equal(t, tlsConfigs.SSLMode, REQUIRE_SSL)
	assert.Equal(t, tlsConfigs.SSLCertPath, basePath+CLIENT_CERT_FILE)
	assert.Equal(t, tlsConfigs.CACertPath, basePath+CA_CERT_FILE)
	assert.Equal(t, tlsConfigs.SSLKeyPath, basePath+CLIENT_KEY_FILE)
	assert.Equal(t, tlsConfigs.SSLPassword, defaultPass)

}

func TestFromEnvs(t *testing.T) {
	globalConfigs = nil
	originalStorageKind := os.Getenv(STORAGE_KIND)
	originalDBProtocol := os.Getenv(DB_PROTOCOL)
	originalDBPass := os.Getenv(DB_PASS)
	originalDBUser := os.Getenv(DB_USER)
	originalDBHost := os.Getenv(DB_HOST)
	originalDBPort := os.Getenv(DB_PORT)
	originalDBName := os.Getenv(DB_NAME)
	originalAppName := os.Getenv(APP_NAME)
	originalDBSSLMode := os.Getenv(DB_SSLMode)
	originalDBSSLPass := os.Getenv(DB_SSL_PASS)
	originalMaxJobs := os.Getenv(MAX_JOBS)
	originalPollInterval := os.Getenv(POLL_INTERVAL_SECONDS)

	os.Setenv(STORAGE_KIND, "")
	os.Setenv(DB_PROTOCOL, "")
	os.Setenv(DB_PASS, "")
	os.Setenv(DB_USER, "")
	os.Setenv(DB_HOST, "")
	os.Setenv(DB_PORT, "")
	os.Setenv(DB_NAME, "")
	os.Setenv(APP_NAME, "")
	os.Setenv(DB_SSLMode, REQUIRE_SSL)
	os.Setenv(DB_SSL_PASS, defaultSSLPass)
	os.Setenv(MAX_JOBS, "")
	os.Setenv(POLL_INTERVAL_SECONDS, "")

	// Clean up environment variables after the test
	defer func() {
		os.Setenv(STORAGE_KIND, originalStorageKind)
		os.Setenv(DB_PROTOCOL, originalDBProtocol)
		os.Setenv(DB_PASS, originalDBPass)
		os.Setenv(DB_USER, originalDBUser)
		os.Setenv(DB_HOST, originalDBHost)
		os.Setenv(DB_PORT, originalDBPort)
		os.Setenv(DB_NAME, originalDBName)
		os.Setenv(APP_NAME, originalAppName)
		os.Setenv(DB_SSLMode, originalDBSSLMode)
		os.Setenv(DB_SSL_PASS, originalDBSSLPass)
		os.Setenv(MAX_JOBS, originalMaxJobs)
		os.Setenv(POLL_INTERVAL_SECONDS, originalPollInterval)
	}()

	// Call the function to be tested
	config := FromEnvs()

	// Check if the values are correctly loaded from the environment variables
	if config.Kind != defaultKind {
		t.Errorf("Expected Kind to be '%s', but got '%s'", defaultKind, config.Kind)
	}

	if config.Protocol != defaultProtocol {
		t.Errorf("Expected Protocol to be '%s', but got '%s'", defaultProtocol, config.Protocol)
	}

	if config.Pass != defaultPass {
		t.Errorf("Expected Pass to be '%s', but got '%s'", defaultPass, config.Pass)
	}

	if config.User != defaultUser {
		t.Errorf("Expected User to be '%s', but got '%s'", defaultUser, config.User)
	}

	if config.Host != defaultHost {
		t.Errorf("Expected Host to be '%s', but got '%s'", defaultHost, config.Host)
	}

	if config.Port != defaultPort {
		t.Errorf("Expected Port to be '%s', but got '%s'", defaultPort, config.Port)
	}

	if config.Dbname != defaultDbname {
		t.Errorf("Expected Dbname to be '%s', but got '%s'", defaultDbname, config.Dbname)
	}

	if config.AppName != defaultAppName {
		t.Errorf("Expected AppName to be '%s', but got '%s'", defaultAppName, config.AppName)
	}

	if config.SSLConfigs.SSLMode != REQUIRE_SSL {
		t.Errorf("Expected SSLMode to be '%s', but got '%s'", REQUIRE_SSL, config.SSLConfigs.SSLMode)
	}

	if config.SSLConfigs.SSLPassword != defaultSSLPass {
		t.Errorf("Expected SSLPassword to be '%s', but got '%s'", defaultSSLPass, config.SSLConfigs.SSLPassword)
	}

	if config.MaxJobs != defaultMaxJobs {
		t.Errorf("Expected MaxJobs to be %d, but got '%d'", defaultMaxJobs, config.MaxJobs)
	}

	if config.PollInterval != defaultPollInterval {
		t.Errorf("Expected PollInterval to be %s, but got '%s'", defaultPollInterval, config.PollInterval)
	}
}

func TestIsInvalidAppName(t *testing.T) {
	testCases := []struct {
		name    string
		invalid bool
	}{
		{"validAppName", false},
		{"invalid@AppName", true},
		{"invalid App Name", true},
		{"", true},
		{"valid_App_Name_123", false},
		{"normalName123", false},
		{"nameWithSpaces 123", true},
		{"name_With_Underscores", false},
		{"name-With-Hyphens", true},
		{"name%With%Special%Characters", true},
		{"nameWithSingleQuote'", true},
		{"nameWithDoubleQuote\"", true},
		{"nameWithBacktick`", true},
		{"nameWithParenthesis(", true},
		{"nameWithParenthesis)", true},
		{"nameWithAmpersand&", true},
		{"nameWithEquals=", true},
		{"nameWithSemicolon;", true},
		{"'; DROP TABLE users; --", true},
		{"OR 1=1 --", true},
		{"UNION SELECT * FROM users --", true},
		{"' OR 'a'='a' --", true},
		{"\" OR \"a\"=\"a\" --", true},
		{"` OR `a`=`a` --", true},
		{"; DROP TABLE users; --", true},
		{"'; SHUTDOWN --", true},
		{"\"; SHUTDOWN --", true},
		{"'; SELECT * FROM users; --", true},
		{"' UNION SELECT password FROM users; --", true},
		{"'; EXEC xp_cmdshell('ls'); --", true},
		{"'; EXEC xp_cmdshell('cat /etc/passwd'); --", true},
		{"'; EXEC master..xp_cmdshell('ls'); --", true},
		{"'; EXEC('DROP TABLE users'); --", true},
		{"' OR 1=1; DROP TABLE users; --", true},
		{"' OR 'a'='a'; DROP TABLE users; --", true},
		{"'; UPDATE users SET password='hacked'; --", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isInvalid := isInvalidAppName(tc.name)
			assert.Equal(t, tc.invalid, isInvalid, "Unexpected result for app name: %s", tc.name)
		})
	}
}
