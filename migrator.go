// Package migrator is a lightweight database migrator package, pass only the *sql.DB, migrate, rollback and migration reports
package migrator

import (
	"database/sql"

	"github.com/olbrichattila/godbmigrator/internal/baseliner"
	"github.com/olbrichattila/godbmigrator/internal/migrate"
	"github.com/olbrichattila/godbmigrator/internal/migrationfile"
)

func NewMigrationProvider(providerType, tablePrefix string, db *sql.DB, createMigrationTables bool) (migrate.MigrationProvider, error) {
	return migrate.NewProvider(providerType, tablePrefix, db, createMigrationTables)
}

// Rollback rolls back last migrated items or all if count is 0
func Rollback(
	db *sql.DB,
	migrationProvider migrate.MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	mm := migrationfile.New(migrationFilePath)
	m := migrate.New(db, mm)
	return m.Rollback(db, migrationProvider, migrationFilePath, count, false)
}

// Refresh runs a full rollback and migrate again
func Refresh(
	db *sql.DB,
	migrationProvider migrate.MigrationProvider,
	migrationFilePath string,
) error {
	mm := migrationfile.New(migrationFilePath)
	m := migrate.New(db, mm)
	err := m.Rollback(db, migrationProvider, migrationFilePath, 0, true)
	if err != nil {
		return err
	}

	return Migrate(db, migrationProvider, migrationFilePath, 0)
}

// Migrate execute migrations
func Migrate(
	db *sql.DB,
	migrationProvider migrate.MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	mm := migrationfile.New(migrationFilePath)
	m := migrate.New(db, mm)
	return m.Migrate(db, migrationProvider, migrationFilePath, count)
}

// Report return a report of the already executed migrations
func Report(
	db *sql.DB,
	migrationProvider migrate.MigrationProvider,
	migrationFilePath string,
) (string, error) {
	mm := migrationfile.New(migrationFilePath)
	m := migrate.New(db, mm)
	return m.Report(db, migrationProvider, migrationFilePath)
}

// AddNewMigrationFiles adds a new blank migration file and a rollback file
func AddNewMigrationFiles(migrationFilePath, customText string) error {
	mf := migrationfile.New(migrationFilePath)
	var err error
	err = mf.CreateNewMigrationFiles(migrationFilePath, customText, false)
	if err != nil {
		return err
	}
	err = mf.CreateNewMigrationFiles(migrationFilePath, customText, true)
	if err != nil {
		return err
	}

	return nil
}

// ChecksumValidation validates if the checksums are correct and nothing changed
func ChecksumValidation(
	db *sql.DB,
	migrationProvider migrate.MigrationProvider,
	migrationFilePath string,
) []string {
	mm := migrationfile.New(migrationFilePath)
	m := migrate.New(db, mm)
	return m.ChecksumValidation(db, migrationProvider, migrationFilePath)
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
