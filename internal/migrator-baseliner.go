package baseliner

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	queryTypeTables = "table"
	queryTypeIndex = "index" // SQLite specific
	queryTypeViews = "view"
	queryTypeProcedures = "procedure"
	queryTypeFunctions = "function"
	queryTypeTriggers = "trigger"
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

type retrievalInstruction struct{
	query string
	fieldPosition int
}

type baselineInstruction struct {
	execute []string
	listerQueries map[string]string
	schemaRetrievalQueries map[string]retrievalInstruction
	activeDatabaseSql string
}

type baselilner struct {
	baselineInstruction baselineInstruction
	db *sql.DB
	databaseName string
}

func (b *baselilner) Save(migrationFilePath string) error {
	baselineInstruction, err := b.getEngineSpecificInstructions()
	if err != nil {
		return err
	}
	b.baselineInstruction = *baselineInstruction

	filename := migrationFilePath + "/baseline.sql"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
		
	}
	defer file.Close()

	err = b.GetSchemaData(func(schemaDef string, useDelimiter bool) error {
		if useDelimiter {
			_, err := file.WriteString("DELIMITER ;;\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline %v", err)
			}	
		}

		_, err := file.WriteString(schemaDef + ";\n\n")
		if err != nil {
			return fmt.Errorf("cannot save baseline %v", err)
		}

		if useDelimiter {
			_, err := file.WriteString("DELIMITER ;\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline %v", err)
			}	
		}

		return nil
	});

	if err != nil {
		return err;
	}

	return nil;

}

func (b *baselilner) Load(migrationFilePath string) error {
	filename := migrationFilePath + "/baseline.sql"

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file opening error %s Error:%v", filename, err)
		
	}
	defer file.Close()

	var statementBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	isDelimiterSeparation := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "--") || line == "" {
			continue 
		}

		if line == "DELIMITER ;;" {
			isDelimiterSeparation = true
			continue
		}

		if line != "DELIMITER ;" {
			statementBuilder.WriteString(line + "\n")	
		}
		
		if b.detectStatementEnd(line, isDelimiterSeparation) {
			query := statementBuilder.String()
			statementBuilder.Reset()

			
			_, err := b.db.Exec(query)
			if err != nil {
				return fmt.Errorf("SQL Execution Error: %v query: %s", err, query)
			}
		}
	}

	return scanner.Err()
}

func (*baselilner) detectStatementEnd(line string, isDelimiterSeparation bool) bool {
	if isDelimiterSeparation {
		return line == "DELIMITER ;"
	}

	return strings.HasSuffix(line, ";")
}

func (b *baselilner) GetSchemaData(callback func(string, bool) error) error {
	databaseName, err := b.getActiveDatabaseName()
	if err != nil {
		return err
	}
	b.databaseName = databaseName

	for _, pType := range b.baselineInstruction.execute {
		tables, err := b.getInformationSchemaList(pType)
		if err != nil {
			return err
		}

		for _, tableName := range tables {
			schemaSql, err := b.getSchemaSql(pType, tableName);
			if err != nil {
				return err
			}

			err = callback(schemaSql, b.useDelimiter(pType))
			if err != nil {
				return err
			}
		}
	}

	return nil;
}

func (b *baselilner) getInformationSchemaList(queryType string) ([]string, error) {
	sql, err := b.getListQuery(queryType)
	if err != nil {
		return nil, err
	}

	sqlParams := make([]any, 0);
	if b.databaseName != "" {
		sqlParams = append(sqlParams, b.databaseName)
	}
	
	rows, err := b.db.Query(sql, sqlParams...)
	if err != nil {
		return nil, fmt.Errorf("cannot get schema definition from mySQL, (%s) error: %v", sql, err)
	}

	defer rows.Close()

	result := make([]string, 0)
	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, fmt.Errorf("cannot get schema row definition from MySql %v", err)
		}
		
		result = append(result, tableName)
	}

	return result, nil;
}

func (b *baselilner) getSchemaSql(queryType string, tableName string) (string, error) {
	sql, fieldIndex, err := b.getSchemaQueryByType(queryType, tableName)
	if err != nil {
		return "", err
	}

	rows, err := b.db.Query(sql);
	if err != nil {
		return  "", fmt.Errorf("cannot get schema data (%s), error: %v", sql, err)
	}
	
	defer rows.Close()
	
	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("could not get schema for table: %s, (%s) error: %v", tableName, sql, err)
	}

	data := make([]interface{}, len(cols))
	pointers := make([]interface{}, len(cols))
	for i := range pointers {
		data[i] = &pointers[i]
	}

	rows.Next()
	err = rows.Scan(data...)
	if err != nil {
		return "", fmt.Errorf("could not get schema for table: %s, (%s) error: %v", tableName, sql, err)
	}

	switch val := pointers[fieldIndex].(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	default:
		return "", fmt.Errorf("index %d, in SQL (%s) returns no data, possibly not sufficient database privilege", fieldIndex, sql)
	}
}

