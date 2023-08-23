package migrator

import "fmt"

type MigrationTableSqlProvider interface {
	CreateSql() string
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
		return nil, fmt.Errorf("Provider %s does not exists", driverName)
	}
}

func (p *PostgresMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at TIMESTAMP,
		deleted_at TIMESTAMP
	)`
}

func (p *SqliteMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *MySqlMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	)`
}

func (p *FirebirdMigrationTableSqlProvider) CreateSql() string {
	return `EXECUTE BLOCK AS BEGIN
		if (not exists(select 1 from rdb$relations where rdb$relation_name = 'MIGRATIONS')) then
		execute statement 'CREATE TABLE MIGRATIONS (
			file_name VARCHAR(25),
			created_at VARCHAR(25),
			deleted_at TIMESTAMP);';
		END`
}
