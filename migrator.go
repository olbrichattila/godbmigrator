// Package migrator is a lightweight database migrator package, pass only the *sql.DB, migrate, rollback and migration reports
package migrator

import (
	"database/sql"

	"github.com/olbrichattila/godbmigrator/internal/baseliner"
	"github.com/olbrichattila/godbmigrator/internal/migrate"
	"github.com/olbrichattila/godbmigrator/internal/migrationfile"
)

// Rollback rolls back last migrated items or all if count is 0
func Rollback(
	db *sql.DB,
	tablePrefix string,
	migrationFilePath string,
	count int,
) error {
	m, provider, err := getMigrator(db, migrationFilePath, tablePrefix, true)
	if err != nil {
		return err
	}

	return m.Rollback(provider, migrationFilePath, count, false)
}

// Refresh runs a full rollback and migrate again
func Refresh(
	db *sql.DB,
	tablePrefix string,
	migrationFilePath string,
) error {
	m, provider, err := getMigrator(db, migrationFilePath, tablePrefix, true)
	if err != nil {
		return err
	}

	err = m.Rollback(provider, migrationFilePath, 0, true)
	if err != nil {
		return err
	}

	return m.Migrate(provider, migrationFilePath, 0)
}

// Migrate execute migrations
func Migrate(
	db *sql.DB,
	tablePrefix string,
	migrationFilePath string,
	count int,
) error {
	m, provider, err := getMigrator(db, migrationFilePath, tablePrefix, true)
	if err != nil {
		return err
	}

	return m.Migrate(provider, migrationFilePath, count)
}

// Report return a report of the already executed migrations
func Report(
	db *sql.DB,
	tablePrefix string,
	migrationFilePath string,
) (string, error) {
	m, provider, err := getMigrator(db, migrationFilePath, tablePrefix, true)
	if err != nil {
		return "", err
	}

	return m.Report(provider, migrationFilePath)
}

// AddNewMigrationFiles adds a new blank migration file and a rollback file
func AddNewMigrationFiles(migrationFilePath, customText string) error {
	mf := migrationfile.New(migrationFilePath)
	err := mf.CreateNewMigrationFiles(migrationFilePath, customText)
	if err != nil {
		return err
	}

	return nil
}

// ChecksumValidation validates if the checksums are correct and nothing changed
func ChecksumValidation(
	db *sql.DB,
	tablePrefix string,
	migrationFilePath string,
) []string {
	m, provider, err := getMigrator(db, migrationFilePath, tablePrefix, true)
	if err != nil {
		return []string{err.Error()}
	}

	return m.ChecksumValidation(provider, migrationFilePath)
}

// SaveBaseline will save the current status of your database as baseline, which means the migration can start from this point
func SaveBaseline(
	db *sql.DB,
	migrationFilePath string,
) error {
	b := baseliner.New(db)

	return b.Save(migrationFilePath)
}

// LoadBaseline loads the backed up baseline schema to the database
func LoadBaseline(
	db *sql.DB,
	migrationFilePath string,
) error {
	b := baseliner.New(db)

	return b.Load(migrationFilePath)
}

func getMigrator(db *sql.DB, migrationFilePath, tablePrefix string, createMigrationTables bool) (migrate.Migrator, migrate.MigrationProvider, error) {
	provider, err := migrate.NewProvider(tablePrefix, db, createMigrationTables)

	if err != nil {
		return nil, nil, err
	}

	migrator := migrate.New(
		db,
		migrationfile.New(migrationFilePath),
	)

	return migrator, provider, nil
}
