package migrator

type SqliteMigrationTableSqlProvider struct {
}

func (p *SqliteMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *SqliteMigrationTableSqlProvider) CreateReportSql() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT
	)`
}
