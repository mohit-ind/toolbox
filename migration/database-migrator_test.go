package migration

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/toolbox/constants"
	"github.com/toolbox/models"

	// PostgreSQL database driver
	_ "github.com/lib/pq"
)

////////////
// Suite //
//////////

// DatabaseMigratorTestSuite extends testify's Suite
type DatabaseMigratorTestSuite struct {
	suite.Suite
	tempDir string
}

func (dmts *DatabaseMigratorTestSuite) setupScriptFolder() {
	tempDir, err := ioutil.TempDir("", "prefix")
	dmts.NoError(err, "Temp dir should have been created")
	dmts.tempDir = tempDir
}

func (dmts *DatabaseMigratorTestSuite) clearScriptFolder() {
	dmts.NoError(os.RemoveAll(dmts.tempDir), "Temp dir should have been removed")
}

func (dmts *DatabaseMigratorTestSuite) stringInList(s string, list []string) {
	for _, elem := range list {
		if strings.Contains(elem, s) {
			return
		}
	}
	dmts.Failf("String is not in list", "String: %s, list: %s", s, list)
}

func (dmts *DatabaseMigratorTestSuite) createMigrationScript(name string, upStatement []string, downStatement []string) {
	filePath := path.Join(dmts.tempDir, name+".sql")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	dmts.NoErrorf(err, "Script file should have been created at: %s", filePath)
	defer func(f *os.File) {
		dmts.NoError(file.Close(), "script file should have been closed without any error")
	}(file)

	datawriter := bufio.NewWriter(file)

	_, err = datawriter.WriteString(fmt.Sprintf(`-- +migrate Up
-- %s

`, name))
	dmts.NoError(err, "Template string should have been written into the datawriter")

	for _, data := range upStatement {
		_, err := datawriter.WriteString(data + "\n")
		dmts.NoError(err, "Up statement should have been written into the datawriter")
	}

	_, err = datawriter.WriteString(fmt.Sprintf(`
-- +migrate Down
-- %s

`, name))
	dmts.NoError(err, "Template string should have been written into the datawriter")

	for _, data := range downStatement {
		_, err = datawriter.WriteString(data + "\n")
		dmts.NoError(err, "Down statement should have been written into the datawriter")
	}

	dmts.NoError(datawriter.Flush(), "Datawriter should have been flushed without any error")
}

func (dmts *DatabaseMigratorTestSuite) execDB(mig *DatabaseMigrator, sql string, objects ...interface{}) error {
	db, err := mig.getDB()
	dmts.NoError(err, "The DatabaseMigrator's underlying database should be available")
	defer func() {
		dmts.NoError(
			db.Close(),
			"The DatabaseMigrator's underlying database should have been closed without any error",
		)
	}()
	return db.QueryRow(sql).Scan(objects)
}

////////////
// Tests //
//////////

func (dmts *DatabaseMigratorTestSuite) TestNewDatabaseMigrator() {
	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()
	migrator := NewDatabaseMigrator("", "", "", "", nil)
	dmts.NotNil(migrator, "DatabaseMigrator should have been created")
	scripts, err := migrator.GetMigrationScripts()
	dmts.Error(err, "Without an SQL migration script source the Migrator should return an error")
	dmts.Nil(scripts, "Failed GetMigrationScripts shouldn't return any scripts but a nil pointer")

	migrations, err := migrator.GetMigrationRows()
	dmts.Error(err, "Without a proper connection string this lookup should fail")
	dmts.Nil(migrations, "Failed GetMigrationRows shouldn't return any migration but a nil pointer")

	info, err := migrator.GetMigrationInfo()
	dmts.Error(err, "Without available migrations and database rows GetMigrationInfo should fail")
	dmts.Nil(info, "Failed GetMigrationInfo shouldn't return any info but a nil pointer")

	stepsApplied, err := migrator.Migrate(Up, 999)
	dmts.Error(err, "Without available migration scripts and database Migrate should fail")
	dmts.Equal(0, stepsApplied, "Failed Migrate should return 0 applied steps")
}

func (dmts *DatabaseMigratorTestSuite) TestGetMigrationScripts_invalid_script() {
	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()
	dmts.createMigrationScript("add_testing_table", []string{"Invalid migration script 🐸"}, nil)

	migrator := NewDatabaseMigrator("", "", "", "", http.Dir(dmts.tempDir))

	_, err := migrator.GetMigrationScripts()

	dmts.Error(err, "The Migrator should return on error on invalid migration script")
}

