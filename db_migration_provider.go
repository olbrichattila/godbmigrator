package migrator

import (
	"database/sql"
	"strconv"
	"time"
)

type DbMigration struct {
	db         *sql.DB
	timeString string
}

func newDbMigration(db *sql.DB) *DbMigration {
	timeString := time.Now().Format("2006-01-02 15:04:05")
	dbMigration := &DbMigration{db: db, timeString: timeString}
	err := dbMigration.Init()
	if err != nil {
		panic(err)
	}
	return dbMigration
}

func (m *DbMigration) LatestMigrations() []string {
	var migrationList []string
	lastMigrationDate := m.lastMigrationDate()

	rows, err := m.db.Query(`
		SELECT
			file_name 
		FROM
			migrations
		WHERE
			deleted_at is null
			AND created_at = $1
	`, lastMigrationDate)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if err == nil {
		var migration string
		for rows.Next() {
			rows.Scan(&migration)
			migrationList = append(migrationList, migration)
		}

	}

	return migrationList
}

func (m *DbMigration) AddToMigration(fileName string) error {
	sql := `INSERT INTO
				migrations 
			(file_name, created_at)
				values 
			($1, $2)`

	_, err := m.db.Exec(sql, fileName, m.timeString)

	return err

}

func (m *DbMigration) RemoveFromMigration(fileName string) error {
	sql := `UPDATE
				migrations 
			SET
				deleted_at = $1
			WHERE
				file_name = $2
				AND deleted_at IS NULL`

	_, err := m.db.Exec(sql, m.timeString, fileName)

	return err
}

func (m *DbMigration) MigrationExistsForFile(fileName string) bool {
	sql := `SELECT
				count(*) as cnt
			FROM
				migrations
			WHERE
				file_name = $2
			AND
				deleted_at is null`

	row := m.db.QueryRow(sql, fileName)
	var count string
	row.Scan(&count)

	cnt, err := strconv.Atoi(count)
	if err != nil {
		return false
	}

	return cnt > 0
}

func (m *DbMigration) Init() error {
	sql := `CREATE TABLE IF NOT EXISTS migrations (
				file_name VARCHAR(255),
				created_at VARCHAR(20),
				deleted_at VARCHAR(20)
			)`

	_, err := m.db.Exec(sql)

	return err
}

func (m *DbMigration) lastMigrationDate() string {
	sql := `SELECT
				max(created_at) as latest_migration
			FROM
				migrations
			WHERE
				deleted_at is null`

	row := m.db.QueryRow(sql)
	var maxdate string
	row.Scan(&maxdate)

	return maxdate
}
