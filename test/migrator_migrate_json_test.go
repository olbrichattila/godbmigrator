package migrator_test

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
)

func (t *TestSuite) TestJsonMigratorMigrateAllTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)
}

func (t *TestSuite) TestJsonMigratorMigrateSpeciedAmountOfTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()

	migrateCount := 2

	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(2, tableCount)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)
}

func (t *TestSuite) TestJsonMigratorRollbackAllTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(0, tableCount)
}

func (t *TestSuite) TestJsonMigratorRollbackSpecificAmountOfTables() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(1, tableCount)
}

func (t *TestSuite) TestJsonMigratorRollsBackTablesInProperBatches() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
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

	t.Equal(5, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(1, tableCount)

	err = migrator.Rollback(db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(0, tableCount)
}

func (t *TestSuite) TestJsonRefresh() {
	resetrTestMigrationPath(testMigrationFilePath)
	resetDatabase()
	resetJsonFile()
	db, err := inMemorySqlite()
	defer db.Close()

	t.Nil(err)

	migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	t.Nil(err)

	err = migrator.Migrate(db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Refresh(db, migrationProvider, testFixtureFolder)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(db)
	t.Nil(err)

	t.Equal(5, tableCount)
}
