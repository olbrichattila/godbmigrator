package migrator_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

type ReportTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestReportTestSuite(t *testing.T) {
	suite.Run(t, new(ReportTestSuite))
}

func (suite *ReportTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
}

func (suite *ReportTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *ReportTestSuite) TestDBMigratorMigrateAllTables() {
	reportSourceFixture := testFixtureFolder + "/migration_reports_fixture.json"
	reportDestinationFixture := testFixtureFolder + "/migration_reports.json"
	err := copyFile(reportSourceFixture, reportDestinationFixture)
	t.Nil(err)

	MigrationProvider, err := migrator.NewMigrationProvider("json", tablePrefix, t.db, true)
	MigrationProvider.SetJSONFilePath(testFixtureFolder)
	t.Nil(err)

	report, err := MigrationProvider.Report()
	t.Nil(err)

	expected := "Created at: 2024-05-27 14:12:31, File Name: 2023-07-27_17_57_47-migrate-fixture.sql, Status: success, Message: ok\nCreated at: 2024-05-27 14:12:31, File Name: 2023-07-27_17_57_50-migrate-fixture.sql, Status: success, Message: ok\n"

	t.Equal(expected, report)
}
