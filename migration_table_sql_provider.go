package migrator

import "fmt"

const (
	defaultTablePrefix                   = "olb"
	defaultMigrationReportCreateTableSQL = `CREATE TABLE IF NOT EXISTS %s_migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT)`

	defaultMigrationCreateTableSQL = `CREATE TABLE IF NOT EXISTS %s_migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME)`
)

type migrationTableSQLProvider interface {
	createMigrationSQL() string
	createReportSQL() string
}

func migrationTableProviderByDriverName(driverName, tablePrefix string) (migrationTableSQLProvider, error) {
	switch driverName {
	case dbTypeSqlite:
		return &sqliteMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbTypePostgres:
		return &postgresMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbTypeMySQL:
		return &mySQLMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbTypeFirebird:
		return &firebirdMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	default:
		return nil, fmt.Errorf("provider %s does not exists", driverName)
	}
}
