package migrator

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DbMigration struct {
	db                  *sql.DB
	timeString          string
	sqlBindingParameter string
}

func newDbMigration(db *sql.DB) (*DbMigration, error) {
	timeString := time.Now().Format("2006-01-02 15:04:05")
	dbMigration := &DbMigration{db: db, timeString: timeString}
	driverType, err := dbMigration.diverType()
	if err != nil {
		return nil, err
	}
	dbMigration.setSqlBindingParameter(driverType)
	createSqlProvider, err := MigrationTableProviderByDriverName(driverType)
	if err != nil {
		return nil, err
	}

	err = dbMigration.Init(createSqlProvider)
	if err != nil {
		return nil, err
	}
	return dbMigration, nil
}

func (m *DbMigration) LatestMigrations() ([]string, error) {
	var migrationList []string
	lastMigrationDate := m.lastMigrationDate()

	rows, err := m.db.Query(fmt.Sprintf(
		`SELECT file_name 
		 FROM migrations
		 WHERE created_at = %s 
		 AND deleted_at IS NULL
		 ORDER BY file_name DESC`,
		m.getBindingParameter(1),
	), lastMigrationDate)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err == nil {
		var migration string
		for rows.Next() {
			rows.Scan(&migration)
			migrationList = append(migrationList, migration)
		}

	}

	return migrationList, nil
}

func (m *DbMigration) AddToMigration(fileName string) error {
	sql := fmt.Sprintf(`INSERT INTO migrations  
			(file_name, created_at)
			VALUES (%s, %s)`,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
	)

	_, err := m.db.Exec(sql, fileName, m.timeString)

	return err

}

func (m *DbMigration) RemoveFromMigration(fileName string) error {
	sql := fmt.Sprintf(`UPDATE migrations 
			SET deleted_at = %s
			WHERE file_name = %s
			AND deleted_at IS NULL`,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
	)

	_, err := m.db.Exec(sql, m.timeString, fileName)

	return err
}

func (m *DbMigration) MigrationExistsForFile(fileName string) bool {
	sql := fmt.Sprintf(`SELECT count(*) as cnt
			FROM migrations
			WHERE file_name = %s
			AND deleted_at IS NULL`,
		m.getBindingParameter(1),
	)

	row := m.db.QueryRow(sql, fileName)

	var count string
	row.Scan(&count)

	cnt, err := strconv.Atoi(count)

	if err != nil {
		return false
	}

	return cnt > 0
}

func (m *DbMigration) Init(createSqlProvider MigrationTableSqlProvider) error {
	sql := createSqlProvider.CreateSql()

	_, err := m.db.Exec(sql)

	return err
}

func (m *DbMigration) lastMigrationDate() string {
	sql := `SELECT max(created_at) as latest_migration
			FROM migrations
			WHERE deleted_at IS NULL`

	row := m.db.QueryRow(sql)
	var maxdate string
	row.Scan(&maxdate)

	return maxdate
}

func (m *DbMigration) setSqlBindingParameter(driverType string) {
	if driverType == "pq" {
		m.sqlBindingParameter = "$1"

		return
	}

	m.sqlBindingParameter = "?"
}

func (m *DbMigration) getBindingParameter(index int) string {
	if m.sqlBindingParameter == "?" {
		return "?"
	}

	return fmt.Sprintf("$%d", index)
}

func (m *DbMigration) diverType() (string, error) {
	driverType := reflect.TypeOf(m.db.Driver()).String()

	if strings.Contains(driverType, "mysql") {
		return "mysql", nil
	}

	if strings.Contains(driverType, "pg") || strings.Contains(driverType, "postgres") {
		return "pg", nil
	}

	if strings.Contains(driverType, "sqlite") {
		return "sqlite", nil
	}

	return "", fmt.Errorf("The driver used %s does not match any known dirver by the application", driverType)
}
