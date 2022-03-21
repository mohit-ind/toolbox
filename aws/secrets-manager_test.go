package aws

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/toolbox/models"
)

// This is the actual name of a secret in appventurez's AWS SecretsManager
// it has invalid username and password and it is only used for test purposes.
const TestSecretName = "devops-integration-test"

// SecretsManagerTestSuite extends testify's Suite.
// It satisfies the SecretsManagerAPI interface, so it can be used as the API of SecretsManager
type SecretsManagerTestSuite struct {
	suite.Suite
	secretValueName  string
	secretValue      *secretsmanager.GetSecretValueOutput
	secretValueError error
}

func (smts *SecretsManagerTestSuite) SetSecretValue(
	name string,
	secret *secretsmanager.GetSecretValueOutput,
	secretError error) {
	smts.secretValueName = name
	smts.secretValue = secret
	smts.secretValueError = secretError
}

func (smts *SecretsManagerTestSuite) SetupTest() {
	smts.secretValueName = ""
	smts.secretValue = nil
	smts.secretValueError = nil
}

func (smts *SecretsManagerTestSuite) GetSecretValue(
	ctx context.Context,
	params *secretsmanager.GetSecretValueInput,
	optFns ...func(*secretsmanager.Options),
) (*secretsmanager.GetSecretValueOutput, error) {
	if smts.secretValueError != nil {
		return nil, smts.secretValueError
	}
	if *params.SecretId == "" {
		//nolint:lll
		return nil, errors.New("operation error Secrets Manager: GetSecretValue, https response error StatusCode: 400, RequestID: c9da13b3-25b6-4d9c-a45f-92666c2331c3, api error ValidationException: Invalid name. Must be a valid name containing alphanumeric characters, or any of the following: -/_+=.@!")
	}
	if *params.SecretId != smts.secretValueName {
		//nolint:lll
		return nil, errors.New("operation error Secrets Manager: GetSecretValue, https response error StatusCode: 400, RequestID: 10ccd2d8-4608-42e1-944c-4fe62eacd7c6, ResourceNotFoundException: Secrets Manager can't find the specified secret.")
	}
	return smts.secretValue, nil
}

func (smts *SecretsManagerTestSuite) TestNewSecretsManager() {
	secretsManager, err := NewSecretsManager(context.TODO())
	smts.NoError(err, "NewSecretsManager should not return any errors")
	smts.NotNil(secretsManager, "New SecretsManager should have been created")
}

func (smts *SecretsManagerTestSuite) TestGetSecret_empty_name() {
	secretsManager := &SecretsManager{
		ctx: context.TODO(),
		api: smts,
	}

	resp, err := secretsManager.GetSecret("")
	smts.Empty(resp, "On error GetSecret should not return any secret")
	smts.Error(err, "An error should be returned when looking up non existing secrets")
	smts.Contains(
		err.Error(),
		"ValidationException: Invalid name.",
		"The error should be ValidationException: Invalid name.")
}

func (smts *SecretsManagerTestSuite) TestGetSecret_wrong_name() {
	secretsManager := &SecretsManager{
		ctx: context.TODO(),
		api: smts,
	}

	resp, err := secretsManager.GetSecret("my-secret")
	smts.Empty(resp, "On error GetSecret should not return any secret")
	smts.Error(err, "An error should be returned when looking up non existing secrets")
	smts.Contains(err.Error(), "ResourceNotFoundException", "The error should be ResourceNotFoundException")
}

