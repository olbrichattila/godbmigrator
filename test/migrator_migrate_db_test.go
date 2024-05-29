package migrator_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

const testFixtureFolder = "./test_fixtures"

type DbTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestDbRunner(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
}

func (suite *DbTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *DbTestSuite) TestDBMigratorMigrateAllTables() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	reportCount, err := rowCountInTable(t.db, "migration_reports")
	t.Nil(err)

	t.Equal(5, reportCount)
}

func (t *DbTestSuite) TestDBMigratorMigrateSpeciedAmountOfTables() {
	migrateCount := 2

	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackAllTables() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackSpecificAmountOfTables() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollsBackTablesInProperBatches() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 1)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 2)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Rollback(t.db, MigrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBRefresh() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Refresh(t.db, MigrationProvider, testFixtureFolder)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}
