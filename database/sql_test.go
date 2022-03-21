package database

import (
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // blank import required for postgres database driver and connection
	"github.com/stretchr/testify/require"
	"github.com/toolbox/constants"
	"github.com/toolbox/models"
)

// This is the actual name of the test secret in our SecretsManager
const secretName = "devops-integration-docker-test"

func TestNewDatabasePool(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "true" {
		t.Skip("Skipping Database Pool test as not running in CI")
	}
	assert := require.New(t)

	db, err := NewDatabasePool("Test Database", constants.DefaultDatabaseDriver, os.Getenv("TEST_DB_CONNECTION_STRING"), 1, 1, time.Second)
	assert.NoError(err, "The database connection pool should have been created without any error")
	defer db.Close()
	assert.NotNil(db, "The database connection pool object should have been returned")
}

func TestNewDatabasePoolFromDBInfo(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "true" {
		t.Skip("Skipping Database Pool test as not running in CI")
	}
	assert := require.New(t)
	db, err := NewDatabasePoolFromDBInfo("Test Database", constants.DefaultDatabaseDriver, &models.DatabaseInfo{
		Database: "postgres",
		Username: "postgres",
		Password: "pass123",
		Host:     "localhost",
		Port:     "5432",
		SSLMode:  "disable",
	}, 1, 1, time.Second)
	assert.NoError(err, "The database connection pool should have been created without any error")
	defer db.Close()
	assert.NotNil(db, "The database connection pool object should have been returned")
}

func TestNewDatabasePoolFromSecret(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "true" {
		t.Skip("Skipping Database Pool test as not running in CI")
	}
	assert := require.New(t)

	db, err := NewDatabasePoolFromSecret("Test Database", constants.DefaultDatabaseDriver, secretName, 1, 1, time.Second)
	assert.NoError(err, "The database connection pool should have been created without any error")
	defer db.Close()
	assert.NotNil(db, "The database connection pool object should have been returned")
}
