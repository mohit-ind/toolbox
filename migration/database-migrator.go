package migration

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/toolbox/models"
)

// The DatabaseMigrator connects to the database with the help of the databaseDriver, the databaseDialect,
// databaseSchema and connectionString. It loads migration scripts from the migrationSource and performs
// step number of them in the direction of choice. Also it can report MigrationInfo with available
// migration scripts, and the list of already applied ones.
type DatabaseMigrator struct {
	databaseDriver   string
	databaseDialect  string
	databaseSchema   string
	connectionString string
	migrationSource  http.FileSystem
}

// NewDatabaseMigrator creates a new DatabaseMigrator instance with the supplied ConnectionString to access the database,
// and MigrationSource to load migration scripts from.
func NewDatabaseMigrator(DatabaseDriver, DatabaseDialect, DatabaseSchema, ConnectionString string, MigrationSource http.FileSystem) *DatabaseMigrator {
	return &DatabaseMigrator{
		databaseDriver:   DatabaseDriver,
		databaseDialect:  DatabaseDialect,
		databaseSchema:   DatabaseSchema,
		connectionString: ConnectionString,
		migrationSource:  MigrationSource,
	}
}

// getDB connects to the database with the DatabaseDriver and the DatabaseMigrator's connectionString.
// It pings the database and returns the connection object. It may return a connection error if
// it failed to connect or ping the database.
func (mig *DatabaseMigrator) getDB() (*sql.DB, error) {
	db, err := sql.Open(mig.databaseDriver, mig.connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the database")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Database is unreachable")
	}
	return db, nil
}

// getSource returns a migrate.HttpFileSystemMigrationSource
// with the DatabaseMigrator's migrationSource as http.FileSystem.
func (mig *DatabaseMigrator) getSource() *migrate.HttpFileSystemMigrationSource {
	return &migrate.HttpFileSystemMigrationSource{
		FileSystem: mig.migrationSource,
	}
}

// Migrate tries to execute the migration. It returns the number of actually applied migration steps,
// and an optional database error.
func (mig *DatabaseMigrator) Migrate(direction Direction, steps int) (int, error) {
	db, err := mig.getDB()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to set up migration")
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close the database properly: %s", err)
		}
	}()

	source := mig.getSource()

	migrate.SetSchema(mig.databaseSchema)

	appliedSteps, err := migrate.ExecMax(
		db,
		mig.databaseDriver,
		source,
		migrate.MigrationDirection(direction),
		steps,
	)
	if err != nil {
		return appliedSteps, errors.Wrap(err, "Failed to execute migration")
	}

	return appliedSteps, nil
}

// GetMigrationScripts returns the map of MigrationScripts, and an optional lookup error.
func (mig *DatabaseMigrator) GetMigrationScripts() ([]*models.MigrationScript, error) {
	if mig.migrationSource == nil {
		return nil, errors.New("GetMigrationScripts is called, but the migration script source is nil")
	}

	scripts := []*models.MigrationScript{}

	migrations, err := mig.getSource().FindMigrations()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot find migration scripts")
	}
	for _, migration := range migrations {
		scripts = append(scripts, &models.MigrationScript{
			Name:          migration.Id,
			UpStatement:   migration.Up,
			DownStatement: migration.Down,
		})
	}

	return scripts, nil
}

// GetMigrationRows gets the list of migration rows from the database, and an optional database error.
func (mig *DatabaseMigrator) GetMigrationRows() ([]*models.MigrationRow, error) {
	db, err := mig.getDB()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the database to get migration rows")
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close the database properly: %s", err)
		}
	}()

	records, err := migrate.GetMigrationRecords(db, mig.databaseDialect)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get migration rows from the database")
	}

	rows := []*models.MigrationRow{}
	for _, record := range records {
		rows = append(rows, &models.MigrationRow{
			Name:      record.Id,
			AppliedAt: record.AppliedAt,
		})
	}

	return rows, nil
}

// GetMigrationInfo collects the available migration scripts and already applied ones
// and builds the MigrationInfo with the them.
func (mig *DatabaseMigrator) GetMigrationInfo() (*models.MigrationInfo, error) {
	scripts, err := mig.GetMigrationScripts()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to look up available migration scripts")
	}

	rows, err := mig.GetMigrationRows()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to look up migration rows in the database")
	}

	return &models.MigrationInfo{
		MigrationScripts: scripts,
		MigrationRows:    rows,
	}, nil
}
