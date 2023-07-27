package migrator_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const testMigrationFilePath = "./test-migrations"
const dataFolder = "./data"

type TestSuite struct {
	suite.Suite
}

func TestRunner(t *testing.T) {
	initFolder(testMigrationFilePath)
	initFolder(dataFolder)

	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) TestRuns() {
	t.True(true)
}
