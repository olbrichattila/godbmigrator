// Package migrate is the main internal package of the migrator
package migrate

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olbrichattila/godbmigrator/internal/helper"
	"github.com/olbrichattila/godbmigrator/internal/migrationfile"
	"github.com/olbrichattila/godbmigrator/messager"
)

// Migrator abstracts migration logic
type Migrator interface {
	Migrate(migrationProvider MigrationProvider, migrationFilePath string, count int) error
	Rollback(migrationProvider MigrationProvider, migrationFilePath string, count int, isCompleteRollback bool) error
	Report(migrationProvider MigrationProvider, migrationFilePath string) (string, error)
	ChecksumValidation(migrationProvider MigrationProvider, migrationFilePath string) []string
}

type migration struct {
	db                   *sql.DB
	migrationProvider    MigrationProvider
	migrationFilePath    string
	migrationFileManager migrationfile.Manager
	msg                  messager.Messager
}

// New creates a new migration
func New(db *sql.DB, migrationFileManager migrationfile.Manager, msg messager.Messager) Migrator {
	return &migration{
		db:                   db,
		migrationFileManager: migrationFileManager,
		msg:                  msg,
	}
}

func (m *migration) Migrate(
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
) error {
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	m.migrationProvider.ResetDate()

	fileNames, err := m.migrationFileManager.OrderedMigrationFiles()
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

	m.messageDispatch(messager.MigratedItems, strconv.Itoa(migrateCount))

	return nil
}

func (m *migration) Rollback(
	migrationProvider MigrationProvider,
	migrationFilePath string,
	count int,
	isCompleteRollback bool,
) error {
	var err error

	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	migrations, err := m.migrationProvider.Migrations(!isCompleteRollback)
	if err != nil {
		return err
	}
	if len(migrations) == 0 {
		m.messageDispatch(messager.NothingToRollback, "")
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

	m.messageDispatch(messager.RolledBack, strconv.Itoa(rollbackCount))

	return nil
}

func (m *migration) Report(
	migrationProvider MigrationProvider,
	migrationFilePath string,
) (string, error) {
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider

	return m.migrationProvider.Report()
}

func (m *migration) ChecksumValidation(
	migrationProvider MigrationProvider,
	migrationFilePath string,
) []string {
	errors := make([]string, 0)
	m.migrationFilePath = migrationFilePath
	m.migrationProvider = migrationProvider
	migrations, err := m.migrationProvider.Migrations(false)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	for _, mig := range migrations {
		filePath := m.migrationFilePath + "/" + mig.Migration
		if !helper.FileExists(filePath) {
			errors = append(errors, fmt.Sprintf("migration file for checksum does not %s exists", mig.Migration))
			continue
		}

		md5, err := helper.CalculateFileMD5(m.migrationFilePath + "/" + mig.Migration)
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

func (m *migration) executeSQLFile(fileName string) (bool, error) {
	exists, err := m.migrationProvider.MigrationExistsForFile(fileName)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	m.messageDispatch(messager.RunningMigrations, fileName)
	content, err := os.ReadFile(m.migrationFilePath + "/" + fileName)
	if err != nil {
		return false, err
	}

	contentString := string(content)
	err = m.executeSQL(contentString)
	if err == nil {
		hash := m.getHash(contentString)
		err = m.migrationProvider.AddToMigration(fileName, hash)
		if err != nil {
			return false, err
		}
	}

	_ = m.migrationProvider.AddToMigrationReport(fileName, err)

	return true, err
}

func (m *migration) executeRollbackSQLFile(fileName string) error {
	rollbackFileName, err := m.migrationFileManager.ResolveRollbackFile(fileName)
	if err != nil {
		m.messageDispatch(messager.SkipRollback, fileName)

		err := m.migrationProvider.RemoveFromMigration(fileName)
		if err != nil {
			return err
		}

		return nil
	}

	m.messageDispatch(messager.RunningRollback, rollbackFileName)

	content, err := os.ReadFile(m.migrationFilePath + "/" + rollbackFileName)
	if err != nil {
		return err
	}

	err = m.executeSQL(string(content))
	if err == nil {
		err = m.migrationProvider.RemoveFromMigration(fileName)
		if err != nil {
			return err
		}
	}

	_ = m.migrationProvider.AddToMigrationReport(rollbackFileName, err)

	return err
}

func (m *migration) executeSQL(sql string) error {
	statements := m.splitSQLStatements(sql)
	for _, singleSQL := range statements {
		err := m.executeSingleSQL(singleSQL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *migration) splitSQLStatements(sqlScript string) []string {
	var statements []string
	var currentStatement bytes.Buffer
	inProcedure := false

	lines := strings.Split(sqlScript, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToUpper(trimmed), "CREATE PROCEDURE") || strings.HasPrefix(strings.ToUpper(trimmed), "CREATE FUNCTION") {
			inProcedure = true
		}

		currentStatement.WriteString(line + "\n")
		if strings.HasSuffix(trimmed, ";") && !inProcedure {
			statements = append(statements, currentStatement.String())
			currentStatement.Reset()
		}

		if inProcedure && strings.HasPrefix(strings.ToUpper(trimmed), "END") && strings.HasSuffix(trimmed, ";") {
			inProcedure = false
			statements = append(statements, currentStatement.String())
			currentStatement.Reset()
		}
	}

	if currentStatement.Len() > 0 {
		statements = append(statements, currentStatement.String())
	}

	return statements
}

func (m *migration) executeSingleSQL(sql string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(sql)

	return err
}

func (m *migration) getHash(sql string) string {
	hash := md5.Sum([]byte(sql))
	return hex.EncodeToString(hash[:])
}

func (m *migration) messageDispatch(eventType int, message string) {
	if m.msg != nil {
		m.msg.Dispatch(eventType, message)
	}
}
