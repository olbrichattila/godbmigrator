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
	tablePrefix         string
	timeString          string
	sqlBindingParameter string
}

type reportRow struct {
	FileName     string
	CreatedAt    string
	ResultStatus string
	Message      string
}

func newDbMigration(db *sql.DB, tablePrefix string) (*dbMigration, error) {
	if tablePrefix == "" {
		tablePrefix = defaultTablePrefix
	}

	dbMigration := &dbMigration{
		db:          db,
		tablePrefix: tablePrefix,
	}

	dbMigration.ResetDate()
	driverType, err := dbMigration.diverType()
	if err != nil {
		return nil, err
	}
	dbMigration.setSQLBindingParameter(driverType)
	createSQLProvider, err := migrationTableProviderByDriverName(driverType, tablePrefix)
	if err != nil {
		return nil, err
	}

	err = dbMigration.init(createSQLProvider)
	if err != nil {
		return nil, err
	}
	return dbMigration, nil
}

func (m *dbMigration) ResetDate() {
	m.timeString = time.Now().Format(timeFormat)
}

func (m *dbMigration) Migrations(isLatest bool) ([]string, error) {
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
		 FROM %s_migrations
		 WHERE created_at = %s 
		 AND deleted_at IS NULL
		 ORDER BY file_name DESC`,
		m.tablePrefix,
		m.getBindingParameter(1),
	), lastMigrationDate)
}

func (m *dbMigration) allMigrations() (*sql.Rows, error) {
	return m.db.Query(
		fmt.Sprintf(
			`SELECT file_name 
			FROM %s_migrations
			WHERE deleted_at IS NULL
			ORDER BY file_name DESC`,
			m.tablePrefix,
		),
	)
}

func (m *dbMigration) AddToMigration(fileName string) error {
	sql := fmt.Sprintf(`INSERT INTO %s_migrations  
			(file_name, created_at)
			VALUES (%s, %s)`,
		m.tablePrefix,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
	)

	_, err := m.db.Exec(sql, fileName, m.timeString)

	return err

}

func (m *dbMigration) RemoveFromMigration(fileName string) error {
	sql := fmt.Sprintf(`UPDATE %s_migrations 
			SET deleted_at = %s
			WHERE file_name = %s
			AND deleted_at IS NULL`,
		m.tablePrefix,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
	)

	_, err := m.db.Exec(sql, m.timeString, fileName)

	return err
}

func (m *dbMigration) MigrationExistsForFile(fileName string) (bool, error) {
	sql := fmt.Sprintf(`SELECT count(*) as cnt
			FROM %s_migrations
			WHERE file_name = %s
			AND deleted_at IS NULL`,
		m.tablePrefix,
		m.getBindingParameter(1),
	)

	row := m.db.QueryRow(sql, fileName)

	var count string
	err := row.Scan(&count)
	if err != nil {
		return false, nil
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
	sql := fmt.Sprintf(
		`SELECT max(created_at) as latest_migration
			FROM %s_migrations
			WHERE deleted_at IS NULL`,
		m.tablePrefix,
	)

	row := m.db.QueryRow(sql)
	var maxdate string
	err := row.Scan(&maxdate)
	if err != nil {
		return "", nil
	}

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

func (m *dbMigration) GetJSONFileName() string {
	// dummy, not used in db version, need due to interface
	return ""
}

func (m *dbMigration) SetJSONFilePath(_ string) {
	// dummy, not used in db version, need due to interface
}

func (m *dbMigration) AddToMigrationReport(fileName string, errorToLog error) error {
	sql := fmt.Sprintf(`INSERT INTO %s_migration_reports
			(file_name, created_at, result_status, message)
			VALUES (%s, %s, %s, %s)`,
		m.tablePrefix,
		m.getBindingParameter(1),
		m.getBindingParameter(2),
		m.getBindingParameter(3),
		m.getBindingParameter(4),
	)

	message := "ok"
	status := statusSuccess
	if errorToLog != nil {
		message = errorToLog.Error()
		status = statusError
	}

	createdAt := time.Now().Format(timeFormat)

	_, err := m.db.Exec(sql, fileName, createdAt, status, message)

	return err
}

func (m *dbMigration) Report() (string, error) {
	rows, err := m.db.Query(
		fmt.Sprintf(
			`SELECT file_name, created_at, result_status, message FROM %s_migration_reports`,
			m.tablePrefix,
		),
	)
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
			reportMessageText,
			row.CreatedAt,
			row.FileName,
			row.ResultStatus,
			row.Message,
		)
		builder.WriteString(str)
	}

	return builder.String(), nil
}
