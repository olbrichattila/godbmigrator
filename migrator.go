package migrator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"time"
)

const (
	migrationFileRegex   = "^.*migrate.*\\.sql$"
	rollbackReplaceRegex = "migrate"
)

type migration struct {
	db                *sql.DB
	migrationProvider MigrationProvider
	migrationFilePath string
}

func newMigrator(db *sql.DB) *migration {
	return &migration{db: db}
}

func Rollback(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	return rollback(db, migrationProvider, migrationFilePath, count, false)
}

func Refresh(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
) error {
	err := rollback(db, migrationProvider, migrationFilePath, 0, true)
	if err != nil {
		return err
	}

	return Migrate(db, migrationProvider, migrationFilePath, 0)
}

func rollback(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
	isCompleteRollback bool,
) error {
	var err error
	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	m.migrationProvider.SetJsonFilePath(migrationFilePath)
	migrations, err := m.migrationProvider.Migrations(!isCompleteRollback)
	if err != nil {
		return err
	}
	if len(migrations) == 0 {
		fmt.Println("Nothing to rollback")
		return nil
	}

	rollbackCount := 0
	for _, fileName := range migrations {
		if count > 0 {
			if rollbackCount == count {
				break
			}

		}

		err = m.executeRollbackSqlFile(fileName)
		if err != nil {
			return err
		}
		rollbackCount++
	}

	fmt.Printf("Rolled back %d items\n", rollbackCount)

	return nil
}

func Migrate(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	m.migrationProvider.SetJsonFilePath(migrationFilePath)
	m.migrationProvider.ResetDate()
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

func Report(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
) (string, error) {

	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	m.migrationProvider.SetJsonFilePath(migrationFilePath)
	return m.migrationProvider.Report()
}

func AddNewMigrationFiles(migrationFilePath, customText string) error {
	var err error
	err = createNewMigrationFiles(migrationFilePath, customText, false)
	if err != nil {
		return err
	}
	err = createNewMigrationFiles(migrationFilePath, customText, true)
	if err != nil {
		return err
	}

	return nil
}

func createNewMigrationFiles(migrationFilePath, customText string, isRollback bool) error {
	alteredCustomText := customText
	mgType := "migrate"

	if customText != "" {
		alteredCustomText = "-" + customText
	}

	if isRollback {
		mgType = "rollback"
	}

	fileName := fmt.Sprintf(
		"%s-%s%s.sql",
		time.Now().Format("2006-01-02_15_04_05"),
		mgType,
		alteredCustomText,
	)

	filePath := migrationFilePath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("Migration file %s created\n", filePath)

	return nil
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
