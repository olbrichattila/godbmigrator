package migrator

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

const migrationJSONFileName = "./migrations/migrations.json"
const migrationJSONReportFileName = "./migrations/migration_report.json"

type jsonMigration struct {
	data              map[string]MigrationRow
	timeString        string
	jsonFileName      string
	jsonReportFileName string
}

type jsonMigrationReport struct {
	FileName     string `json:"fileName"`
	CreatedAt    string `json:"createdAt"`
	ResultStatus string `json:"resultStatus"`
	Message      string `json:"message"`
}

func newJSONMigration() (*jsonMigration, error) {
	jsonMigration := &jsonMigration{}
	jsonMigration.ResetDate()
	err := jsonMigration.loadMigrationFile()

	return jsonMigration, err
}

func (m *jsonMigration) ResetDate() {
	m.timeString = time.Now().Format(timeFormat)
}

func (m *jsonMigration) Migrations(isLatest bool) ([]MigrationRow, error) {
	var latestDate string
	var filtered []MigrationRow

	for _, dateString := range m.data {
		if latestDate == "" || dateString.Migration > latestDate {
			latestDate = dateString.Migration
		}
	}

	for fileName, dateString := range m.data {
		if dateString.Migration == latestDate || !isLatest {
			filtered = append(filtered, MigrationRow{Migration: fileName})
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Migration > filtered[j].Migration
	})

	return filtered, nil
}

func (m *jsonMigration) loadMigrationFile() error {
	m.data = make(map[string]MigrationRow)
	jsonFileName := m.GetJSONFileName()
	if !fileExists(jsonFileName) {
		return nil
	}

	jsonData, err := os.ReadFile(jsonFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &m.data)
	if err != nil {
		return err
	}

	return nil
}

func (m *jsonMigration) saveMigrationFile() error {
	jsonData, err := json.Marshal(&m.data)
	if err != nil {
		return err
	}

	jsonFileName := m.GetJSONFileName()
	return os.WriteFile(jsonFileName, jsonData, 0644)
}

func (m *jsonMigration) AddToMigration(fileName, checksum string) error {
	migration := MigrationRow{Migration: m.timeString, Checksum: checksum}
	m.data[fileName] = migration

	return m.saveMigrationFile()
}

func (m *jsonMigration) RemoveFromMigration(fileName string) error {
	delete(m.data, fileName)

	return m.saveMigrationFile()
}

func (m *jsonMigration) MigrationExistsForFile(fileName string) (bool, error) {
	return m.data[fileName].Migration != "", nil
}

func (m *jsonMigration) GetJSONFileName() string {
	if m.jsonFileName == "" {
		return migrationJSONFileName
	}

	return m.jsonFileName
}

func (m *jsonMigration) getJSONReportFileName() string {
	if m.jsonFileName == "" {
		return migrationJSONReportFileName
	}

	return m.jsonReportFileName
}

func (m *jsonMigration) SetJSONFilePath(filePath string) {
	m.jsonFileName = filePath + "/migrations.json"
	m.jsonReportFileName = filePath + "/migration_reports.json"
}

func (m *jsonMigration) AddToMigrationReport(fileName string, errorToLog error) error {
	storeFileName := m.getJSONReportFileName()
	message := "ok"
	status := statusSuccess
	if errorToLog != nil {
		message = errorToLog.Error()
		status = statusError
	}

	file, err := os.OpenFile(storeFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	newReportItem := jsonMigrationReport{
		FileName:     fileName,
		ResultStatus: status,
		CreatedAt:    time.Now().Format(timeFormat),
		Message:      message,
	}

	newData, err := json.Marshal(newReportItem)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > 0 {
		newData = append([]byte(","), newData...)
	}

	if _, err := file.Write(newData); err != nil {
		return err
	}

	return nil
}

func (m *jsonMigration) Report() (string, error) {
	storeFileName := m.getJSONReportFileName()

	_, err := os.Stat(storeFileName)
	if os.IsNotExist(err) {
		return "", nil
	}

	// Open the JSON file
	jsonFile, err := os.Open(storeFileName)
	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	// Read the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	byteValue = append([]byte{'['}, byteValue...)
	byteValue = append(byteValue, ']')

	if err != nil {
		return "", err
	}

	var collection []jsonMigrationReport
	err = json.Unmarshal(byteValue, &collection)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, item := range collection {
		str := fmt.Sprintf(
			reportMessageText,
			item.CreatedAt,
			item.FileName,
			item.ResultStatus,
			item.Message,
		)
		builder.WriteString(str)
	}

	return builder.String(), nil
}
