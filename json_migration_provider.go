package migrator

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"time"
)

type JsonMigration struct {
	data       map[string]string
	timeString string
}

func newJsonMigration() *JsonMigration {
	return &JsonMigration{}
}

func (m *JsonMigration) LatestMigrations() []string {
	m.loadMigrationFile()
	var latestDate string
	var filtered []string

	for _, dateString := range m.data {
		if latestDate == "" || dateString > latestDate {
			latestDate = dateString
		}
	}

	for fileName, dateString := range m.data {
		if dateString == latestDate {
			filtered = append(filtered, fileName)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(filtered)))

	return filtered
}

func (m *JsonMigration) loadMigrationFile() error {
	m.timeString = time.Now().Format("2006-01-02 15:04:05")
	m.data = make(map[string]string)
	jsonData, err := ioutil.ReadFile(migragionJsonFileName)
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

	return ioutil.WriteFile(migragionJsonFileName, jsonData, 0644)
}

func (m *JsonMigration) AddToMigration(fileName string) error {
	if len(m.data) == 0 {
		m.loadMigrationFile()
	}
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
