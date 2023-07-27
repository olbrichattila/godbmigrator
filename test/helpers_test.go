package migrator_test

import (
	"database/sql"
	"os"
)

func inMemorySqlite() (*sql.DB, error) {
	// return sql.Open("sqlite3", ":memory:")
	return sql.Open("sqlite3", "./data/data2.sqlite")
}

func tableCountInDatabase(db *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table'"

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func resetDatabase() error {
	return os.Remove("./data/data2.sqlite")
}

func resetJsonFile() error {
	return os.Remove("./test_fixtures/migrations.json")
}

func initFolder(fullPath string) error {
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
