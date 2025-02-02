package migrator_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

type baselineTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestBaselineSuite(t *testing.T) {
	suite.Run(t, new(baselineTestSuite))
}

func (suite *baselineTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
}

func (suite *baselineTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *baselineTestSuite) TestDBMigratorMigrateAllTables() {
	// Load and test if the correct number of tables, views, indexes and triggers are created
	err := migrator.LoadBaseline(t.db, "./test_fixtures_baseliner")
	t.NoError(err)

	tableCount, err := countInSqliteMasterForType(t.db, "table")
	t.NoError(err)
	t.Equal(15, tableCount)

	indexCount, err := countInSqliteMasterForType(t.db, "index")
	t.NoError(err)
	t.Equal(2, indexCount)

	viewCount, err := countInSqliteMasterForType(t.db, "view")
	t.NoError(err)
	t.Equal(1, viewCount)

	triggerCount, err := countInSqliteMasterForType(t.db, "trigger")
	t.NoError(err)
	t.Equal(1, triggerCount)

	// Save it back and test if the saved file is not empty
	err = migrator.SaveBaseline(t.db, ".")
	t.NoError(err)

	fileSize, err := getFileSize("./baseline.sql")
	t.NoError(err)
	t.Greater(fileSize, int64(0))
}
