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
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
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

	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(4, tableCount)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(6, tableCount)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, migrateCount)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackAllTables() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollbackSpecificAmountOfTables() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 2)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(3, tableCount)
}

func (t *DbTestSuite) TestDBMigratorRollsBackTablesInProperBatches() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
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

	t.Equal(7, tableCount)

	err = migrator.Rollback(t.db, migrationProvider, testFixtureFolder, 0)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
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

	t.Equal(2, tableCount)
}

func (t *DbTestSuite) TestDBRefresh() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	tableCount, err := tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(5, tableCount)

	err = migrator.Refresh(t.db, migrationProvider, testFixtureFolder)
	t.Nil(err)

	tableCount, err = tableCountInDatabase(t.db)
	t.Nil(err)

	t.Equal(7, tableCount)
}

func (t *DbTestSuite) TestDBChecksum() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	testFixtureFile := "2023-07-27_17_57_47-migrate-fixture.sql"

	checksum, err := getChecksumFromTable(t.db, testFixtureFile)

	t.Nil(err)

	hash, err := calculateFileMD5(testFixtureFolder + "/" + testFixtureFile)

	t.Nil(err)
	t.Equal(hash, checksum)
}

func (t *DbTestSuite) TestDBChecksumValidator() {
	migrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	errors := migrator.ChecksumValidation(t.db, migrationProvider, testChecksumFixtureFolder)
	t.Len(errors, 1)
}

func (t *DbTestSuite) TestJSONChecksumValidator() {
	migrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, t.db, true)
	t.Nil(err)

	err = migrator.Migrate(t.db, migrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	// It cross saves the file, need to resolve this in test
	// errors := migrator.ChecksumValidation(t.db, migrationProvider, testChecksumFixtureFolder)
	// t.Len(errors, 1)
}
