package migrator_test

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
)

const testFixtureFolder = "./test_fixtures"

func (t *TestSuite) TestDBMigratorMigrateAllTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)
}

func (t *TestSuite) TestDBMigratorMigrateSpeciedAmountOfTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()

	migrateCount := 2

	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)
}

func (t *TestSuite) TestDBMigratorRollbackAllTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(1, tableCount)
}

func (t *TestSuite) TestDBMigratorRollbackSpecificAmountOfTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *TestSuite) TestDBMigratorRollsBackTablesInProperBatches() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 1)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(2, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(1, tableCount)
}

func (t *TestSuite) TestDBRefresh() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Refresh(db, migrationProvider, testFixtureFolder)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(6, tableCount)
}
