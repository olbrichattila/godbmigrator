package migrator_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

type JsonTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestJsonTestRunner(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
}

func (suite *JsonTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
	resetJsonFile()
}

func (suite *JsonTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *JsonTestSuite) TestJsonMigratorMigrateAllTables() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)
}

func (t *JsonTestSuite) TestJsonMigratorMigrateSpeciedAmountOfTables() {
	migrateCount := 2

	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)
}

func (t *JsonTestSuite) TestJsonMigratorRollbackAllTables() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(0, tableCount)
}

func (t *JsonTestSuite) TestJsonMigratorRollbackSpecificAmountOfTables() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(1, tableCount)
}

func (t *JsonTestSuite) TestJsonMigratorRollsBackTablesInProperBatches() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 1)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)
	time.Sleep(time.Second)
	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(1, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(0, tableCount)
}

func (t *JsonTestSuite) TestJsonRefresh() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, nil, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = migrator.Refresh(t.db, migrationProvider, testFixtureFolder)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)
}
