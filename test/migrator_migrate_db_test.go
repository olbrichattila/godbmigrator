package migrator_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

const (
	testFixtureFolder         = "./test_fixtures"
	testChecksumFixtureFolder = "./test_fixtures_checksum"
	tablePrefix               = "olb"
)

type DbTestSuite struct {
	suite.Suite
	db               *sql.DB
	migrator         migrator.DBMigrator
	checksumMigrator migrator.DBMigrator
}

func TestDbRunner(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
	suite.migrator = migrator.New(suite.db, testFixtureFolder, tablePrefix)
	suite.checksumMigrator = migrator.New(suite.db, testChecksumFixtureFolder, tablePrefix)
}

func (suite *DbTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *DbTestSuite) TestDBMigratorMigrateAllTables() {
	err := t.migrator.Migrate(0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	reportCount, err := rowCountInTable(t.db, tablePrefix+"_migration_reports")
	t.Nil(err)

	t.Equal(5, reportCount)
}

func (t *DbTestSuite) TestDBMigratorMigrateSpecifiedAmountOfTables() {
	migrateCount := 2

	err := t.migrator.Migrate(migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = t.migrator.Migrate(migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = t.migrator.Migrate(migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackAllTables() {
	err := t.migrator.Migrate(0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = t.migrator.Rollback(0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackSpecificAmountOfTables() {
	err := t.migrator.Migrate(0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = t.migrator.Rollback(2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = t.migrator.Rollback(2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollsBackTablesInProperBatches() {
	err := t.migrator.Migrate(1)
	t.Nil(err)
	time.Sleep(time.Second)
	err = t.migrator.Migrate(2)
	t.Nil(err)
	time.Sleep(time.Second)
	err = t.migrator.Migrate(2)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = t.migrator.Rollback(0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = t.migrator.Rollback(0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)

	err = t.migrator.Rollback(0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBRefresh() {
	err := t.migrator.Migrate(3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = t.migrator.Refresh()
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}

func (t *DbTestSuite) TestDBChecksum() {
	err := t.checksumMigrator.Migrate(3)
	t.Nil(err)
	// TODO probably will be broken

	testFixtureFile := "2023-07-27_17_57_47-fixture.sql"

	checksum, err := getChecksumFromTable(t.db, testFixtureFile)
	t.Nil(err)

	hash, err := calculateFileMD5(testFixtureFolder + "/" + testFixtureFile)

	t.Nil(err)
	t.Equal(hash, checksum)
}

func (t *DbTestSuite) TestDBChecksumValidator() {
	err := t.checksumMigrator.Migrate(3)
	t.Nil(err)

	errors := t.checksumMigrator.ChecksumValidation()
	t.Len(errors, 0)
}
