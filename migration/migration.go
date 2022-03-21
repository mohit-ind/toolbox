// Package migration provides a Migrator object which can execute sql migration scripts against a PostgreSQL database
package migration

import (
	migrate "github.com/rubenv/sql-migrate"
)

// Direction tells the Migrator if it should apply migration scripts forward or backward.
type Direction migrate.MigrationDirection

const (
	// Up migrates forward
	Up = Direction(migrate.Up)

	// Down migrates backwards
	Down = Direction(migrate.Down)
)
