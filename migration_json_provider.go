package migrator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

const migrationJsonFileName = "./migrations/migrations.json"
const migrationJsonReportFileName = "./migrations/migration_report.json"

type JsonMigration struct {
	data              map[string]string
	timeString        string
	jsonFileName      string
	jsonReporFileName string
}

type JsonMigrationReport struct {
	FileName  string `json:"fileName`
	CreatedAt string `json:"createdAt`
	Message   string `json:message`
}

func newJsonMigration() *JsonMigration {
	jsonMigration := &JsonMigration{}
	jsonMigration.ResetDate()
	jsonMigration.loadMigrationFile()

	return jsonMigration
}

func (m *JsonMigration) ResetDate() {
	m.timeString = time.Now().Format("2006-01-02 15:04:05")
}

func (m *JsonMigration) Migrations(isLatest bool) ([]string, error) {
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

func (m *JsonMigration) loadMigrationFile() error {
	m.data = make(map[string]string)
	jsonFileName := m.GetJsonFileName()
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

func (m *JsonMigration) saveMigrationFile() error {
	jsonData, err := json.Marshal(&m.data)
	if err != nil {
		return err
	}

	jsonFileName := m.GetJsonFileName()
	return ioutil.WriteFile(jsonFileName, jsonData, 0644)
}

func (m *JsonMigration) AddToMigration(fileName string) error {
	m.data[fileName] = m.timeString

	return m.saveMigrationFile()
}

func (m *JsonMigration) RemoveFromMigration(fileName string) error {
	delete(m.data, fileName)

	return m.saveMigrationFile()
}

func (m *JsonMigration) MigrationExistsForFile(fileName string) bool {
	return m.data[fileName] != ""
}

func (m *JsonMigration) GetJsonFileName() string {
	if m.jsonFileName == "" {
		return migrationJsonFileName
	}

	return m.jsonFileName
}

func (m *JsonMigration) getJsonReportFileName() string {
	if m.jsonFileName == "" {
		return migrationJsonReportFileName
	}

	return m.jsonReporFileName
}

func (m *JsonMigration) SetJsonFileName(filePath string) {
	m.jsonFileName = filePath + "/migrations.json"
	m.jsonReporFileName = filePath + "/migration_reports.json"
}

func (m *JsonMigration) AddToMigrationReport(fileName, message string) error {
	storeFileName := m.getJsonReportFileName()
	file, err := os.OpenFile(storeFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	newReportItem := JsonMigrationReport{FileName: fileName, CreatedAt: time.Now().Format("2006-01-02 15:04:05"), Message: message}
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
