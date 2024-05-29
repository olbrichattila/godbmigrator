package migrator

type postgresMigrationTableSQLProvider struct {
}

func (p *postgresMigrationTableSQLProvider) createMigrationSQL() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at TIMESTAMP,
		deleted_at TIMESTAMP
	)`
}

func (p *postgresMigrationTableSQLProvider) createReportSQL() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at TIMESTAMP,
		message TEXT
	)`
}