func (dmts *DatabaseMigratorTestSuite) TestGetMigrationScripts() {
	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()
	dmts.createMigrationScript("add_testing_table", []string{"CREATE TABLE testing;"}, []string{"DROP TABLE testing;"})

	migrator := NewDatabaseMigrator("", "", "", "", http.Dir(dmts.tempDir))

	scripts, err := migrator.GetMigrationScripts()
	dmts.NoError(err, "The Migrator should return the migration scripts without any error")
	dmts.Equal(1, len(scripts), "Exactly one migration script should have been returned")
	dmts.Equal("add_testing_table.sql", scripts[0].Name, "The script's Name should be: add_testing_table")
	dmts.stringInList("CREATE TABLE testing;", scripts[0].UpStatement)
	dmts.stringInList("DROP TABLE testing;", scripts[0].DownStatement)
}

func (dmts *DatabaseMigratorTestSuite) TestGetMigrationInfo_no_database() {
	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()
	migrator := NewDatabaseMigrator(
		"",
		"",
		"",
		"",
		http.Dir(dmts.tempDir),
	)

	info, err := migrator.GetMigrationInfo()
	dmts.Error(err, "GetMigrationInfo return an error when unable to fetch migration rows from the database")
	dmts.Nil(info, "If GetMigrationInfo fails it shouldn't return any MigrationInfo, but a nil pointer")
}

////////////////////////
// Integration tests //
//////////////////////

func (dmts *DatabaseMigratorTestSuite) TestGetMigrationRows_empty() {
	if os.Getenv("CI") != "true" {
		dmts.T().Skip("Skipping DatabaseMigrator integration tests due to CI is not 'true'")
	}

	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()

	migrator := NewDatabaseMigrator(
		constants.DefaultDatabaseDriver,
		constants.DefaultDatabaseDialect,
		constants.DefaultDatabaseSchema,
		os.Getenv("TEST_DB_CONNECTION_STRING"),
		nil,
	)

	dmts.Equal(
		sql.ErrNoRows,
		dmts.execDB(migrator, "DROP TABLE IF EXISTS gorp_migrations"),
		"This operation shouldn't return any error",
	)

	rows, err := migrator.GetMigrationRows()
	dmts.NoError(err, "Migration rows should be queried without any error")
	dmts.Equal(0, len(rows), "There should be no migration row in a fresh database")

	dmts.Equal(
		sql.ErrNoRows,
		dmts.execDB(migrator, "SELECT * FROM gorp_migrations"),
		"Empty gorp_migrations table should have been created",
	)
	dmts.Equal(
		sql.ErrNoRows,
		dmts.execDB(migrator, "DROP TABLE IF EXISTS gorp_migrations"),
		"This operation shouldn't return any error",
	)
}

func (dmts *DatabaseMigratorTestSuite) TestGetMigrationInfo() {
	if os.Getenv("CI") != "true" {
		dmts.T().Skip("Skipping DatabaseMigrator integration tests due to CI is not 'true'")
	}

	dmts.setupScriptFolder()
	defer dmts.clearScriptFolder()

	migrator := NewDatabaseMigrator(
		constants.DefaultDatabaseDriver,
		constants.DefaultDatabaseDialect,
		constants.DefaultDatabaseSchema,
		os.Getenv("TEST_DB_CONNECTION_STRING"),
		http.Dir(dmts.tempDir),
	)

	info, err := migrator.GetMigrationInfo()
	dmts.NoError(err, "MigrationInfo should have been returned without any error")
	dmts.NotNil(info.MigrationScripts, "MigrationInfo should have a list of available migration scripts")
	dmts.NotNil(info.MigrationRows, "MigrationInfo should have a list of migration rows")
}

