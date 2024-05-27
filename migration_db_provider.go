package migrator

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	dbTypeSqlite   = "sqlite"
	dbTypePostgres = "pg"
	dbTypeMySql    = "mysql"
	dbTypeFirebird = "firebird"
)

type DbMigration struct {
	db                  *sql.DB
	timeString          string
	sqlBindingParameter string
}

func newDbMigration(db *sql.DB) (*DbMigration, error) {
	dbMigration := &DbMigration{db: db}
	dbMigration.ResetDate()
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

func (m *DbMigration) ResetDate() {
	m.timeString = time.Now().Format("2006-01-02 15:04:05")
}

func (m *DbMigration) Migrations(isLatest bool) ([]string, error) {
	var migrationList []string
	var rows *sql.Rows
	var err error

	lastMigrationDate := m.lastMigrationDate()
	if lastMigrationDate == "" {
		return migrationList, nil
	}

	if isLatest {
		rows, err = m.latestMigrations(lastMigrationDate)
	} else {
		rows, err = m.allMigrations()
	}

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

func (m *DbMigration) latestMigrations(lastMigrationDate string) (*sql.Rows, error) {
	return m.db.Query(fmt.Sprintf(
		`SELECT file_name 
		 FROM migrations
		 WHERE created_at = %s 
		 AND deleted_at IS NULL
		 ORDER BY file_name DESC`,
		m.getBindingParameter(1),
	), lastMigrationDate)
}

func (m *DbMigration) allMigrations() (*sql.Rows, error) {
	return m.db.Query(
		`SELECT file_name 
		 FROM migrations
		 WHERE deleted_at IS NULL
		 ORDER BY file_name DESC`,
	)
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
	sql := createSqlProvider.CreateMigrationSql()

	_, err := m.db.Exec(sql)
	if err != nil {
		return err
	}

	sql = createSqlProvider.CreateReportSql()
	_, err = m.db.Exec(sql)

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
	if driverType == dbTypePostgres {
		m.sqlBindingParameter = "$"

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
		return dbTypeMySql, nil
	}

	if strings.Contains(driverType, "pq") || strings.Contains(driverType, "postgres") {
		return dbTypePostgres, nil
	}

	if strings.Contains(driverType, "sqlite") {
		return dbTypeSqlite, nil
	}

	if strings.Contains(driverType, "firebirdsql") {
		return dbTypeFirebird, nil
	}

	return "", fmt.Errorf("the driver used %s does not match any known dirver by the application", driverType)
}

func (m *DbMigration) GetJsonFileName() string {
	// dummy, not used in db version, need due to interface
	return ""
}

func (m *DbMigration) SetJsonFileName(filePath string) {
	// dummy, not used in db version, need due to interface
}

func (m *DbMigration) AddToMigrationReport(fileName string, errorToLog error) error {
	sql := fmt.Sprintf(`INSERT INTO migration_reports
			(file_name, created_at, result_status, message)
			VALUES (%s, %s, %s, %s)`,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
		m.getBindingParameter(3),
		m.getBindingParameter(4),
	)

	message := "ok"
	status := "success"
	if errorToLog != nil {
		message = errorToLog.Error()
		status = "error"
	}

	createdAt := time.Now().Format("2006-01-02 15:04:05")

	_, err := m.db.Exec(sql, fileName, createdAt, status, message)

	return err
}
