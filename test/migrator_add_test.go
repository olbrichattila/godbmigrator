package migrator_test

import (
	"os"
	"strings"
	"testing"

	migrator "github.com/olbrichattila/godbmigrator"
	"github.com/stretchr/testify/suite"
)

const testMigrationFilePath = "./test-migrations"

type AddTestSuite struct {
	suite.Suite
	migrator migrator.DBMigrator
}

func TestAddTestRunner(t *testing.T) {
	suite.Run(t, new(AddTestSuite))
}

func (suite *AddTestSuite) SetupTest() {
	resetTestMigrationPath()
	db := initMemorySqlite()
	suite.migrator = migrator.New(db, testMigrationFilePath, tablePrefix)
}

func (t *AddTestSuite) TestMigrationAdded() {
	err := t.migrator.AddNewMigrationFiles("")
	t.Nil(err)

	count, err := countFilesInDirectory(testMigrationFilePath)
	t.Nil(err)

	t.Equal(2, count)
}

func (t *AddTestSuite) TestMigrationAddedWithCustomName() {
	customText := "custom-text"

	err := t.migrator.AddNewMigrationFiles(customText)
	t.Nil(err)

	count, err := countFilesInDirectory(testMigrationFilePath)
	t.Nil(err)

	t.Equal(2, count)

	exists, err := checkStringInFileNames(testMigrationFilePath, customText)
	t.Nil(err)

	t.True(exists)
}

func resetTestMigrationPath() error {
	err := os.RemoveAll(testMigrationFilePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(testMigrationFilePath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func countFilesInDirectory(path string) (int, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}

	return count, nil
}

func checkStringInFileNames(dirPath, searchString string) (bool, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if !strings.Contains(string(file.Name()), searchString) {
			return false, nil
		}
	}

	return true, nil
}