func (dmts *DatabaseMigratorTestSuite) TestMigrate() {
	if os.Getenv("CI") != "true" {
		dmts.T().Skip("Skipping DatabaseMigrator integration tests due to CI is not 'true'")
	}

	type operation struct {
		name          string
		direction     Direction
		steps         int
		expectedSteps int
		expectedRows  int
		expectedError string
	}

	testCases := map[string]struct {
		scripts    []models.MigrationScript
		operations []operation
	}{
		"null test": {
			scripts:    nil,
			operations: nil,
		},
		"invalid SQL statement": {
			scripts: []models.MigrationScript{
				models.MigrationScript{
					Name:        "01-an-invalid-statement",
					UpStatement: []string{"THIS WONT WORK"},
				},
			},
			operations: []operation{
				{
					name:          "A failed migration",
					expectedError: "Failed to execute migration: Error while parsing 01-an-invalid-statement.sql: Error parsing migration (01-an-invalid-statement.sql): ERROR: The last statement must be ended by a semicolon or '-- +migrate StatementEnd' marker.\n\t\t\tSee https://github.com/rubenv/sql-migrate for details.",
				},
			},
		},
		"add a table": {
			scripts: []models.MigrationScript{
				models.MigrationScript{
					Name: "02-add-table-users",
					UpStatement: []string{
						"CREATE TABLE IF NOT EXISTS users (",
						"id SERIAL PRIMARY KEY",
						");",
					},
					DownStatement: []string{"DROP TABLE IF EXISTS users;"},
				},
			},
			operations: []operation{
				{
					name:          "Migrate up",
					direction:     Up,
					steps:         1,
					expectedSteps: 1,
					expectedRows:  1,
				},
				{
					name:          "Migrate down",
					direction:     Down,
					steps:         1,
					expectedSteps: 1,
					expectedRows:  0,
				},
			},
		},
		"multi step migration": {
			scripts: []models.MigrationScript{
				models.MigrationScript{
					Name:          "03-add-table-users",
					UpStatement:   []string{"CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY);"},
					DownStatement: []string{"DROP TABLE IF EXISTS users;"},
				},
				models.MigrationScript{
					Name:          "04-add-username-to-users",
					UpStatement:   []string{"ALTER TABLE users ADD COLUMN username VARCHAR(32) NOT NULL;"},
					DownStatement: []string{"ALTER TABLE users DROP COLUMN username;"},
				},
				models.MigrationScript{
					Name:          "05-add-email-to-users",
					UpStatement:   []string{"ALTER TABLE users ADD COLUMN email VARCHAR(64) NOT NULL;"},
					DownStatement: []string{"ALTER TABLE users DROP COLUMN email;"},
				},
				models.MigrationScript{
					Name:          "06-create-table-todos",
					UpStatement:   []string{"CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY);"},
					DownStatement: []string{"DROP TABLE IF EXISTS todos;"},
				},
			},
			operations: []operation{
				{
					name:          "Migrate up one step",
					direction:     Up,
					steps:         1,
					expectedSteps: 1,
					expectedRows:  1,
				},
				{
					name:          "Migrate all up",
					direction:     Up,
					steps:         0,
					expectedSteps: 3,
					expectedRows:  4,
				},
				{
					name:          "Migrate one step down",
					direction:     Down,
					steps:         1,
					expectedSteps: 1,
					expectedRows:  3,
				},
				{
					name:          "Migrate all down",
					direction:     Down,
					steps:         0,
					expectedSteps: 3,
					expectedRows:  0,
				},
			},
		},
	}

	for testCaseName, testCase := range testCases {
		dmts.T().Logf("Integration Tests: Migrations: %s", testCaseName)
		dmts.setupScriptFolder()
		for _, script := range testCase.scripts {
			dmts.createMigrationScript(script.Name, script.UpStatement, script.DownStatement)
		}
		migrator := NewDatabaseMigrator(
			constants.DefaultDatabaseDriver,
			constants.DefaultDatabaseDialect,
			constants.DefaultDatabaseSchema,
			os.Getenv("TEST_DB_CONNECTION_STRING"),
			http.Dir(dmts.tempDir),
		)
		for _, operation := range testCase.operations {
			dmts.T().Logf("Operation: %s", operation.name)
			stepsTaken, err := migrator.Migrate(operation.direction, operation.steps)
			if operation.expectedError != "" {
				dmts.EqualError(err, operation.expectedError, "This operation should fail")
				continue
			}
			dmts.NoError(err, "This operation should have been executed without any error")
			dmts.Equalf(
				operation.expectedSteps,
				stepsTaken,
				"Number of expected taken steps: %d actual number of steps taken: %d",
				operation.expectedSteps,
				stepsTaken,
			)

			info, err := migrator.GetMigrationInfo()
			dmts.NoError(err, "MigrationInfo should have been returned without any error")
			dmts.Equal(
				operation.expectedRows,
				len(info.MigrationRows),
				"Number of expected MigrationRows: %d actual: %d",
				operation.expectedRows,
				len(info.MigrationRows),
			)
		}
		dmts.clearScriptFolder()
		dmts.Equal(
			sql.ErrNoRows,
			dmts.execDB(migrator, "DROP TABLE IF EXISTS gorp_migrations, users, todos"),
			"This operation shouldn't return any error",
		)
	}
}

// TestDatabaseMigrator runs the whole test suite
func TestDatabaseMigrator(t *testing.T) {
	suite.Run(t, new(DatabaseMigratorTestSuite))
}
