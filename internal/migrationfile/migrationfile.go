// Package migrationfile manages migration file creation
package migrationfile

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

// Manager encapsulates the migration file management methods
type Manager interface {
	CreateNewMigrationFiles(migrationFilePath, customText string, isRollback bool) error
	IsMigration(fileName string) bool
	ResolveRollbackFile(migrationFileName string) (string, error)
}

const (
	typeRollback         = "rollback"
	migrationFileRegex   = "^.*migrate.*\\.sql$"
	rollbackReplaceRegex = "migrate"
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
func (*mFile) CreateNewMigrationFiles(migrationFilePath, customText string, isRollback bool) error {
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

func (m *mFile) IsMigration(fileName string) bool {
	regex := regexp.MustCompile(migrationFileRegex)
	matches := regex.FindStringSubmatch(fileName)

	return len(matches) > 0
}

func (m *mFile) ResolveRollbackFile(migrationFileName string) (string, error) {
	regex := regexp.MustCompile(rollbackReplaceRegex)

	result := regex.ReplaceAllStringFunc(migrationFileName, func(match string) string {
		if match == rollbackReplaceRegex {
			return typeRollback
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
