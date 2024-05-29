package migrator

type mySQLMigrationTableSQLProvider struct {
}

func (p *mySQLMigrationTableSQLProvider) createMigrationSQL() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *mySQLMigrationTableSQLProvider) createReportSQL() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT
	)`
}