func (b *baselilner) useDelimiter(typeText string) bool {
	if typeText == queryTypeProcedures || typeText == queryTypeFunctions || typeText == queryTypeTriggers {
		return true
	}

	return false
}

func (b *baselilner) getActiveDatabaseName() (string, error) {
	if b.baselineInstruction.activeDatabaseSql == "" {
		return "", nil
	}

	var dbName string
	err := b.db.QueryRow(b.baselineInstruction.activeDatabaseSql).Scan(&dbName)
	if err != nil {
		log.Fatal(err)
	}
	return dbName, nil
}

func (b *baselilner) getListQuery(queryType string) (string, error) {
	if query, ok := b.baselineInstruction.listerQueries[queryType]; ok {
		return query, nil
	}
	
	return "", fmt.Errorf("query type %s not implemented", queryType)
	
}

func (b *baselilner) getSchemaQueryByType(queryType string, name string) (string, int, error) {
	if query, ok := b.baselineInstruction.schemaRetrievalQueries[queryType]; ok {
		return fmt.Sprintf(query.query, name), query.fieldPosition, nil
	}

	return "", 0, fmt.Errorf("query type %s not implemented", queryType)
}

func (b *baselilner) GetDb() *sql.DB {
	return b.db
}

func (b *baselilner) getEngineSpecificInstructions() (*baselineInstruction, error) {
	driverType := reflect.TypeOf(b.db.Driver()).String()

	if strings.Contains(driverType, "sqlite") {
		return b.getSQLiteInstruction(), nil
	}

	if strings.Contains(driverType, "mysql") {
		return b.getMySQLInstruction(), nil
	}

	if strings.Contains(driverType, "pq") || strings.Contains(driverType, "postgres") {
		return b.getPostgreSQLInstruction(), nil
	}
	if strings.Contains(driverType, "firebirdsql") {
		return b.getFirebirdSQLInstruction(), nil
	}

	return nil, fmt.Errorf("the driver used %s does not match any known driver by the application", driverType)
}

func (b *baselilner) getSQLiteInstruction() *baselineInstruction {
	return &baselineInstruction{
		execute: []string{queryTypeTables, queryTypeIndex, queryTypeViews, queryTypeTriggers},
		listerQueries: map[string]string{
			queryTypeTables: "SELECT name FROM sqlite_master WHERE type = \"table\"",
			queryTypeIndex: "SELECT name FROM sqlite_master WHERE type = \"index\"",
			queryTypeViews: "SELECT name FROM sqlite_master WHERE type = \"view\"",
			queryTypeTriggers: "SELECT name FROM sqlite_master WHERE type = \"trigger\"",
			
		},
		schemaRetrievalQueries: map[string]retrievalInstruction{
			queryTypeTables: {query: "SELECT sql FROM sqlite_master WHERE type = \"table\" and name = \"%s\"", fieldPosition: 0},
			queryTypeIndex: {query: "SELECT sql FROM sqlite_master WHERE type = \"index\" and name = \"%s\"", fieldPosition: 0},
			queryTypeViews: {query: "SELECT sql FROM sqlite_master WHERE type = \"view\" and name = \"%s\"", fieldPosition: 0},
			queryTypeTriggers: {query: "SELECT sql FROM sqlite_master WHERE type = \"trigger\" and name = \"%s\"", fieldPosition: 0},
		},
	}
}

func (b *baselilner) getMySQLInstruction() *baselineInstruction {
	return &baselineInstruction{
		execute: []string{queryTypeTables, queryTypeViews, queryTypeProcedures, queryTypeFunctions, queryTypeTriggers},
		listerQueries: map[string]string{
			queryTypeTables: "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = ?",
			queryTypeViews: "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ?",
			queryTypeProcedures: "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'PROCEDURE' AND ROUTINE_SCHEMA = ?",
			queryTypeFunctions: "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'FUNCTION' AND ROUTINE_SCHEMA = ?",
			queryTypeTriggers: "SELECT TRIGGER_NAME FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ?",
		},
		schemaRetrievalQueries: map[string]retrievalInstruction{
			queryTypeTables: {query: "SHOW CREATE TABLE `%s`", fieldPosition: 1},
			queryTypeViews: {query: "SHOW CREATE VIEW `%s`", fieldPosition: 1},
			queryTypeProcedures: {query: "SHOW CREATE PROCEDURE `%s`", fieldPosition: 2},
			queryTypeFunctions: {query: "SHOW CREATE FUNCTION `%s`", fieldPosition: 2},
			queryTypeTriggers: {query: "SHOW CREATE TRIGGER `%s`", fieldPosition: 2},
		},
		activeDatabaseSql: "SELECT DATABASE()",
	}
}

func (b *baselilner) getPostgreSQLInstruction() *baselineInstruction {
	return &baselineInstruction{}
}

func (b *baselilner) getFirebirdSQLInstruction() *baselineInstruction {
	return &baselineInstruction{}
}