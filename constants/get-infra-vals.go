package constants

import (
	"os"
)

// GetBaseURLbyEnv returns the API Gateway's external URL based on the supplied environment
func GetBaseURLbyEnv(env string) string {
	switch env {
	case ENV_TEST:
		return BaseURLTest
	case ENV_STAGING:
		return BaseURLStaging
	case "old-staging":
		return BaseURLOldStaging
	case ENV_ACCEPTANCE:
		return BaseURLAcceptance
	case ENV_PRODUCTION:
		return BaseURLProduction
	}
	return "localhost"
}

// GetBaseURL calls GetBaseURLbyEnv with the actual value of the ENV environment variable
func GetBaseURL() string {
	return GetBaseURLbyEnv(os.Getenv("ENV"))
}
