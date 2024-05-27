package migrator

type MySqlMigrationTableSqlProvider struct {
}

func (p *MySqlMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *MySqlMigrationTableSqlProvider) CreateReportSql() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT
	)`
}
