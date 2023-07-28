package migrator_test

import (
	"database/sql"
	"os"
)

func initMemorySqlite() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	return db
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

func resetJsonFile() error {

	return os.Remove(testFixtureFolder + "/migrations.json")
}

func initFolder(fullPath string) error {
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
