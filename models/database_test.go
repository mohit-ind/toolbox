package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatabaseInfo(t *testing.T) {
	assert := require.New(t)

	dbInfo := NewDatabaseInfo("test-db", "test-user", "test-password", "testhost", "1212", "require")
	assert.NoError(dbInfo.Validate(), "This should be a valid database connection info")
	assert.Equal("test-db", dbInfo.Database, "The database name should be test-db")
	assert.Equal("test-user", dbInfo.Username, "The username should be test-user")
	assert.Equal("test-password", dbInfo.Password, "The password should be test-password")
	assert.Equal("testhost", dbInfo.Host, "The host should be testhost")
	assert.Equal("1212", dbInfo.Port, "The port should be 1212")
	assert.Equal("require", dbInfo.SSLMode, "The sslmode should be require")
	assert.Equal("dbname=test-db user=test-user password=test-password host=testhost port=1212 sslmode=require",
		dbInfo.ConnectionString())
}

func TestNewDatabaseInfoFromEnvs(t *testing.T) {
	assert := require.New(t)

	assert.NoError(os.Setenv("DBNAME", "test-db"), "Database name should have been set")
	assert.NoError(os.Setenv("DBUSER", "test-user"), "Database username should have been set")
	assert.NoError(os.Setenv("DBPASS", "test-password"), "Database Password should have been set")
	assert.NoError(os.Setenv("DBHOST", "testhost"), "Database Host should have been set")
	assert.NoError(os.Setenv("DBPORT", "9393"), "Database Port should have been set")

	dbInfo := DBInfoFromEnvs("DBNAME", "DBUSER", "DBPASS", "DBHOST", "DBPORT", "an-env-that-does-not-exists")
	assert.NoError(dbInfo.Validate(), "This should be a valid database connection info")
	assert.Equal("test-db", dbInfo.Database, "The database name should be test-db")
	assert.Equal("test-user", dbInfo.Username, "The username should be test-user")
	assert.Equal("test-password", dbInfo.Password, "The password should be test-password")
	assert.Equal("testhost", dbInfo.Host, "The host should be testhost")
	assert.Equal("9393", dbInfo.Port, "The port should be 9393")
	assert.Equal("disable", dbInfo.SSLMode, "The sslmode should be disable")
	assert.Equal("dbname=test-db user=test-user password=test-password host=testhost port=9393 sslmode=disable",
		dbInfo.ConnectionString())
}

func TestDBInfoFromSecret_malformed_json(t *testing.T) {
	assert := require.New(t)
	//nolint:lll
	secretString := "🐸"

	dbInfo, err := DBInfoFromSecret(secretString)
	assert.Error(err, "An error should have been returned")
	assert.Nil(dbInfo, "On malformed secret string no Database Info should be returned")
}

func TestDBInfoFromSecret(t *testing.T) {
	assert := require.New(t)
	//nolint:lll
	secretString := "{\"username\":\"db_user\",\"password\":\"db_password\",\"engine\":\"postgres\",\"host\":\"test_db_hostname\",\"port\":5432,\"dbname\":\"db_name\",\"dbInstanceIdentifier\":\"test_database\"}"

	dbInfo, err := DBInfoFromSecret(secretString)
	assert.NoError(err, "No error should be returned on successful lookup")
	assert.Equal(
		"dbname=db_name user=db_user password=db_password host=test_db_hostname port=5432 sslmode=disable",
		dbInfo.ConnectionString(),
		"A correct DSN like connection string should have been returned")
}

func TestDBInfoFromConnectionString_empty_string(t *testing.T) {
	assert := require.New(t)

	connectionString := ""

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.EqualError(err, "Malformed connection string element: ")
	assert.Nil(dbInfo, "Empty connection string should yield a nil pointer")
}

func TestDBInfoFromConnectionString_malformed_string(t *testing.T) {
	assert := require.New(t)

	connectionString := "port=1234 host=localhost userTomas password=pass1234 dbname=test-db sslmode=disable"

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.EqualError(err, "Malformed connection string element: userTomas")
	assert.Nil(dbInfo, "Malformed connection string should yield a nil pointer")
}

func TestDBInfoFromConnectionString_invalid_port(t *testing.T) {
	assert := require.New(t)

	connectionString := "port=ioio host=localhost user=Tomas password=pass1234 dbname=test-db sslmode=disable"

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.EqualError(err, "Invalid DatabaseInfo: port: must be a valid port number.")
	assert.Nil(dbInfo, "Invalid connection string should yield a nil pointer")
}

func TestDBInfoFromConnectionString_invalid_sslmode(t *testing.T) {
	assert := require.New(t)

	connectionString := "port=1717 host=localhost user=Tomas password=pass1234 dbname=test-db sslmode=hypercube"

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.EqualError(err, "Invalid DatabaseInfo: sslmode: must be a valid value.")
	assert.Nil(dbInfo, "Invalid connection string should yield a nil pointer")
}

func TestDBInfoFromConnectionString_invalid_password(t *testing.T) {
	assert := require.New(t)

	connectionString := "port=1010 host=localhost user=Tomas password=A dbname=test-db sslmode=disable"

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.EqualError(err, "Invalid DatabaseInfo: password: the length must be between 3 and 64.")
	assert.Nil(dbInfo, "Invalid connection string should yield a nil pointer")
}

func TestDBInfoFromConnectionString(t *testing.T) {
	assert := require.New(t)

	connectionString := "port=1010 host=localhost user=Tomas password=ABA1234 dbname=test-db sslmode=disable"

	dbInfo, err := DBInfoFromConnectionString(connectionString)

	assert.NoError(err, "On valid connection string the there should be no error")
	assert.Equal("test-db", dbInfo.Database, "The database name should be test-db")
	assert.Equal("Tomas", dbInfo.Username, "The username should be Tomas")
	assert.Equal("ABA1234", dbInfo.Password, "The password should be ABA1234")
	assert.Equal("localhost", dbInfo.Host, "The host should be localhost")
	assert.Equal("1010", dbInfo.Port, "The port should be 1010")
	assert.Equal("disable", dbInfo.SSLMode, "The sslmode should be disable")
	assert.Equal("dbname=test-db user=Tomas password=ABA1234 host=localhost port=1010 sslmode=disable",
		dbInfo.ConnectionString())
}
