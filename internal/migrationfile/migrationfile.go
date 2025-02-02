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

const (
	sqlFileExt = ".sql"
)

// Manager encapsulates the migration file management methods
type Manager interface {
	CreateNewMigrationFiles(migrationFilePath, customText string) ([]string, error)
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
func (m *mFile) CreateNewMigrationFiles(migrationFilePath, customText string) ([]string, error) {
	datePart := time.Now().Format("2006-01-02_15_04_05")
	file1, err := m.createNewMigrationFile(migrationFilePath, customText, datePart, false)
	if err != nil {
		return nil, err
	}

	file2, err := m.createNewMigrationFile(migrationFilePath, customText, datePart, true)
	if err != nil {
		return nil, err
	}

	return []string{file1, file2}, nil

}

func (*mFile) createNewMigrationFile(migrationFilePath, customText, datePart string, isRollback bool) (string, error) {
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
		return "", err
	}
	defer file.Close()

	return filePath, nil
}

func (m *mFile) ResolveRollbackFile(migrationFileName string) (string, error) {

	lastIndex := strings.LastIndex(migrationFileName, sqlFileExt)
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
		fileName := file.Name()
		if strings.HasSuffix(fileName, sqlFileExt) && m.isMigration(fileName) {
			fileNames = append(fileNames, fileName)
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
