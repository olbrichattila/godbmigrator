package migrator_test

import (
	"io/ioutil"
	"os"
	"strings"

	migrator "github.com/olbrichattila/godbmigrator"
)

func (t *TestSuite) TestMigrationAdded() {
	resetrTestMigrationPath(testMigrationFilePath)
	err := migrator.AddNewMigrationFiles(testMigrationFilePath, "")
	t.Nil(err)

	count, err := countFilesInDirectory(testMigrationFilePath)
	t.Nil(err)

	t.Equal(2, count)
}

func (t *TestSuite) TestMigrationAddedWithCustomName() {
	resetrTestMigrationPath(testMigrationFilePath)
	customText := "custom-text"

	err := migrator.AddNewMigrationFiles(testMigrationFilePath, customText)
	t.Nil(err)

	count, err := countFilesInDirectory(testMigrationFilePath)
	t.Nil(err)

	t.Equal(2, count)

	exists, err := checkStringInFileNames(testMigrationFilePath, customText)
	t.Nil(err)

	t.True(exists)
}

func resetrTestMigrationPath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func countFilesInDirectory(path string) (int, error) {
	files, err := ioutil.ReadDir(path)
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
	files, err := ioutil.ReadDir(dirPath)
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
