package migration

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	// PostgreSQL database driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"

	constants "github.com/toolboxconstants"
	logger "github.com/toolboxlogger"
	models "github.com/toolboxmodels"
)

////////////
// Suite //
//////////

// MigratorCLITestSuite extends testify's Suite
type MigratorCLITestSuite struct {
	suite.Suite
	tempDir string
}

func (mcts *MigratorCLITestSuite) setupScriptFolder() {
	tempDir, err := ioutil.TempDir("", "prefix")
	mcts.NoError(err, "Temp dir should have been created")
	mcts.tempDir = tempDir
}

func (mcts *MigratorCLITestSuite) clearScriptFolder() {
	mcts.NoError(os.RemoveAll(mcts.tempDir), "Temp dir should have been removed")
}

func (mcts *MigratorCLITestSuite) createMigrationScript(name string, upStatement []string, downStatement []string) {
	filePath := path.Join(mcts.tempDir, name+".sql")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	mcts.NoErrorf(err, "Script file should have been created at: %s", filePath)
	defer func(f *os.File) {
		mcts.NoError(file.Close(), "script file should have been closed without any error")
	}(file)

	datawriter := bufio.NewWriter(file)

	_, err = datawriter.WriteString(fmt.Sprintf(`-- +migrate Up
-- %s

`, name))
	mcts.NoError(err, "Template string should have been written into the datawriter")

	for _, data := range upStatement {
		_, err := datawriter.WriteString(data + "\n")
		mcts.NoError(err, "Up statement should have been written into the datawriter")
	}

	_, err = datawriter.WriteString(fmt.Sprintf(`
-- +migrate Down
-- %s

`, name))
	mcts.NoError(err, "Template string should have been written into the datawriter")

	for _, data := range downStatement {
		_, err = datawriter.WriteString(data + "\n")
		mcts.NoError(err, "Down statement should have been written into the datawriter")
	}

	mcts.NoError(datawriter.Flush(), "Datawriter should have been flushed without any error")
}

func (mcts *MigratorCLITestSuite) execDB(mig *DatabaseMigrator, sql string, objects ...interface{}) error {
	db, err := mig.getDB()
	mcts.NoError(err, "The DatabaseMigrator's underlying database should be available")
	defer func() {
		mcts.NoError(
			db.Close(),
			"The DatabaseMigrator's underlying database should have been closed without any error",
		)
	}()
	return db.QueryRow(sql).Scan(objects)
}

func (mcts *MigratorCLITestSuite) execAndCapture(f func() error) (string, error) {
	// Save os.Stdout to a var, and switch it with a buffer writer
	oldStdout := os.Stdout
	reader, writer, pipeErr := os.Pipe()
	mcts.NoError(pipeErr, "OS Pipe should have been created")
	os.Stdout = writer

	// Execute the function, save the error
	err := f()

	// Read the custom writer's buffer and restore os.Stdout
	writer.Close()
	out, readErr := ioutil.ReadAll(reader)
	mcts.NoError(readErr, "The buffer should be readable")
	os.Stdout = oldStdout

	// return captured output and error
	return string(out), err
}

////////////
// Tests //
//////////

func (mcts *MigratorCLITestSuite) TestNewMigratorCLI() {
	mCLI := NewMigratorCLI(MigratorCLIOptions{})
	mcts.NotNil(mCLI, "Migrator CLI should have been created")

	cmd := mCLI.BuildMigrationCommand()
	mcts.NotNil(cmd, "Migrator CLI root command should have been created")
}

func (mcts *MigratorCLITestSuite) TestNewMigratorCLI_empty() {
	cmd := NewMigratorCLI(MigratorCLIOptions{}).BuildMigrationCommand()

	out, err := mcts.execAndCapture(func() error { return cmd.Execute(nil) })

	mcts.Error(err, "Executing the root command without arguments should result in an error")
	mcts.Contains(err.Error(), "Command is missing")
	for _, elem := range []string{
		"This utility command is used to execute SQL migration scripts against the database.",
		"$ ./app migrate info",
		"generate - Generate one or more new migration script files",
	} {
		mcts.Containsf(out, elem, "The output: %s should contain: %s", out, elem)
	}
}

type failMigrator struct{}

func (fm *failMigrator) Migrate(direction Direction, steps int) (int, error) {
	return 0, errors.New("Failed to migrate")
}

func (fm *failMigrator) GetMigrationScripts() ([]*models.MigrationScript, error) {
	return nil, errors.New("Failed to get MigrationScripts")
}

func (fm *failMigrator) GetMigrationRows() ([]*models.MigrationRow, error) {
	return nil, errors.New("Failed to get MigrationRows")
}

func (fm *failMigrator) GetMigrationInfo() (*models.MigrationInfo, error) {
	return nil, errors.New("Failed to get MigrationInfo")
}

