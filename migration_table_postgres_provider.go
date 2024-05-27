package migrator

type PostgresMigrationTableSqlProvider struct {
}

func (p *PostgresMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at TIMESTAMP,
		deleted_at TIMESTAMP
	)`
}

func (p *PostgresMigrationTableSqlProvider) CreateReportSql() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at TIMESTAMP,
		message TEXT
	)`
}
