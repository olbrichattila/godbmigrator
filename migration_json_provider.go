package migrator

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"time"
)

const migragionJsonFileName = "./migrations/migrations.json"

type JsonMigration struct {
	data         map[string]string
	timeString   string
	jsonFileName string
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
		if dateString == latestDate || isLatest == false {
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
		return migragionJsonFileName
	}

	return m.jsonFileName
}

func (m *JsonMigration) SetJsonFileName(fileName string) {
	m.jsonFileName = fileName + "/migrations.json"
}
