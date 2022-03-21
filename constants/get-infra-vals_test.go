package constants

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBaseURL(t *testing.T) {
	req := require.New(t)

	req.NoError(os.Setenv("ENV", ENV_DEV), "The ENV env-var should have been set to 'dev'")
	url := GetBaseURL()
	req.Equal("localhost", url, "The Dev environment's URL should have been returned")

	req.NoError(os.Setenv("ENV", ENV_TEST), "The ENV env-var should have been set to 'test'")
	url = GetBaseURL()
	req.Equal(BaseURLTest, url, "The Test environment's URL should have been returned")

	req.NoError(os.Setenv("ENV", ENV_STAGING), "The ENV env-var should have been set to 'staging'")
	url = GetBaseURL()
	req.Equal(BaseURLStaging, url, "The Staging environment's URL should have been returned")

	req.NoError(os.Setenv("ENV", "old-staging"), "The ENV env-var should have been set to 'old-staging'")
	url = GetBaseURL()
	req.Equal(BaseURLOldStaging, url, "The Old Staging environment's URL should have been returned")

	req.NoError(os.Setenv("ENV", ENV_ACCEPTANCE), "The ENV env-var should have been set to 'acceptance'")
	url = GetBaseURL()
	req.Equal(BaseURLAcceptance, url, "The Acceptance environment's URL should have been returned")

	req.NoError(os.Setenv("ENV", ENV_PRODUCTION), "The ENV env-var should have been set to 'production'")
	url = GetBaseURL()
	req.Equal(BaseURLProduction, url, "In production environment the base Production URL should be returned")

	req.NoError(os.Unsetenv("ENV"), "The ENV env-var should have been unset")
	url = GetBaseURL()
	req.Equal("localhost", url, "In case if ENV is not set 'localhost' should be returned")

	req.NoError(
		os.Setenv("ENV", "invalid-environment"),
		"The ENV env-var should have been set to 'invalid-environment'",
	)
	url = GetBaseURL()
	req.Equal("localhost", url, "In case if ENV is not a valid Environment, 'localhost' should be returned")

	req.NoError(os.Setenv("ENV", ENV_PRODUCTION), "The ENV env-var should have been set to 'production'")
	url = GetBaseURL()
	req.Equal(BaseURLProduction, url, "In production environment the base Production URL should be returned")
}
