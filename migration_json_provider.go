package migrator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

const migrationJSONFileName = "./migrations/migrations.json"
const migrationJSONReportFileName = "./migrations/migration_report.json"

type jsonMigration struct {
	data              map[string]string
	timeString        string
	jsonFileName      string
	jsonReporFileName string
}

type jsonMigrationReport struct {
	FileName     string `json:"fileName"`
	CreatedAt    string `json:"createdAt"`
	ResultStatus string `json:"resultStatus"`
	Message      string `json:"message"`
}

func newJSONMigration() (*jsonMigration, error) {
	jsonMigration := &jsonMigration{}
	jsonMigration.resetDate()
	err := jsonMigration.loadMigrationFile()

	return jsonMigration, err
}

func (m *jsonMigration) resetDate() {
	m.timeString = time.Now().Format("2006-01-02 15:04:05")
}

func (m *jsonMigration) migrations(isLatest bool) ([]string, error) {
	var latestDate string
	var filtered []string

	for _, dateString := range m.data {
		if latestDate == "" || dateString > latestDate {
			latestDate = dateString
		}
	}

	for fileName, dateString := range m.data {
		if dateString == latestDate || !isLatest {
			filtered = append(filtered, fileName)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(filtered)))

	return filtered, nil
}

func (m *jsonMigration) loadMigrationFile() error {
	m.data = make(map[string]string)
	jsonFileName := m.getJSONFileName()
	if !fileExists(jsonFileName) {
		return nil
	}

	jsonData, err := ioutil.ReadFile(jsonFileName)
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

	jsonFileName := m.getJSONFileName()
	return ioutil.WriteFile(jsonFileName, jsonData, 0644)
}

func (m *jsonMigration) addToMigration(fileName string) error {
	m.data[fileName] = m.timeString

	return m.saveMigrationFile()
}

func (m *jsonMigration) removeFromMigration(fileName string) error {
	delete(m.data, fileName)

	return m.saveMigrationFile()
}

func (m *jsonMigration) migrationExistsForFile(fileName string) (bool, error) {
	return m.data[fileName] != "", nil
}

func (m *jsonMigration) getJSONFileName() string {
	if m.jsonFileName == "" {
		return migrationJSONFileName
	}

	return m.jsonFileName
}

func (m *jsonMigration) getJSONReportFileName() string {
	if m.jsonFileName == "" {
		return migrationJSONReportFileName
	}

	return m.jsonReporFileName
}

func (m *jsonMigration) SetJSONFilePath(filePath string) {
	m.jsonFileName = filePath + "/migrations.json"
	m.jsonReporFileName = filePath + "/migration_reports.json"
}

func (m *jsonMigration) AddToMigrationReport(fileName string, errorToLog error) error {
	storeFileName := m.getJSONReportFileName()
	message := "ok"
	status := "success"
	if errorToLog != nil {
		message = errorToLog.Error()
		status = "error"
	}

	file, err := os.OpenFile(storeFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	newReportItem := jsonMigrationReport{
		FileName:     fileName,
		ResultStatus: status,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
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
	byteValue, err := ioutil.ReadAll(jsonFile)
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
			"Created at: %s, File Name: %s, Status: %s, Message: %s\n",
			item.CreatedAt,
			item.FileName,
			item.ResultStatus,
			item.Message,
		)
		builder.WriteString(str)
	}

	return builder.String(), nil
}
