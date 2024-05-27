package migrator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
)

type migration struct {
	db                *sql.DB
	migrationProvider MigrationProvider
	migrationFilePath string
}

const (
	migrationFileRegex   = "^.*migrate.*\\.sql$"
	rollbackReplaceRegex = "migrate"
)

func newMigrator(db *sql.DB) *migration {
	return &migration{db: db}
}

func (m *migration) orderedMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(m.migrationFilePath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if m.isMigration(file.Name()) {
			fileNames = append(fileNames, file.Name())
		}
	}

	sort.Strings(fileNames)

	return fileNames, nil
}

func (m *migration) executeSqlFile(fileName string) (bool, error) {
	if m.migrationProvider.MigrationExistsForFile(fileName) {
		return false, nil
	}

	fmt.Printf("Running migration '%s'\n", fileName)
	content, err := ioutil.ReadFile(m.migrationFilePath + "/" + fileName)
	if err != nil {
		return false, err
	}

	err = m.executeSql(string(content))
	if err == nil {
		err = m.migrationProvider.AddToMigration(fileName)
		if err != nil {
			return false, err
		}
	}

	m.migrationProvider.AddToMigrationReport(fileName, err)

	return true, err
}

func (m *migration) executeRollbackSqlFile(fileName string) error {
	rollbackFileName, err := m.resolveRollbackFile(fileName)
	if err != nil {
		fmt.Printf("Skip rollback for %s as rollback file does not exists\n", fileName)
		m.migrationProvider.RemoveFromMigration(fileName)

		return nil
	}

	fmt.Printf("Running rollback '%s'\n", rollbackFileName)
	content, err := ioutil.ReadFile(m.migrationFilePath + "/" + rollbackFileName)
	if err != nil {
		return err
	}

	err = m.executeSql(string(content))
	if err == nil {
		m.migrationProvider.RemoveFromMigration(fileName)
	}

	return err
}

func (m *migration) executeSql(sql string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(sql)

	return err
}

func (m *migration) isMigration(fileName string) bool {
	regex := regexp.MustCompile(migrationFileRegex)

	mathces := regex.FindStringSubmatch(fileName)

	return len(mathces) > 0
}

func (m *migration) resolveRollbackFile(migrationFileName string) (string, error) {
	regex := regexp.MustCompile(rollbackReplaceRegex)

	result := regex.ReplaceAllStringFunc(migrationFileName, func(match string) string {
		if match == rollbackReplaceRegex {
			return "rollback"
		}
		return "unknown"
	})

	if !fileExists(m.migrationFilePath + "/" + result) {
		return "", fmt.Errorf("file does not %s exists", result)
	}

	return result, nil
}
