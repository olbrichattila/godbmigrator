package migrator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

const (
	migrationFolder      = "./migrations"
	migrationFileRegex   = "^.*migrate.*\\.sql$"
	rollbackReplaceRegex = "migrate"
)

type migration struct {
	db                *sql.DB
	migrationProvider MigrationProvider
}

func newMigrator(db *sql.DB) *migration {
	return &migration{db: db}
}

func Rollback(db *sql.DB, migrationProvider MigrationProvider, count int) error {
	m := newMigrator(db)
	m.migrationProvider = migrationProvider
	latestMigrations := m.migrationProvider.LatestMigrations()
	if len(latestMigrations) == 0 {
		fmt.Println("Nothing to rollback")
		return nil
	}

	rollbackCount := 0
	for _, fileName := range latestMigrations {
		if count > 0 {
			if rollbackCount == count {
				break
			}

		}

		err := m.executeRollbackSqlFile(fileName)
		if err != nil {
			return err
		}
		rollbackCount++
	}

	fmt.Printf("Rolled back %d items\n", rollbackCount)

	return nil
}

func Migrate(db *sql.DB, migrationProvider MigrationProvider, count int) error {
	m := newMigrator(db)
	m.migrationProvider = migrationProvider

	fileNames, err := m.orderedMigrationFiles()
	if err != nil {
		return err
	}

	migrateCount := 0
	for _, fileName := range fileNames {
		if count > 0 {
			if migrateCount == count {
				break
			}
		}
		migrated, err := m.executeSqlFile(fileName)
		if err != nil {
			return err
		}

		if migrated {
			migrateCount++
		}
	}

	fmt.Printf("Migrated %d items\n", migrateCount)

	return nil
}

func (m *migration) orderedMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(migrationFolder)
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
	content, err := ioutil.ReadFile(migrationFolder + "/" + fileName)
	if err != nil {
		return false, err
	}

	err = m.executeSql(string(content))
	if err == nil {
		m.migrationProvider.AddToMigration(fileName)
	}

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
	content, err := ioutil.ReadFile(migrationFolder + "/" + rollbackFileName)
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
		tx.Commit()
	}()

	_, err = m.db.Exec(sql)

	return err
}

func (m *migration) isMigration(fileName string) bool {
	regex, err := regexp.Compile(migrationFileRegex)
	if err != nil {
		fmt.Println("Regex pattern error: " + err.Error())
	}

	mathces := regex.FindStringSubmatch(fileName)

	return len(mathces) > 0
}

func (m *migration) resolveRollbackFile(migrationFileName string) (string, error) {
	regex, err := regexp.Compile(rollbackReplaceRegex)
	if err != nil {
		return "", err
	}

	result := regex.ReplaceAllStringFunc(migrationFileName, func(match string) string {
		if match == rollbackReplaceRegex {
			return "rollback"
		}
		return "unknown"
	})

	if !fileExists(migrationFolder + "/" + result) {
		return "", fmt.Errorf("File does not %s exists", result)
	}

	return result, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
