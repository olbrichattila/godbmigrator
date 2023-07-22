package migrator

import (
	"database/sql"
	"fmt"
)

type MigrationProvider interface {
	LatestMigrations() []string
	AddToMigration(string) error
	RemoveFromMigration(string) error
	MigrationExistsForFile(string) bool
}

func NewMigrationProvider(providerType string, db *sql.DB) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJsonMigration(), nil
	case "db":
		return newDbMigration(db), nil
	default:
		return nil, fmt.Errorf("Invalid migration provider type: %s", providerType)
	}
}
