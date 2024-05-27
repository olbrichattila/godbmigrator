package migrator

import "fmt"

type MigrationTableSqlProvider interface {
	CreateMigrationSql() string
	CreateReportSql() string
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
