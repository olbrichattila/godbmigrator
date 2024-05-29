package migrator

import (
	"database/sql"
	"fmt"
)

type MigrationProvider interface {
	migrations(bool) ([]string, error)
	addToMigration(string) error
	removeFromMigration(string) error
	migrationExistsForFile(string) (bool, error)
	resetDate()
	getJsonFileName() string
	SetJsonFilePath(string)
	AddToMigrationReport(string, error) error
	Report() (string, error)
}

func NewMigrationProvider(providerType string, db *sql.DB) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJsonMigration()
	case "db":
		return newDbMigration(db)
	default:
		return nil, fmt.Errorf("invalid migration provider type: %s", providerType)
	}
}
