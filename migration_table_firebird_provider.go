package migrator

type firebirdMigrationTableSQLProvider struct {
}

func (p *firebirdMigrationTableSQLProvider) createMigrationSQL() string {
	return `EXECUTE BLOCK AS BEGIN
		if (not exists(select 1 from rdb$relations where rdb$relation_name = 'MIGRATIONS')) then
		execute statement 'CREATE TABLE MIGRATIONS (
			file_name VARCHAR(35),
			created_at VARCHAR(35),
			deleted_at TIMESTAMP);';
		END`
}

func (p *firebirdMigrationTableSQLProvider) createReportSQL() string {
	return `EXECUTE BLOCK AS BEGIN
		if (not exists(select 1 from rdb$relations where rdb$relation_name = 'MIGRATION_REPORTS')) then
		execute statement 'CREATE TABLE MIGRATION_REPORTS (
			file_name VARCHAR(35),
			result_status VARCHAR(12),
			created_at VARCHAR(35),
			message BLOB SUB_TYPE TEXT);';
		END`
}
