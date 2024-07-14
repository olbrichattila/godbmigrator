package migrator

type sqliteMigrationTableSQLProvider struct {
}

func (p *sqliteMigrationTableSQLProvider) createMigrationSQL() string {
	return defaultMigrationCreateTableSQL
}

func (p *sqliteMigrationTableSQLProvider) createReportSQL() string {
	return defaultMigrationReportCreateTableSQL
}
