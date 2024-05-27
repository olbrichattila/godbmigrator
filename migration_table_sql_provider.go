package migrator

import "fmt"

type MigrationTableSqlProvider interface {
	CreateMigrationSql() string
	CreateReportSql() string
}

type SqliteMigrationTableSqlProvider struct {
}

type PostgresMigrationTableSqlProvider struct {
}

type MySqlMigrationTableSqlProvider struct {
}

type FirebirdMigrationTableSqlProvider struct {
}

func MigrationTableProviderByDriverName(driverName string) (MigrationTableSqlProvider, error) {
	switch driverName {
	case dbTypeSqlite:
		return &SqliteMigrationTableSqlProvider{}, nil
	case dbTypePostgres:
		return &PostgresMigrationTableSqlProvider{}, nil
	case dbTypeMySql:
		return &MySqlMigrationTableSqlProvider{}, nil
	case dbTypeFirebird:
		return &FirebirdMigrationTableSqlProvider{}, nil
	default:
		return nil, fmt.Errorf("provider %s does not exists", driverName)
	}
}

// Migration table SQL
func (p *PostgresMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at TIMESTAMP,
		deleted_at TIMESTAMP
	)`
}

func (p *SqliteMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *MySqlMigrationTableSqlProvider) CreateMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *FirebirdMigrationTableSqlProvider) CreateMigrationSql() string {
	return `EXECUTE BLOCK AS BEGIN
		if (not exists(select 1 from rdb$relations where rdb$relation_name = 'MIGRATIONS')) then
		execute statement 'CREATE TABLE MIGRATIONS (
			file_name VARCHAR(35),
			created_at VARCHAR(35),
			deleted_at TIMESTAMP);';
		END`
}

// Reporting table SQL
func (p *PostgresMigrationTableSqlProvider) CreateReportSql() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at TIMESTAMP,
		message TEXT
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

func (p *MySqlMigrationTableSqlProvider) CreateReportSql() string {
	return `CREATE TABLE IF NOT EXISTS migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT
	)`
}

func (p *FirebirdMigrationTableSqlProvider) CreateReportSql() string {
	return `EXECUTE BLOCK AS BEGIN
		if (not exists(select 1 from rdb$relations where rdb$relation_name = 'MIGRATION_REPORTS')) then
		execute statement 'CREATE TABLE REPORTS (
			file_name VARCHAR(35),
			result_status VARCHAR(12),
			created_at VARCHAR(35),
			message BLOB SUB_TYPE TEXT);';
		END`
}
