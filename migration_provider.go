package migrator

import (
	"database/sql"
	"fmt"
)

// MigrationProvider is the base migrator interface
type MigrationProvider interface {
	Migrations(bool) ([]string, error)
	AddToMigration(string) error
	RemoveFromMigration(string) error
	MigrationExistsForFile(string) (bool, error)
	ResetDate()
	GetJSONFileName() string
	SetJSONFilePath(string)
	AddToMigrationReport(string, error) error
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
