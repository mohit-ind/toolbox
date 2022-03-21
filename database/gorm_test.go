package database

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	logger "github.com/toolboxlogger"
)

func TestNewMSSQLGormDB(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "true" {
		t.Skip("Skipping NewMSSQLGormDB test as not running in CI")
	}
	assert := require.New(t)

	db, err := NewMSSQLGormDB(
		logger.NewLogger(logrus.New(), nil),
		nil,
		"sqlserver://sa:Pass1234@localhost:1433?database=master",
		1,
		1,
		time.Minute)
	assert.NoError(err, "The database connection pool should have been created without any error")
	assert.NotNil(db, "The database connection pool object should have been returned")
}
