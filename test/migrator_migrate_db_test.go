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
	testFixtureFolder = "./test_fixtures"
	testChecksumFixtureFolder = "./test_fixtures_checksum"
	tablePrefix = "olb"
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
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 0)
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

	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
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
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
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
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
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
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
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
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
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


func (t *DbTestSuite) TestDBChecksum() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	testFixtureFile := "2023-07-27_17_57_47-migrate-fixture.sql";

	checksum, err := getChecksumFromTable(t.db, testFixtureFile)

	t.Nil(err)

	hash, err := calculateFileMD5(testFixtureFolder + "/" + testFixtureFile)

	t.Nil(err)
	t.Equal(hash, checksum)
}

func (t *DbTestSuite) TestDBChecksumValidator() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", tablePrefix, t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 3)
	t.Nil(err)

	errors := migrator.ChecksumValidation(t.db, MigrationProvider, testChecksumFixtureFolder)
	t.Len(errors, 1)
}

func (t *DbTestSuite) TestJSONChecksumValidator() {
	MigrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, t.db)
	t.Nil(err)

	err = migrator.Migrate(t.db, MigrationProvider, testFixtureFolder, 3)
	t.Nil(err)
	
	// It cross saves the file, need to resolve this in test
	// errors := migrator.ChecksumValidation(t.db, MigrationProvider, testChecksumFixtureFolder)
	// t.Len(errors, 1)
}