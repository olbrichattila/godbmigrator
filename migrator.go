package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

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
	m.migrationProvider.resetDate()
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
	migrations, err := m.migrationProvider.migrations(!isCompleteRollback)
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