func (smts *SecretsManagerTestSuite) TestGetSecret() {
	secretName := "test-secret"
	//nolint:lll
	secretString := "{\"username\":\"db_user\",\"password\":\"db_password\",\"engine\":\"postgres\",\"host\":\"test_db_hostname\",\"port\":5432,\"dbname\":\"db_name\",\"dbInstanceIdentifier\":\"test_database\"}"
	smts.SetSecretValue(secretName, &secretsmanager.GetSecretValueOutput{
		Name:         &secretName,
		SecretString: &secretString,
	}, nil)

	secretsManager := &SecretsManager{
		ctx: context.TODO(),
		api: smts,
	}

	resp, err := secretsManager.GetSecret("test-secret")
	smts.NoError(err, "No error should be returned on successful lookup")
	smts.NotNil(resp, "The map of secrets should be returned")
	dbInfo, err := models.DBInfoFromSecret(resp)
	smts.NoError(err, "A Database Info object should have been created from the secret string, without any error")
	smts.Equal("dbname=db_name user=db_user password=db_password host=test_db_hostname port=5432 sslmode=disable",
		dbInfo.ConnectionString(),
		"Correct DSN like connectins string should be built from the secret string")
}

func (smts *SecretsManagerTestSuite) TestSecretsManager_integration() {
	if ci := os.Getenv("CI"); ci != "true" {
		smts.T().Skip("Skipping SecretsManager integration tests, as we are not in CI.")
	}

	secretsManager, err := NewSecretsManager(context.Background())
	smts.NoError(err, "SecretsManager client should have been created without any errors")

	secretString, err := secretsManager.GetSecret(TestSecretName)
	smts.NoError(err, "The test secret should have been fetched without any error")

	dbInfo, err := models.DBInfoFromSecret(secretString)
	smts.NoError(err, "Database Info should have been built from the secret string without any errors")

	smts.Equal(
		//nolint:lll
		"dbname=devops user=rds-test-database-user password=rds-test-database-password host=devops-test.cifo0surlhze.eu-central-1.rds.amazonaws.com port=5432 sslmode=disable",
		dbInfo.ConnectionString(),
		"The DSN like connection string should have been built from the secret")
}

func (smts *SecretsManagerTestSuite) TestSecretsManager_connection_string_integration() {
	if ci := os.Getenv("CI"); ci != "true" {
		smts.T().Skip("Skipping SecretsManager integration tests, as we are not in CI.")
	}

	secretsManager, err := NewSecretsManager(context.Background())
	smts.NoError(err, "SecretsManager client should have been created without any errors")

	connectionString, err := secretsManager.GetConnectionString(TestSecretName)
	smts.NoError(err, "The secret name should be converted into a valid connection string")

	smts.Equal(
		//nolint:lll
		"dbname=devops user=rds-test-database-user password=rds-test-database-password host=devops-test.cifo0surlhze.eu-central-1.rds.amazonaws.com port=5432 sslmode=disable",
		connectionString,
		"The DSN like connection string should have been built from the secret")
}

func (smts *SecretsManagerTestSuite) TestSecretsManager_integration_error() {
	if ci := os.Getenv("CI"); ci != "true" {
		smts.T().Skip("Skipping SecretsManager integration tests, as we are not in CI.")
	}

	secretsManager, err := NewSecretsManager(context.Background())
	smts.NoError(err, "SecretsManager client should have been created without any errors")

	secretString, err := secretsManager.GetSecret("non-existing-secret-name")
	smts.Error(err, "An error should be returned on non existant secret name")
	smts.Empty(secretString, "On error, the returned connection string should be empty")
}

func (smts *SecretsManagerTestSuite) TestSecretsManager_integration_error_malformed() {
	if ci := os.Getenv("CI"); ci != "true" {
		smts.T().Skip("Skipping SecretsManager integration tests, as we are not in CI.")
	}

	secretsManager, err := NewSecretsManager(context.Background())
	smts.NoError(err, "SecretsManager client should have been created without any errors")

	secretString, err := secretsManager.GetConnectionString("devops-integration-test-empty")
	//nolint:lll
	smts.EqualError(err, "Invalid DatabaseInfo: dbname: cannot be blank; host: cannot be blank; password: cannot be blank; port: cannot be blank; user: cannot be blank.")
	smts.Empty(secretString, "On error, the returned connection string should be empty")
}

// SecretsManager runs the suite
func TestSecretsManager(t *testing.T) {
	suite.Run(t, new(SecretsManagerTestSuite))
}
