package migrator_test

import (
	"database/sql"
	"fmt"
	"io"
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

	return count, err
}

func rowCountInTable(db *sql.DB, tableName string) (int, error) {
	query := "SELECT count(*) from " + tableName

	var count int
	err := db.QueryRow(query).Scan(&count)

	return count, err
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

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file, truncating it if it already exists
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Flush the file content to disk
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func haveReportrecord(db *sql.DB, fileName, createdAt, status, message string) error {
	sql := fmt.Sprintf(`INSERT INTO migration_reports
			(file_name, created_at, result_status, message)
			VALUES (?,?,?,?)`,
	)

	_, err := db.Exec(sql, fileName, createdAt, status, message)

	return err
}
