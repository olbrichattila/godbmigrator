// Package migrationfile manages migration file creation
package migrationfile

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/olbrichattila/godbmigrator/internal/helper"
)

// Manager encapsulates the migration file management methods
type Manager interface {
	CreateNewMigrationFiles(migrationFilePath, customText string, isRollback bool) error
	ResolveRollbackFile(migrationFileName string) (string, error)
	OrderedMigrationFiles() ([]string, error)
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

func (m *mFile) ResolveRollbackFile(migrationFileName string) (string, error) {
	regex := regexp.MustCompile(rollbackReplaceRegex)

	result := regex.ReplaceAllStringFunc(migrationFileName, func(match string) string {
		if match == rollbackReplaceRegex {
			return typeRollback
		}
		return "unknown"
	})

	if !helper.FileExists(m.migrationFilePath + "/" + result) {
		return "", fmt.Errorf("file does not %s exists", result)
	}

	return result, nil
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
	regex := regexp.MustCompile(migrationFileRegex)
	matches := regex.FindStringSubmatch(fileName)

	return len(matches) > 0
}
