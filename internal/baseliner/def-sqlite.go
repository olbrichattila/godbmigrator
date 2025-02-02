package baseliner

func (b *baselilner) getSQLiteInstruction() *baselineInstruction {
	return &baselineInstruction{
		execute: []string{queryTypeTables, queryTypeIndex, queryTypeViews, queryTypeTriggers},
		listerQueries: map[string]string{
			queryTypeTables:   "SELECT name FROM sqlite_master WHERE type = \"table\"",
			queryTypeIndex:    "SELECT name FROM sqlite_master WHERE type = \"index\"",
			queryTypeViews:    "SELECT name FROM sqlite_master WHERE type = \"view\"",
			queryTypeTriggers: "SELECT name FROM sqlite_master WHERE type = \"trigger\"",
		},
		schemaRetrievalQueries: map[string]retrievalInstruction{
			queryTypeTables: {
				query: "SELECT sql FROM sqlite_master WHERE type = \"table\" and name = \"%s\"",
			},
			queryTypeIndex: {
				query: "SELECT sql FROM sqlite_master WHERE type = \"index\" and name = \"%s\"",
			},
			queryTypeViews: {
				query: "SELECT sql FROM sqlite_master WHERE type = \"view\" and name = \"%s\"",
			},
			queryTypeTriggers: {
				query: "SELECT sql FROM sqlite_master WHERE type = \"trigger\" and name = \"%s\"",
			},
		},
	}
}
