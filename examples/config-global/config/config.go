package config

import (
	"sync"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"

	toolboxConfig "github.com/toolboxconfig"
	constants "github.com/toolboxconstants"
)

const envfile = ".env"

var (
	global *toolboxConfig.AppConfig

	setupOnce sync.Once

	defaultsConfigs = map[string]*toolboxConfig.Variable{
		constants.APP_PORT: {
			DefaultValue: "8080",
			Description:  "TCP/IP Port where the application listens",
			Rules: map[string]validation.Rule{
				"Required":   validation.Required,
				"Valid port": is.Port,
			},
		},
		constants.APP_ENV: {
			DefaultValue: constants.ENV_TEST,
			Description:  "The environment of the application",
			Rules: map[string]validation.Rule{
				"Required":          validation.Required,
				"Valid environment": validation.In(constants.ValidEnvironments...),
			},
		},
		constants.APP_DEBUG: {
			DefaultValue: "true",
			Description:  "Debug mode",
			Rules: map[string]validation.Rule{
				"Truthy value": validation.In(constants.TruthyValues...),
			},
		},
		constants.APP_LOG_LEVEL: {
			DefaultValue: constants.LOG_LEVEL_DEBUG,
			Description:  "Level of logging",
			Rules: map[string]validation.Rule{
				"Required":        validation.Required,
				"Valid log level": validation.In(constants.ValidLogLevels...),
			},
		},
		constants.APP_LOG_DEV: {
			Description: "Log development mode (Test formatter instead of JSON)",
			Rules: map[string]validation.Rule{
				"Truthy value": validation.In(constants.TruthyValues...),
			},
		},
		constants.APP_LOG_FORMAT_ERRORS: {
			Description: "Format error log entries by switching newlines to --- and tabs to spaces",
			Rules: map[string]validation.Rule{
				"Truthy value": validation.In(constants.TruthyValues...),
			},
		},
		constants.APP_DB_SECRET_NAME: {
			Description: "The Database's secret's name in AWS SecretsManager",
		},
	}
)

func SetupConfigs() error {
	var err error
	setupOnce.Do(func() {
		global = toolboxConfig.NewConfig(defaultsConfigs)
		err = global.Setup(envfile)
	})
	if err != nil {
		if val, ok := err.(validation.Errors); ok {
			err = val.Filter()
		}
		err = errors.Wrap(err, "Failed to setup configs")
	}
	return err
}

func Global() *toolboxConfig.AppConfig {
	return global
}
