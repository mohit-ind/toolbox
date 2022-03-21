package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cli "github.com/toolboxcli"
	constants "github.com/toolboxconstants"
	logger "github.com/toolboxlogger"

	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	migration "github.com/toolboxmigration"

	// Blank import of the statik package is required. So we can access the migration scripts.
	_ "github.com/toolboxexamples/migration/statik"

	// Blank import of lib/pq of 'postgres' database driver.
	_ "github.com/lib/pq"
)

const (
	version = "v0.0.1"
)

func main() {
	// get the migration scripts's source from the statik filesystem
	scriptSource, err := fs.New()
	if err != nil {
		log.Fatalf("Failed to get Statik folder: %s", err)
	}

	// create the Migrator with the scriptSource
	migrator := migration.NewDatabaseMigrator(
		constants.DefaultDatabaseDriver,
		constants.DefaultDatabaseDialect,
		constants.DefaultDatabaseSchema,
		os.Getenv("TEST_DB_CONNECTION_STRING"),
		scriptSource,
	)

	// create the Application Logger
	log := getLogger()

	// create the Migration CLI with the Migrator and the Logger
	migrationCLI := migration.NewMigratorCLI(migration.MigratorCLIOptions{
		Migrator:  migrator,
		Logger:    log,
		ScriptDir: "sql-migration-scripts",
	})

	// versionCommand prints out the Version string
	versionCommand := cli.NewCommand("version").
		WithAliases("ver", "Version", "--version").
		WithTask(func(args []string) error {
			fmt.Println(version)
			return nil
		})

	// root command prints out the usage
	rootCommand := cli.NewCommand("root").
		WithSubCommands(
			versionCommand,
			migrationCLI.BuildMigrationCommand(),
		).
		WithTask(cli.EndWithMessage(usage()))

	if err := rootCommand.Execute(os.Args[1:]); err != nil {
		log.WithError(err).Fatal("Failed to execute command")
	}
}

// getLogger creates a logger instance for the Migrator CLI
func getLogger() *logger.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
	})
	return logger.NewLogger(log, logrus.Fields{
		"service": "example-service",
		"version": version,
		"env":     constants.ENV_DEV,
		"host":    "localhost",
	})
}

func usage() string {
	return `
Example Migrator CLI

To spin up an empty PostgreSQL Docker container for testing execute:
	$ docker run --name migration-test-db -e POSTGRES_PASSWORD=pass123 -p 5432:5432 -d postgres

Usage:
	$ app <command> <subcommand>

Example:
	$ app migrate info

Commands:
	version - Prints the Short Semantic Version string
	migrate - Manage SQL database migrations
`
}
