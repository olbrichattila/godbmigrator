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
	migrator migrator.DBMigrator
	db       *sql.DB
}

func TestBaselineSuite(t *testing.T) {
	suite.Run(t, new(baselineTestSuite))
}

func (suite *baselineTestSuite) SetupTest() {
	suite.db = initMemorySqlite()
	suite.migrator = migrator.New(suite.db, "./test_fixtures_baseliner", tablePrefix)
}

func (suite *baselineTestSuite) TearDownTest() {
	suite.db.Close()
}

func (t *baselineTestSuite) TestBaselineCreatedAndRestored() {
	// Load and test if the correct number of tables, views, indexes and triggers are created
	err := t.migrator.LoadBaseline()
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
	err = t.migrator.SaveBaseline("test-fixture-baseliner-save")
	t.NoError(err)

	fileSize, err := getFileSize("./test-fixture-baseliner-save/baseline.sql")
	t.NoError(err)
	t.Greater(fileSize, int64(0))
}
