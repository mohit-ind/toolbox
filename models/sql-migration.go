package models

import "time"

// MigrationScript is a Go representation of an SQL migration script.
type MigrationScript struct {
	// Name is the sql migration script's name
	Name string `json:"name"`

	// UpStatement is the Up part of the migration.
	UpStatement []string `json:"up_statement"`

	// DownStatement is the Down part of the migration.
	DownStatement []string `json:"down_statement"`
}

// MigrationRow is one migrations script's status in the database.
type MigrationRow struct {
	// Name is the sql migration script's name
	Name string `json:"name"`

	// AppliedAt is the date when the scrip were applied.
	AppliedAt time.Time `json:"applied_at"`
}

// MigrationInfo is the collection of migration scripts and database rows.
type MigrationInfo struct {
	// MigrationScripts is a list of available MigrationScripts
	MigrationScripts []*MigrationScript `json:"migration_scripts"`

	// MigrationRows is a list of already applied migrations (from the db table gorp_migrations)
	MigrationRows []*MigrationRow `json:"migration_rows"`
}
