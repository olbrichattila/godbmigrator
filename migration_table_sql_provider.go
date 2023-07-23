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

func MigrationTableProviderByDriverName(driverName string) (MigrationTableSqlProvider, error) {
	switch driverName {
	case dbTypeSqlite:
		return &SqliteMigrationTableSqlProvider{}, nil
	case dbTypePostgres:
		return &PostgresMigrationTableSqlProvider{}, nil
	case dbTypeMySql:
		return &MySqlMigrationTableSqlProvider{}, nil
	default:
		return nil, fmt.Errorf("Provider %s does not exists", driverName)
	}
}

func (p *PostgresMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at VARCHAR(20),
		deleted_at VARCHAR(20)
	)`
}

func (p *SqliteMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at VARCHAR(20),
		deleted_at VARCHAR(20)
	)`
}

func (p *MySqlMigrationTableSqlProvider) CreateSql() string {
	return `CREATE TABLE IF NOT EXISTS migrations (
		file_name VARCHAR(255),
		created_at VARCHAR(20),
		deleted_at VARCHAR(20)
	)`
}
