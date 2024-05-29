package migrator

import "fmt"

type migrationTableSQLProvider interface {
	createMigrationSQL() string
	createReportSQL() string
}

func migrationTableProviderByDriverName(driverName string) (migrationTableSQLProvider, error) {
	switch driverName {
	case dbTypeSqlite:
		return &sqliteMigrationTableSQLProvider{}, nil
	case dbTypePostgres:
		return &postgresMigrationTableSQLProvider{}, nil
	case dbTypeMySQL:
		return &mySQLMigrationTableSQLProvider{}, nil
	case dbTypeFirebird:
		return &firebirdMigrationTableSQLProvider{}, nil
	default:
		return nil, fmt.Errorf("provider %s does not exists", driverName)
	}
}
