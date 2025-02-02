package migrate

import "fmt"

type postgresMigrationTableSQLProvider struct {
	tablePrefix string
}

func (p *postgresMigrationTableSQLProvider) createMigrationSQL() string {
	sql := `CREATE TABLE IF NOT EXISTS %s_migrations (
		file_name VARCHAR(255),
		created_at TIMESTAMP,
		deleted_at TIMESTAMP,
		checksum CHAR(32)
	)`

	return fmt.Sprintf(sql, p.tablePrefix)
}

func (p *postgresMigrationTableSQLProvider) createReportSQL() string {
	sql := `CREATE TABLE IF NOT EXISTS %s_migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at TIMESTAMP,
		message TEXT
	)`

	return fmt.Sprintf(sql, p.tablePrefix)
}
