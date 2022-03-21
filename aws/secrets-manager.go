package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
	models "github.com/toolboxmodels"
)

// SecretsManagerAPI is the interface of object which can GetSecretValue
type SecretsManagerAPI interface {
	GetSecretValue(
		ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManager is a client to AWS SecretsManager, for retrieve secret configuration items.
type SecretsManager struct {
	ctx context.Context
	api SecretsManagerAPI
}

// NewSecretsManager creates a new SecretsManager in the provided context.
// It may return on optional error if the underlying AWS Client failed to initialize.
func NewSecretsManager(ctx context.Context) (*SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create AWS Client")
	}
	return &SecretsManager{
		ctx: ctx,
		api: secretsmanager.NewFromConfig(cfg),
	}, nil
}

// GetSecret looks up AWS SecretsManager for the named secret and returns it as a string.
// The returned sting is the JSON representation of the secret-map.
// GetSecret may return an error if it fails to look up AWS.
func (sm *SecretsManager) GetSecret(secretName string) (string, error) {
	secretsValues, err := sm.api.GetSecretValue(sm.ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to retrieve secret")
	}
	return *secretsValues.SecretString, nil
}

// GetConnectionString gets the secret named databaseSecretName and builds a DSN like connection string with it.
// It may return an error if it fails to look up the secret in AWS, cannot convert it to DatabaseInfo or
// the info itself is incorrect.
func (sm *SecretsManager) GetConnectionString(databaseSecretName string) (string, error) {
	secretMap, err := sm.GetSecret(databaseSecretName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get connection string")
	}

	dbInfo, err := models.DBInfoFromSecret(secretMap)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create DatabaseInfo from secret")
	}

	if err := dbInfo.Validate(); err != nil {
		return "", errors.Wrap(err, "Invalid DatabaseInfo")
	}

	return dbInfo.ConnectionString(), nil
}
