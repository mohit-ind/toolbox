package migration

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	cli "github.com/toolboxcli"
	logger "github.com/toolboxlogger"
	models "github.com/toolboxmodels"
)

// Migrator is the interface of objects which can execute migrations
type Migrator interface {
	Migrate(direction Direction, steps int) (int, error)
	GetMigrationScripts() ([]*models.MigrationScript, error)
	GetMigrationRows() ([]*models.MigrationRow, error)
	GetMigrationInfo() (*models.MigrationInfo, error)
}

// MigratorCLI is a wrapper used to generate a CLI interface for a Migrator in the context of a Services.
type MigratorCLI struct {
	migrator  Migrator
	logger    *logger.Logger
	scriptDir string
}

// MigratorCLIOptions are the possible options for the NewMigratorCLI to create a MigratorCLI.
type MigratorCLIOptions struct {
	Migrator  Migrator
	Logger    *logger.Logger
	ScriptDir string
}

// NewMigratorCLI creates a new MigratorCLI with the supplied MigratorCLIOptions
func NewMigratorCLI(opts MigratorCLIOptions) *MigratorCLI {
	return &MigratorCLI{
		migrator:  opts.Migrator,
		logger:    opts.Logger,
		scriptDir: opts.ScriptDir,
	}
}

// BuildMigrationCommand builds the migration toolbox/cli command. Which is a command line interface of the Migrator.
func (mcli *MigratorCLI) BuildMigrationCommand() *cli.Command {

	migrationCommand := cli.NewCommand("migration").
		WithAliases("migrate", "migrator").
		WithTask(cli.EndWithMessage(mcli.migrationUsage())).
		WithSubCommands(
			mcli.migrateInfoCommand(),
			mcli.migrateGenerateCommand(),
			mcli.migrateUpallCommand(),
			mcli.migrateUpCommand(),
			mcli.migrateDownCommand(),
			mcli.migrateResetCommand(),
		)

	return migrationCommand
}

func (mcli *MigratorCLI) migrateInfoCommand() *cli.Command {
	return cli.NewCommand("info").
		WithAliases("status").
		WithTask(func(args []string) error {
			info, err := mcli.migrator.GetMigrationInfo()
			if err != nil {
				return errors.Wrap(err, "Cannot get migration info")
			}
			mcli.logger.Entry().Info(mcli.renderMigrationInfo(info))
			return nil
		})
}

func (mcli *MigratorCLI) migrateGenerateCommand() *cli.Command {
	return cli.NewCommand("generate").
		WithAliases("gen", "new").
		WithTask(func(args []string) error {
			if len(args) == 0 {
				return errors.New("generate needs at least one script-name as an argument")
			}
			for step, arg := range args {
				fileName, err := mcli.generateMigrationScript(arg)
				if err != nil {
					return errors.Wrap(err, "Failed to generate new migration script")
				}
				mcli.logger.WithField("Path", fileName).Info("New migration script file generated")
				// Wait between steps, so the second in the filename will be different
				if step < len(args)-1 {
					time.Sleep(time.Millisecond * 1333)
				}
			}
			return nil
		})
}

func (mcli *MigratorCLI) migrateUpallCommand() *cli.Command {
	return cli.NewCommand("upall").WithTask(func(args []string) error {
		res, err := mcli.migrator.Migrate(Up, 0)
		if err != nil {
			return errors.Wrap(err, "Failed to migrate Up all")
		}
		mcli.logger.WithField("Steps taken", res).Info("Migrate upall succeeded")
		return nil
	})
}

func (mcli *MigratorCLI) migrateUpCommand() *cli.Command {
	return cli.NewCommand("up").WithTask(func(args []string) error {
		res, err := mcli.migrator.Migrate(Up, 1)
		if err != nil {
			return errors.Wrap(err, "Failed to migrate Up")
		}
		mcli.logger.WithField("Steps taken", res).Info("Migrate up succeeded")
		return nil
	})
}

func (mcli *MigratorCLI) migrateDownCommand() *cli.Command {
	return cli.NewCommand("down").WithTask(func(args []string) error {
		res, err := mcli.migrator.Migrate(Down, 1)
		if err != nil {
			return errors.Wrap(err, "Failed to migrate Down")
		}
		mcli.logger.WithField("Steps taken", res).Info("Migrate down succeeded")
		return nil
	})
}

func (mcli *MigratorCLI) migrateResetCommand() *cli.Command {
	return cli.NewCommand("reset").WithTask(func(args []string) error {
		res, err := mcli.migrator.Migrate(Down, 0)
		if err != nil {
			return errors.Wrap(err, "Failed to reset migrations")
		}
		mcli.logger.WithField("Steps taken", res).Info("Migrate reset succeeded")
		return nil
	})
}

// migrationUsage returns the usage string of the migration command in the Services context.
func (mcli *MigratorCLI) migrationUsage() string {
	return fmt.Sprintf(`
This utility command is used to execute SQL migration scripts against the database. A separate table
called gorp_migrations is created in the database, to keep track of the allready applied migrations.

Migration script directory: %s

Usage: app migrate <subcommand>

Example:

$ ./app migrate info

Available subcommands:
	info     - Prints the available migration scripts and the time they were applied
	upall    - Migrate database all the way up
	up       - Migrate database one step up
	down     - Migrate database one step down
	reset    - Migrate database all the way down
	generate - Generate one or more new migration script files
`, mcli.scriptDir)
}

func (mcli *MigratorCLI) renderMigrationInfo(info *models.MigrationInfo) string {
	migrated := map[string]time.Time{}

	for _, row := range info.MigrationRows {
		migrated[row.Name] = row.AppliedAt
	}

	data := [][]string{}
	for _, script := range info.MigrationScripts {
		appliedAt := "Not Applied!"
		if val, ok := migrated[script.Name]; ok {
			appliedAt = val.Format(time.RFC3339)
		}
		data = append(data, []string{script.Name, appliedAt})
	}

	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf(`
Migration script directory: %s

`, mcli.scriptDir))

	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Migration Script", "Applied At"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.Render()

	tableString.WriteString("\n")
	return tableString.String()
}

// generateMigrationScript a new migration script file and return its path
func (mcli *MigratorCLI) generateMigrationScript(scriptName string) (string, error) {
	// Create migration script file using the UTC time as the first part so the order is guaranteed
	// Format: path/yyyymmddhhmmss-script_name.sql
	// Filename example: scripts/migrations/20201021180150-add_table_users.sql
	fileName := fmt.Sprintf("%s-%s.sql", time.Now().UTC().Format("20060102150405"), strings.TrimSpace(scriptName))
	pathName := path.Join(mcli.scriptDir, fileName)
	f, err := os.Create(pathName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create migration script file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			mcli.logger.WithError(err).Errorf("Failed to close migration file: %s", pathName)
		}
	}()

	// Create template
	templateContent := `-- +migrate Up
-- {{.}}


-- +migrate Down
-- {{.}}
`
	tpl := template.Must(template.New("new_migration").Parse(templateContent))

	if err := tpl.Execute(f, scriptName); err != nil {
		return "", err
	}

	return pathName, nil
}
