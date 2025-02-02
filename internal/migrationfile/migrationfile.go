// Package migrationfile manages migration file creation
package migrationfile

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/olbrichattila/godbmigrator/internal/helper"
)

// Manager encapsulates the migration file management methods
type Manager interface {
	CreateNewMigrationFiles(migrationFilePath, customText string) error
	ResolveRollbackFile(migrationFileName string) (string, error)
	OrderedMigrationFiles() ([]string, error)
}

const (
	typeRollback          = "rollback"
	nonMigrationFileRegex = "^.*-rollback\\.sql$"
	rollbackReplaceRegex  = "^.*\\.sql$"
)

// New returns with a new file manager instance
func New(migrationFilePath string) Manager {
	return &mFile{
		migrationFilePath: migrationFilePath,
	}
}

type mFile struct {
	migrationFilePath string
}

// CreateNewMigrationFiles responsible for creating migration files
func (m *mFile) CreateNewMigrationFiles(migrationFilePath, customText string) error {
	datePart := time.Now().Format("2006-01-02_15_04_05")
	err := m.createNewMigrationFile(migrationFilePath, customText, datePart, false)
	if err != nil {
		return err
	}

	return m.createNewMigrationFile(migrationFilePath, customText, datePart, true)
}

func (*mFile) createNewMigrationFile(migrationFilePath, customText, datePart string, isRollback bool) error {
	suffix := ""
	prefix := ""

	if customText != "" {
		prefix = "-" + customText
	}

	if isRollback {
		suffix = "-" + typeRollback
	}

	migrationfileName := fmt.Sprintf(
		"%s%s%s.sql",
		datePart,
		prefix,
		suffix,
	)

	filePath := migrationFilePath + "/" + migrationfileName
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("Migration file %s created\n", filePath)

	return nil
}

func (m *mFile) ResolveRollbackFile(migrationFileName string) (string, error) {

	lastIndex := strings.LastIndex(migrationFileName, ".sql")
	if lastIndex == -1 {
		return "", fmt.Errorf("non sql file provided for rollback %s exists", migrationFileName)
	}

	rollbackFile := migrationFileName[:lastIndex] + "-rollback.sql"
	if !helper.FileExists(m.migrationFilePath + "/" + rollbackFile) {
		return "", fmt.Errorf("file does not %s exists", rollbackFile)
	}

	return rollbackFile, nil
}

func (m *mFile) OrderedMigrationFiles() ([]string, error) {
	files, err := os.ReadDir(m.migrationFilePath)
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

func (m *mFile) isMigration(fileName string) bool {
	if fileName == "baseline.sql" {
		return false
	}
	regex := regexp.MustCompile(nonMigrationFileRegex)
	matches := regex.FindStringSubmatch(fileName)

	return len(matches) == 0
}
