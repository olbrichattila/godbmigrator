package migrator

import (
	"database/sql"
	"fmt"
)

// MigrationProvider is the base migrator interface
type MigrationProvider interface {
	migrations(bool) ([]string, error)
	addToMigration(string) error
	removeFromMigration(string) error
	migrationExistsForFile(string) (bool, error)
	resetDate()
	getJSONFileName() string
	SetJSONFilePath(string)
	addToMigrationReport(string, error) error
	Report() (string, error)
}

// NewMigrationProvider returns a migration provider, which follows the provider type
// The provider type can be json or db, error returned if the type incorrectly provided
// db should be your database *sql.DB, which can be MySQL, Postgres, Sqlite or Firebird
func NewMigrationProvider(providerType string, db *sql.DB) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJSONMigration()
	case "db":
		return newDbMigration(db)
	default:
		return nil, fmt.Errorf("invalid migration provider type: %s", providerType)
	}
}
