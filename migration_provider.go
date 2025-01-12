package migrator

import (
	"database/sql"
	"fmt"
)

const (
	statusError       = "error"
	statusSuccess     = "success"
	reportMessageText = "Created at: %s, File Name: %s, Status: %s, Message: %s\n"
	timeFormat        = "2006-01-02 15:04:05"
)

// MigrationProvider is the base migrator interface
type MigrationProvider interface {
	Migrations(bool) ([]MigrationRow, error)
	AddToMigration(string, string) error
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
func NewMigrationProvider(providerType, tablePrefix string, db *sql.DB) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJSONMigration()
	case "db":
		return newDbMigration(db, tablePrefix)
	default:
		return nil, fmt.Errorf("invalid migration provider type: %s", providerType)
	}
}
