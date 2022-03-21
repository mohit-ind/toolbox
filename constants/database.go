package constants

const (
	// DefaultDatabaseDriver is postgres as we will deal with PostgreSQL.
	// Blank import of _ "github.com/lib/pq" registers the 'postgres' database driver.
	DefaultDatabaseDriver = "postgres"

	// DefaultDatabaseDialect is postgres as we will deal with PostgreSQL.
	DefaultDatabaseDialect = "postgres"

	// DefaultDatabaseSchema is the default database schema.
	DefaultDatabaseSchema = "public"
)
