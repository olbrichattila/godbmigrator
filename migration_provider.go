package migrator

import (
	"database/sql"
	"fmt"
)

type MigrationProvider interface {
	LatestMigrations() ([]string, error)
	AddToMigration(string) error
	RemoveFromMigration(string) error
	MigrationExistsForFile(string) bool
}

func NewMigrationProvider(providerType string, db *sql.DB) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJsonMigration(), nil
	case "db":
		dbMigration, err := newDbMigration(db)
		if err != nil {
			return nil, err
		}
		return dbMigration, nil
	default:
		return nil, fmt.Errorf("Invalid migration provider type: %s", providerType)
	}
}
