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
	dbTypeMySQL    = "mysql"
	dbTypeFirebird = "firebird"
)

type dbMigration struct {
	db                  *sql.DB
	timeString          string
	sqlBindingParameter string
}

type reportRow struct {
	FileName     string
	CreatedAt    string
	ResultStatus string
	Message      string
}

func newDbMigration(db *sql.DB) (*dbMigration, error) {
	dbMigration := &dbMigration{db: db}
	dbMigration.resetDate()
	driverType, err := dbMigration.diverType()
	if err != nil {
		return nil, err
	}
	dbMigration.setSQLBindingParameter(driverType)
	createSQLProvider, err := migrationTableProviderByDriverName(driverType)
	if err != nil {
		return nil, err
	}

	err = dbMigration.init(createSQLProvider)
	if err != nil {
		return nil, err
	}
	return dbMigration, nil
}

func (m *dbMigration) resetDate() {
	m.timeString = time.Now().Format("2006-01-02 15:04:05")
}

func (m *dbMigration) migrations(isLatest bool) ([]string, error) {
	var migrationList []string
	var rows *sql.Rows
	var err error

	lastMigrationDate, err := m.lastMigrationDate()
	if err != nil {
		return nil, err
	}

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
			err := rows.Scan(&migration)
			if err != nil {
				return nil, err
			}
			migrationList = append(migrationList, migration)
		}

	}

	return migrationList, nil
}

func (m *dbMigration) latestMigrations(lastMigrationDate string) (*sql.Rows, error) {
	return m.db.Query(fmt.Sprintf(
		`SELECT file_name 
		 FROM migrations
		 WHERE created_at = %s 
		 AND deleted_at IS NULL
		 ORDER BY file_name DESC`,
		m.getBindingParameter(1),
	), lastMigrationDate)
}

func (m *dbMigration) allMigrations() (*sql.Rows, error) {
	return m.db.Query(
		`SELECT file_name 
		 FROM migrations
		 WHERE deleted_at IS NULL
		 ORDER BY file_name DESC`,
	)
}

func (m *dbMigration) addToMigration(fileName string) error {
	sql := fmt.Sprintf(`INSERT INTO migrations  
			(file_name, created_at)
			VALUES (%s, %s)`,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
	)

	_, err := m.db.Exec(sql, fileName, m.timeString)

	return err

}

func (m *dbMigration) removeFromMigration(fileName string) error {
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

func (m *dbMigration) migrationExistsForFile(fileName string) (bool, error) {
	sql := fmt.Sprintf(`SELECT count(*) as cnt
			FROM migrations
			WHERE file_name = %s
			AND deleted_at IS NULL`,
		m.getBindingParameter(1),
	)

	row := m.db.QueryRow(sql, fileName)

	var count string
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	cnt, err := strconv.Atoi(count)

	if err != nil {
		return false, nil
	}

	return cnt > 0, nil
}

func (m *dbMigration) init(createSQLProvider migrationTableSQLProvider) error {
	sql := createSQLProvider.createMigrationSQL()

	_, err := m.db.Exec(sql)
	if err != nil {
		return err
	}

	sql = createSQLProvider.createReportSQL()
	_, err = m.db.Exec(sql)

	return err
}

func (m *dbMigration) lastMigrationDate() (string, error) {
	sql := `SELECT max(created_at) as latest_migration
			FROM migrations
			WHERE deleted_at IS NULL`

	row := m.db.QueryRow(sql)
	var maxdate string
	err := row.Scan(&maxdate)

	return maxdate, err
}

func (m *dbMigration) setSQLBindingParameter(driverType string) {
	if driverType == dbTypePostgres {
		m.sqlBindingParameter = "$"

		return
	}

	m.sqlBindingParameter = "?"
}

func (m *dbMigration) getBindingParameter(index int) string {
	if m.sqlBindingParameter == "?" {
		return "?"
	}

	return fmt.Sprintf("$%d", index)
}

func (m *dbMigration) diverType() (string, error) {
	driverType := reflect.TypeOf(m.db.Driver()).String()

	if strings.Contains(driverType, "mysql") {
		return dbTypeMySQL, nil
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

func (m *dbMigration) getJSONFileName() string {
	// dummy, not used in db version, need due to interface
	return ""
}

func (m *dbMigration) SetJSONFilePath(_ string) {
	// dummy, not used in db version, need due to interface
}

func (m *dbMigration) AddToMigrationReport(fileName string, errorToLog error) error {
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

func (m *dbMigration) Report() (string, error) {
	rows, err := m.db.Query(`SELECT  file_name, created_at, result_status, message FROM migration_reports`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var row reportRow
	var builder strings.Builder
	for rows.Next() {
		err := rows.Scan(&row.FileName, &row.CreatedAt, &row.ResultStatus, &row.Message)
		if err != nil {
			return "", err
		}
		str := fmt.Sprintf(
			"Created at: %s, File Name: %s, Status: %s, Message: %s\n",
			row.CreatedAt,
			row.FileName,
			row.ResultStatus,
			row.Message,
		)
		builder.WriteString(str)
	}

	return builder.String(), nil
}
