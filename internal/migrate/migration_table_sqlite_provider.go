package migrate

import "fmt"

type sqliteMigrationTableSQLProvider struct {
	tablePrefix string
}

func (p *sqliteMigrationTableSQLProvider) createMigrationSQL() string {
	return fmt.Sprintf(defaultMigrationCreateTableSQL, p.tablePrefix)
}

func (p *sqliteMigrationTableSQLProvider) createReportSQL() string {
	return fmt.Sprintf(defaultMigrationReportCreateTableSQL, p.tablePrefix)
}
