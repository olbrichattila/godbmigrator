package baseliner

func (b *baselilner) getPostgreSQLInstruction() *baselineInstruction {
	return &baselineInstruction{
		execute: []string{queryTypeTables, queryTypeIndex, queryTypeViews, queryTypeMaterialViews, queryTypeFunctions, queryTypeProcedures},
		listerQueries: map[string]string{
			queryTypeTables:        "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = $1",
			queryTypeIndex:         "SELECT x.indexrelid FROM pg_index x JOIN pg_class c ON c.oid = x.indrelid JOIN pg_class i ON i.oid = x.indexrelid JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = $1 ORDER BY c.relname, i.relname",
			queryTypeViews:         "SELECT viewname FROM pg_catalog.pg_views WHERE schemaname = $1",
			queryTypeMaterialViews: "SELECT matviewname FROM pg_catalog.pg_matviews WHERE schemaname = $1",
			queryTypeFunctions:     "SELECT routine_name FROM information_schema.routines WHERE routine_type = 'FUNCTION' AND specific_schema = $1",
			queryTypeProcedures:    "SELECT routine_name FROM information_schema.routines WHERE routine_type = 'PROCEDURE' AND specific_schema = $1",
		},
		schemaRetrievalQueries: map[string]retrievalInstruction{
			queryTypeTables: {
				query:                "SELECT 'CREATE TABLE ' || c.relname || E'(\\n' || array_to_string(array_agg(E'\\t' || a.attname || ' ' || pg_catalog.format_type(a.atttypid, a.atttypmod)), E',\\n') || E'\\n)' AS create_table_sql FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace JOIN pg_attribute a ON a.attrelid = c.oid WHERE c.relkind = 'r' AND n.nspname = '%s' AND a.attnum > 0 and c.relname = '%s' GROUP BY c.relname",
				dbNameShouldBePassed: true,
			},
			queryTypeIndex: {
				query:                "SELECT pg_catalog.pg_get_indexdef(i.oid) AS create_index_sql FROM pg_index x JOIN pg_class i ON i.oid = x.indexrelid JOIN pg_namespace n ON n.oid = i.relnamespace WHERE n.nspname = '%s' and x.indexrelid = '%s'",
				dbNameShouldBePassed: true,
			},
			queryTypeViews: {
				query:                "SELECT ('CREATE VIEW ' || viewname || ' AS ' || definition) as sql FROM pg_catalog.pg_views WHERE schemaname = '%s' and viewname = '%s';",
				dbNameShouldBePassed: true,
			},
			queryTypeMaterialViews: {
				query:                "SELECT ('CREATE MATERIALIZED VIEW ' || matviewname || ' AS ' || definition) as sql FROM pg_catalog.pg_matviews WHERE schemaname = '%s' AND matviewname = '%s'",
				dbNameShouldBePassed: true,
			},
			queryTypeFunctions: {
				query:                "SELECT pg_get_functiondef(p.oid) as sql FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = '%s' and p.proname = '%s'",
				dbNameShouldBePassed: true,
			},
			queryTypeProcedures: {
				query:                "SELECT pg_get_functiondef(p.oid) as sql FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = '%s' AND p.prokind = 'p' and p.proname = '%s'",
				dbNameShouldBePassed: true,
			},
		},
		activeDatabaseSQL: "SELECT current_schema()",
	}
}
