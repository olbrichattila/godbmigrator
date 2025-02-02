package migrator

import (
	"fmt"

	"github.com/olbrichattila/godbmigrator/internal/dbtypemanager"
)

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
		deleted_at DATETIME,
		checksum CHAR(32)
		)`
)

type migrationTableSQLProvider interface {
	createMigrationSQL() string
	createReportSQL() string
}

func migrationTableProviderByDriverName(driverName, tablePrefix string) (migrationTableSQLProvider, error) {
	switch driverName {
	case dbtypemanager.DbTypeSqlite:
		return &sqliteMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbtypemanager.DbTypePostgres:
		return &postgresMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbtypemanager.DbTypeMySQL:
		return &mySQLMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	case dbtypemanager.DbTypeFirebird:
		return &firebirdMigrationTableSQLProvider{tablePrefix: tablePrefix}, nil
	default:
		return nil, fmt.Errorf("provider %s does not exists", driverName)
	}
}
