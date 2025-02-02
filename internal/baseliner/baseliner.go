// Package baseliner saves and restores current database structure
package baseliner

import (
	"database/sql"
)

const (
	queryTypeTables        = "table"
	queryTypeIndex         = "index" // SQLite specific
	queryTypeViews         = "view"
	queryTypeMaterialViews = "materialView"
	queryTypeProcedures    = "procedure"
	queryTypeFunctions     = "function"
	queryTypeTriggers      = "trigger"

	// SQL file Delimiters
	openingDelimiter = "DELIMITER ;"
	closingDelimiter = "DELIMITER ;;"
)

// New baseliner, which saves and restores database structure
func New(db *sql.DB) Baseliner {
	return &baselilner{
		db: db,
	}
}

// Baseliner implements Save and Load
type Baseliner interface {
	Save(migrationFilePath string) error
	Load(migrationFilePath string) error
}

type retrievalInstruction struct {
	query                string
	fieldPosition        int
	dbNameShouldBePassed bool
}

type baselineInstruction struct {
	execute                []string
	listerQueries          map[string]string
	schemaRetrievalQueries map[string]retrievalInstruction
	activeDatabaseSQL      string
}

type baselilner struct {
	baselineInstruction baselineInstruction
	db                  *sql.DB
	databaseName        string
}
