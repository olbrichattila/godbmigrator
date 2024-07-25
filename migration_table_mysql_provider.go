package migrator

import "fmt"

type mySQLMigrationTableSQLProvider struct {
	tablePrefix string
}

func (p *mySQLMigrationTableSQLProvider) createMigrationSQL() string {
	return fmt.Sprintf(defaultMigrationCreateTableSQL, p.tablePrefix)
}

func (p *mySQLMigrationTableSQLProvider) createReportSQL() string {
	return fmt.Sprintf(defaultMigrationReportCreateTableSQL, p.tablePrefix)
}
