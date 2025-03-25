// Package migrator is a lightweight database migrator package, pass only the *sql.DB, migrate, rollback and migration reports
package migrator

import (
	"database/sql"

	"github.com/olbrichattila/godbmigrator/internal/baseliner"
	"github.com/olbrichattila/godbmigrator/internal/messager"
	"github.com/olbrichattila/godbmigrator/internal/migrate"
	"github.com/olbrichattila/godbmigrator/internal/migrationfile"
)

func New(
	db *sql.DB,
	migrationFilePath,
	tablePrefix string,
) DBMigrator {
	return &dbmigrate{
		db:                db,
		migrationFilePath: migrationFilePath,
		tablePrefix:       tablePrefix,
		messDispatch:      messager.New(),
	}
}

// DBMigrator encapsulates migrator functions
type DBMigrator interface {
	SubscribeToMessages(callback messager.CallbackFunc)
	Rollback(count int) error
	Refresh() error
	Migrate(count int) error
	Report() (string, error)
	AddNewMigrationFiles(customText string) error
	ChecksumValidation() []string
	SaveBaseline(files ...string) error
	LoadBaseline(files ...string) error
}

type dbmigrate struct {
	db                *sql.DB
	migrationFilePath string
	tablePrefix       string
	messDispatch      messager.Messager
}

// SubscribeToMessages receive messages from the migrator, events happening
func (d *dbmigrate) SubscribeToMessages(callback messager.CallbackFunc) {
	d.messDispatch.Register(callback)
}

// Rollback rolls back last migrated items or all if count is 0
func (d *dbmigrate) Rollback(
	count int,
) error {
	m, provider, err := d.getMigrator()
	if err != nil {
		return err
	}

	return m.Rollback(provider, d.migrationFilePath, count, false)
}

// Refresh runs a full rollback and migrate again
func (d *dbmigrate) Refresh() error {
	m, provider, err := d.getMigrator()
	if err != nil {
		return err
	}

	err = m.Rollback(provider, d.migrationFilePath, 0, true)
	if err != nil {
		return err
	}

	return m.Migrate(provider, d.migrationFilePath, 0)
}

// Migrate execute migrations
func (d *dbmigrate) Migrate(
	count int,
) error {
	m, provider, err := d.getMigrator()
	if err != nil {
		return err
	}

	return m.Migrate(provider, d.migrationFilePath, count)
}

// Report return a report of the already executed migrations
func (d *dbmigrate) Report() (string, error) {
	m, provider, err := d.getMigrator()
	if err != nil {
		return "", err
	}

	return m.Report(provider, d.migrationFilePath)
}

// AddNewMigrationFiles adds a new blank migration file and a rollback file
func (d *dbmigrate) AddNewMigrationFiles(customText string) error {
	mf := migrationfile.New(d.migrationFilePath)
	files, err := mf.CreateNewMigrationFiles(d.migrationFilePath, customText)
	if err != nil {
		return err
	}

	for _, fileName := range files {
		d.messDispatch.Dispatch(messager.MigrationFileCreated, fileName)
	}

	return nil
}

// ChecksumValidation validates if the checksums are correct and nothing changed
func (d *dbmigrate) ChecksumValidation() []string {
	m, provider, err := d.getMigrator()
	if err != nil {
		return []string{err.Error()}
	}

	return m.ChecksumValidation(provider, d.migrationFilePath)
}

// SaveBaseline will save the current status of your database as baseline, which means the migration can start from this point
func (d *dbmigrate) SaveBaseline(files ...string) error {
	b := baseliner.New(d.db)
	if len(files) == 0 {
		return b.Save(d.migrationFilePath)
	}

	return b.Save(files[0])
}

// LoadBaseline loads the backed up baseline schema to the database
func (d *dbmigrate) LoadBaseline(files ...string) error {
	b := baseliner.New(d.db)

	if len(files) == 0 {
		return b.Load(d.migrationFilePath)
	}

	return b.Load(files[0])
}

func (d *dbmigrate) getMigrator() (migrate.Migrator, migrate.MigrationProvider, error) {
	provider, err := migrate.NewProvider(d.tablePrefix, d.db)

	if err != nil {
		return nil, nil, err
	}

	migrator := migrate.New(
		d.db,
		migrationfile.New(d.migrationFilePath),
		d.messDispatch,
	)

	return migrator, provider, nil
}
