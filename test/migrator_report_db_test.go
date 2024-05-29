package migrator_test

import (
	"database/sql"
	"testing"

	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

type ReportDbTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestReportDbTestSuite(t *testing.T) {
	suite.Run(t, new(ReportDbTestSuite))
}

func (suite *ReportDbTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
}

func (suite *ReportDbTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *ReportDbTestSuite) TestDBMigratorMigrateAllTables() {
	MigrationProvider, err := migrator.NewMigrationProvider("db", t.db)
	t.Nil(err)

	err = haveReportrecord(t.db, "FN1", "2006-01-01 00:00:00", "success", "ok")
	t.Nil(err)

	err = haveReportrecord(t.db, "FN2", "2006-01-02 00:00:00", "error", "table not exists")
	t.Nil(err)

	report, err := MigrationProvider.Report()
	t.Nil(err)

	expected := "Created at: 2006-01-01T00:00:00Z, File Name: FN1, Status: success, Message: ok\nCreated at: 2006-01-02T00:00:00Z, File Name: FN2, Status: error, Message: table not exists\n"

	t.Equal(expected, report)
}
