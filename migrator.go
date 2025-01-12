// Package migrator is a lightweight database migrator package, pass only the *sql.DB, migrate, rollback and migration reports
package migrator

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// Rollback rolls back last migrated items or all if count is 0
func Rollback(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	return rollback(db, migrationProvider, migrationFilePath, count, false)
}

// Refresh runs a full rollback and migrate again
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

// Migrate execute migrations
func Migrate(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.MigrationProvider = migrationProvider
	m.MigrationProvider.SetJSONFilePath(migrationFilePath)
	m.MigrationProvider.ResetDate()
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
		migrated, err := m.executeSQLFile(fileName)
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

// Report return a report of the already executed migrations
func Report(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
) (string, error) {

	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.MigrationProvider = migrationProvider
	m.MigrationProvider.SetJSONFilePath(migrationFilePath)
	return m.MigrationProvider.Report()
}

// AddNewMigrationFiles adds a new blank migration file and a rollback file
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

// ChecksumValidation validates if the checksums are correct and nothing changed
func ChecksumValidation(
	db *sql.DB,
	migrationProvider MigrationProvider,
	migrationFilePath string,
) []string {
	errors := make([]string, 0)
	m := newMigrator(db)
	m.migrationFilePath = migrationFilePath
	m.MigrationProvider = migrationProvider
	m.MigrationProvider.SetJSONFilePath(migrationFilePath)
	migrations, err := m.MigrationProvider.Migrations(false)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	for _, mig := range migrations {
		filePath := m.migrationFilePath + "/" + mig.Migration
		if !fileExists(filePath) {
			errors = append(errors, fmt.Sprintf("migration file for checksum does not %s exists", mig.Migration))
			continue
		}

		md5, err := calculateFileMD5(m.migrationFilePath + "/" + mig.Migration)
		if err != nil {
			errors = append(errors, fmt.Sprintf("migration file for checksum could not be opened %s exists", mig.Migration))
			continue
		}

		if md5 != mig.Checksum {
			errors = append(errors, fmt.Sprintf("md5 error for file %s, md5 %s/%s", mig.Migration, md5, mig.Checksum))
		}
	}

	return errors
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
	m.MigrationProvider = migrationProvider
	m.MigrationProvider.SetJSONFilePath(migrationFilePath)
	migrations, err := m.MigrationProvider.Migrations(!isCompleteRollback)
	if err != nil {
		return err
	}
	if len(migrations) == 0 {
		fmt.Println("Nothing to rollback")
		return nil
	}

	rollbackCount := 0
	for _, mig := range migrations {
		if count > 0 {
			if rollbackCount == count {
				break
			}
		}

		err = m.executeRollbackSQLFile(mig.Migration)
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
		mgType = typeRollback
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

func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}
	
	return hex.EncodeToString(hash.Sum(nil)), nil
}