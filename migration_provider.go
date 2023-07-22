package migrator

import "fmt"

type MigrationProvider interface {
	LatestMigrations() []string
	AddToMigration(string) error
	RemoveFromMigration(string) error
	MigrationExistsForFile(string) bool
}

func NewMigrationProvider(providerType string) (MigrationProvider, error) {
	switch providerType {
	case "json":
		return newJsonMigration(), nil
	default:
		return nil, fmt.Errorf("Invalid migration provider type: %s", providerType)
	}
}
