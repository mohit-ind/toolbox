package constants

const (
	// Amazon Web Services EC2 Identifier
	EC2_ID = "EC2_ID"

	// APP_ENV is the applications's environment.
	APP_ENV = "APP_ENV"

	// APP_PORT is the TCP/IP port where the application listens.
	APP_PORT = "APP_PORT"

	// APP_LOG_LEVEL is the level of logging in the application.
	APP_LOG_LEVEL = "APP_LOG_LEVEL"

	// APP_LOG_DEV indicates if Development Mode is on and logs should be rendered as Text and not as JSON.
	APP_LOG_DEV = "APP_LOG_DEV"

	// APP_LOG_FORMAT_ERRORS indicates if Error formating is enabled, so newlines and tabs should be converted.
	APP_LOG_FORMAT_ERRORS = "APP_LOG_FORMAT_ERRORS"

	// APP_DEBUG indicates if the Debug Mode is enabled.
	APP_DEBUG = "APP_DEBUG"

	// APP_DB_SECRET_NAME is the name of the entry in AWS SecretsManager,
	// with the connection info to the application's database.
	APP_DB_SECRET_NAME = "APP_DB_SECRET_NAME"
)

var (
	TruthyValues = []interface{}{"1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE"}
)

var (
	// BasicEnvs is the list of environment variables all service should use.
	BasicEnvs = []string{
		EC2_ID,
		APP_ENV,
		APP_PORT,
		APP_LOG_LEVEL,
		APP_LOG_DEV,
		APP_LOG_FORMAT_ERRORS,
		APP_DEBUG,
		APP_DB_SECRET_NAME,
	}
)

// lib/pq ssl modes https://www.postgresql.org/docs/current/libpq-ssl.html
const (
	// SSL_MODE_DISABLE disables the checking of SSL.
	SSL_MODE_DISABLE = "disable"

	// SSL_MODE_ALLOW allows the checking of SSL.
	SSL_MODE_ALLOW = "allow"

	// SSL_MODE_PREFER prefers  the checking of SSL.
	SSL_MODE_PREFER = "prefer"

	// SSL_MODE_REQUIRE requires  the checking of SSL.
	SSL_MODE_REQUIRE = "require"

	// SSL_MODE_REQUIRE requires  the checking of SSL and verifies it with a CA.
	SSL_MODE_VERIFY_CA = "verify-ca"

	// SSL_MODE_REQUIRE requires  the checking of SSL and verifies it with a CA and validates if the domain exists.
	SSL_MODE_VERIFY_FULL = "verify-full"
)

var (
	// ValidSSLModes are the valid SSL modes. Used in validation.
	ValidSSLModes = []interface{}{
		SSL_MODE_DISABLE,
		SSL_MODE_ALLOW,
		SSL_MODE_PREFER,
		SSL_MODE_REQUIRE,
		SSL_MODE_VERIFY_CA,
		SSL_MODE_VERIFY_FULL,
	}
)

const (
	// LOG_LEVEL_DEBUG - debug level and above
	LOG_LEVEL_DEBUG = "debug"

	// LOG_LEVEL_INFO - info level and above
	LOG_LEVEL_INFO = "info"

	// LOG_LEVEL_WARN - warn level and above
	LOG_LEVEL_WARN = "warn"

	// LOG_LEVEL_ERROR - error level and above
	LOG_LEVEL_ERROR = "error"
)

var (
	// ValidLogLevels are the valid logging levels of the application
	ValidLogLevels = []interface{}{
		LOG_LEVEL_DEBUG,
		LOG_LEVEL_INFO,
		LOG_LEVEL_WARN,
		LOG_LEVEL_ERROR,
	}
)

const (
	// ENV_DEV is the development environment
	ENV_DEV = "dev"

	// ENV_TEST is the testing environment
	ENV_TEST = "test"

	// ENV_STAGING is the staging environment
	ENV_STAGING = "staging"

	// ENV_ACCEPTANCE is the acceptance environment
	ENV_ACCEPTANCE = "acceptance"

	// ENV_PRODUCTION is the production environment
	ENV_PRODUCTION = "production"
)

var (
	// ValidEnvironments are the valid environments of the application
	ValidEnvironments = []interface{}{
		ENV_DEV,
		ENV_TEST,
		ENV_STAGING,
		ENV_ACCEPTANCE,
		ENV_PRODUCTION,
	}
)
