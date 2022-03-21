package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/toolbox/aws"
	"github.com/toolbox/models"
)

// NewDatabasePool creates a database connection pool with the supplied inputs,
// pings it and returns the connection.
// Inputs:
// databaseFamiliarName string   this name will be used in error messages
// databaseDriver       string   name of the pre-registered database driver (e.g.: 'postgres')
// connectionString     string   DSN-like connection string used to connect to the database
// maxOpenConn          int      number of maximum open database connections
// maxIdleConn          int      number of maximum idle database connections
// maxConnLifeTime      duration maximum lifetime of a database connection (on 0 input this wont be set)
// Outputs:
// a *sql.DB database connection pool
// and an optional error occurred during the look up of the database secret
// or during the connection to the database itself.
func NewDatabasePool(
	databaseFamiliarName,
	databaseDriver,
	connectionString string,
	maxOpenConn,
	maxIdleConn int,
	maxConnLifeTime time.Duration) (*sql.DB, error) {

	// Create the database connection with the connection string.
	database, err := sql.Open(databaseDriver, connectionString)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to %s", databaseFamiliarName)
	}

	// Ping the database.
	err = database.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to ping %s", databaseFamiliarName)
	}

	// Set the database pool.
	database.SetMaxOpenConns(maxOpenConn)
	database.SetMaxIdleConns(maxIdleConn)
	if maxConnLifeTime > 0 {
		database.SetConnMaxLifetime(maxConnLifeTime)
	}

	return database, nil
}

// NewDatabasePoolFromDBInfo creates a database connection pool with the supplied inputs,
// pings it and returns the connection.
// Inputs:
// databaseFamiliarName string               this name will be used in error messages
// databaseDriver       string               name of the pre-registered database driver (e.g.: 'postgres')
// dbInfo               *models.DatabaseInfo DatabaseInfo object, used to get the ConnectionString
// maxOpenConn          int                  number of maximum open database connections
// maxIdleConn          int                  number of maximum idle database connections
// maxConnLifeTime      duration             maximum lifetime of a database connection (on 0 input this wont be set)
// Outputs:
// a *sql.DB database connection pool
// and an optional error occurred during the look up of the database secret
// or during the connection to the database itself.
func NewDatabasePoolFromDBInfo(
	databaseFamiliarName,
	databaseDriver string,
	dbInfo *models.DatabaseInfo,
	maxOpenConn,
	maxIdleConn int,
	maxConnLifeTime time.Duration) (*sql.DB, error) {
	if dbInfo == nil {
		return nil, errors.New("Cannot create database connection with nil DatabaseInfo object")
	}
	if err := dbInfo.Validate(); err != nil {
		return nil, errors.Wrap(err, "Invalid DatabaseInfo")
	}
	return NewDatabasePool(
		databaseFamiliarName,
		databaseDriver,
		dbInfo.ConnectionString(),
		maxOpenConn,
		maxIdleConn,
		maxConnLifeTime,
	)
}

// NewDatabasePoolFromSecret creates a database connection pool with the supplied inputs,
// pings it and returns the connection.
// Inputs:
// databaseDriver       string   name of the pre-registered database driver (e.g.: 'postgres')
// databaseFamiliarName string   this name will be used in error messages
// secretName           string   SecretsManager Secret name, used to fetch database info
// maxOpenConn          int      number of maximum open database connections
// maxIdleConn          int      number of maximum idle database connections
// maxConnLifeTime      duration maximum lifetime of a database connection (on 0 input this wont be set)
// Outputs:
// a *sql.DB database connection pool
// and an optional error occurred during the look up of the database secret
// or during the connection to the database itself.
func NewDatabasePoolFromSecret(
	databaseFamiliarName,
	databaseDriver,
	secretName string,
	maxOpenConn,
	maxIdleConn int,
	maxConnLifeTime time.Duration) (*sql.DB, error) {
	secretsManager, err := aws.NewSecretsManager(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect ot SecretsManager")
	}
	connectionString, err := secretsManager.GetConnectionString(secretName)
	if err != nil {
		return nil, errors.New("Failed to get connection string from SecretsManager")
	}
	return NewDatabasePool(
		databaseFamiliarName,
		databaseDriver,
		connectionString,
		maxOpenConn,
		maxIdleConn,
		maxConnLifeTime,
	)
}
