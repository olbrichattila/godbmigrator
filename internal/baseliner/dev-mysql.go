package baseliner

func (b *baselilner) getMySQLInstruction() *baselineInstruction {
	return &baselineInstruction{
		execute: []string{queryTypeTables, queryTypeViews, queryTypeProcedures, queryTypeFunctions, queryTypeTriggers},
		listerQueries: map[string]string{
			queryTypeTables:     "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = ?",
			queryTypeViews:      "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ?",
			queryTypeProcedures: "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'PROCEDURE' AND ROUTINE_SCHEMA = ?",
			queryTypeFunctions:  "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'FUNCTION' AND ROUTINE_SCHEMA = ?",
			queryTypeTriggers:   "SELECT TRIGGER_NAME FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ?",
		},
		schemaRetrievalQueries: map[string]retrievalInstruction{
			queryTypeTables: {
				query:         "SHOW CREATE TABLE `%s`",
				fieldPosition: 1,
			},
			queryTypeViews: {
				query:         "SHOW CREATE VIEW `%s`",
				fieldPosition: 1,
			},
			queryTypeProcedures: {
				query:         "SHOW CREATE PROCEDURE `%s`",
				fieldPosition: 2,
			},
			queryTypeFunctions: {
				query:         "SHOW CREATE FUNCTION `%s`",
				fieldPosition: 2,
			},
			queryTypeTriggers: {
				query:         "SHOW CREATE TRIGGER `%s`",
				fieldPosition: 2,
			},
		},
		activeDatabaseSQL: "SELECT DATABASE()",
	}
}
