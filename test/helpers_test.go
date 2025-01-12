package migrator_test

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	if err != nil {
		return 0, err
	}

	return count, err
}

func rowCountInTable(db *sql.DB, tableName string) (int, error) {
	query := "SELECT count(*) from " + tableName

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getChecksumFromTable(db *sql.DB, fileName string) (string, error) {
	query := fmt.Sprintf("SELECT checksum from %s_migrations WHERE file_name=?", tablePrefix)

	var checksum string
	err := db.QueryRow(query, fileName).Scan(&checksum)
	if err != nil {
		return "", err
	}

	return checksum, nil
}


func resetJsonFile() error {
	return os.Remove(testFixtureFolder + "/migrations.json")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func haveReportRecord(db *sql.DB, fileName, createdAt, status, message string) error {
	sql := `INSERT INTO %s_migration_reports
			(file_name, created_at, result_status, message)
			VALUES (?,?,?,?)`
	sql = fmt.Sprintf(sql, tablePrefix)

	_, err := db.Exec(sql, fileName, createdAt, status, message)

	return err
}

func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}
	
	return hex.EncodeToString(hash.Sum(nil)), nil
}