func (mcts *MigratorCLITestSuite) TestMigratorCLI_errors() {
	nullLogger, _ := logrusTest.NewNullLogger()

	log := logger.NewLogger(nullLogger, nil)

	cmd := NewMigratorCLI(MigratorCLIOptions{
		Migrator:  &failMigrator{},
		Logger:    log,
		ScriptDir: "/a/folder/that/does/not/exists",
	}).BuildMigrationCommand()

	testCases := map[string]struct {
		args        []string
		expectedErr string
	}{
		"Generate script failure": {
			args:        []string{"generate", "new-script"},
			expectedErr: "Failed to generate new migration script: Failed to create migration script file: open /a/folder/that/does/not/exists",
		},
		"Get MigrationInfo failure": {
			args:        []string{"info"},
			expectedErr: "Cannot get migration info: Failed to get MigrationInfo",
		},
		"Migrate up failure": {
			args:        []string{"up"},
			expectedErr: "Failed to migrate Up: Failed to migrate",
		},
		"Migrate upall failure": {
			args:        []string{"upall"},
			expectedErr: "Failed to migrate Up all: Failed to migrate",
		},
		"Migrate down failure": {
			args:        []string{"down"},
			expectedErr: "Failed to migrate Down: Failed to migrate",
		},
		"Migrate reset failure": {
			args:        []string{"reset"},
			expectedErr: "Failed to reset migrations: Failed to migrate",
		},
	}

	for testCaseName, testCase := range testCases {
		mcts.T().Logf("Migrator CLI failMigrator test: %s", testCaseName)

		err := cmd.Execute(testCase.args)

		mcts.Error(err, "This operation should return an error")
		mcts.Containsf(
			err.Error(),
			testCase.expectedErr,
			"The Migrator CLI error: %s should contain: %s",
			err.Error(),
			testCase.expectedErr,
		)
	}
}

////////////////////////
// Integration Tests //
//////////////////////

func (mcts *MigratorCLITestSuite) TestMigratorCLI() {
	if os.Getenv("CI") != "true" {
		mcts.T().Skip("Skipping Migrator CLI integration tests due to CI is not 'true'")
	}

	type operation struct {
		name             string
		args             []string
		expectedError    string
		expectedInOutput []string
	}

	testCases := map[string]struct {
		scripts    []models.MigrationScript
		operations []operation
	}{
		"null test": {
			scripts:    []models.MigrationScript{models.MigrationScript{}},
			operations: []operation{},
		},
		"empty test": {
			scripts: []models.MigrationScript{models.MigrationScript{}},
			operations: []operation{
				{
					name: "Check MigrationInfo",
					args: []string{"info"},
					expectedInOutput: []string{
						"MIGRATION SCRIPT | APPLIED AT",
					},
				},
			},
		},
		"script generation test": {
			scripts: []models.MigrationScript{models.MigrationScript{}},
			operations: []operation{
				{
					name:          "Without arg",
					args:          []string{"generate"},
					expectedError: "generate needs at least one script-name as an argument",
				},
				{
					name: "With two args",
					args: []string{"generate", "one", "two"},
					expectedInOutput: []string{
						"New migration script file generated",
					},
				},
			},
		},
		"multi script migration test": {
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
					name: "Migrate up",
					args: []string{"up"},
					expectedInOutput: []string{
						"Migrate up succeeded",
					},
				},
				{
					name: "Check MigrationInfo",
					args: []string{"info"},
					expectedInOutput: []string{
						"03-add-table-users.sql       | 20",
						"04-add-username-to-users.sql | Not Applied!",
						"06-create-table-todos.sql    | Not Applied!",
					},
				},
				{
					name: "Migrate upall",
					args: []string{"upall"},
					expectedInOutput: []string{
						"Migrate upall succeeded",
					},
				},
				{
					name: "Migrate down",
					args: []string{"down"},
					expectedInOutput: []string{
						"Migrate down succeeded",
					},
				},
				{
					name: "Check MigrationInfo",
					args: []string{"info"},
					expectedInOutput: []string{
						"03-add-table-users.sql       | 20",
						"05-add-email-to-users.sql    | 20",
						"06-create-table-todos.sql    | Not Applied!",
					},
				},
				{
					name: "Migrate reset",
					args: []string{"reset"},
					expectedInOutput: []string{
						"Migrate reset succeeded",
					},
				},
				{
					name: "Check MigrationInfo",
					args: []string{"info"},
					expectedInOutput: []string{
						"03-add-table-users.sql       | Not Applied!",
						"06-create-table-todos.sql    | Not Applied!",
					},
				},
			},
		},
	}

	for testCaseName, testCase := range testCases {
		mcts.T().Logf("Integration Tests: Migrations CLI: %s", testCaseName)
		mcts.setupScriptFolder()
		for _, script := range testCase.scripts {
			mcts.createMigrationScript(script.Name, script.UpStatement, script.DownStatement)
		}

		migrator := NewDatabaseMigrator(
			constants.DefaultDatabaseDriver,
			constants.DefaultDatabaseDialect,
			constants.DefaultDatabaseSchema,
			os.Getenv("TEST_DB_CONNECTION_STRING"),
			http.Dir(mcts.tempDir),
		)

		nullLogger, hook := logrusTest.NewNullLogger()

		log := logger.NewLogger(nullLogger, nil)

		cmd := NewMigratorCLI(MigratorCLIOptions{
			Migrator:  migrator,
			Logger:    log,
			ScriptDir: mcts.tempDir,
		}).BuildMigrationCommand()

		for _, operation := range testCase.operations {
			mcts.T().Logf("Operation: %s", operation.name)

			err := cmd.Execute(operation.args)

			if operation.expectedError != "" {
				mcts.Error(err, "This operation should return an error")
				mcts.Containsf(
					err.Error(),
					operation.expectedError,
					"The error: %s should contain: %s",
					err.Error(),
					operation.expectedError,
				)
				continue
			}

			out := hook.LastEntry().Message
			for _, elem := range operation.expectedInOutput {
				mcts.Containsf(
					out,
					elem,
					"The output of the operation: %s should contain: %s",
					out,
					elem,
				)
			}
		}

		mcts.clearScriptFolder()
		mcts.Equal(
			sql.ErrNoRows,
			mcts.execDB(migrator, "DROP TABLE IF EXISTS gorp_migrations, users, todos"),
			"This operation shouldn't return any error",
		)
	}
}

// TestMigratorCLI runs the whole test suite
func TestMigratorCLI(t *testing.T) {
	suite.Run(t, new(MigratorCLITestSuite))
}
