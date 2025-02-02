// Package baseliner saves and restores current database structure
package baseliner

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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
		schemaAppend := ";\n"
		if useDelimiter {
			schemaAppend = "\n"
			_, err := file.WriteString(openingDelimiter + "\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline when string opening delimiter, error: %v", err)
			}
		}

		_, err := file.WriteString(schemaDef + schemaAppend)
		if err != nil {
			return fmt.Errorf("cannot save baseline schema sql, error: %v", err)
		}

		if useDelimiter {
			_, err := file.WriteString(closingDelimiter + "\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline when string closing delimiter, error: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

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

		if line == openingDelimiter {
			isDelimiterSeparation = true
			continue
		}

		if line != closingDelimiter {
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
		return line == closingDelimiter
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

		fmt.Println(tables)

		for _, tableName := range tables {
			schemaSQL, err := b.getSchemaSQL(pType, tableName)
			if err != nil {
				return err
			}

			err = callback(schemaSQL, b.useDelimiter(pType))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *baselilner) getInformationSchemaList(queryType string) ([]string, error) {
	sql, err := b.getListQuery(queryType)
	if err != nil {
		return nil, err
	}

	sqlParams := make([]any, 0)
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

	return result, nil
}

func (b *baselilner) getSchemaSQL(queryType string, tableName string) (string, error) {
	sql, fieldIndex, err := b.getSchemaQueryByType(queryType, tableName)
	if err != nil {
		return "", err
	}

	rows, err := b.db.Query(sql)
	if err != nil {
		return "", fmt.Errorf("cannot get schema data (%s), error: %v", sql, err)
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("could not get columns for schema for table: %s, (%s) error: %v", tableName, sql, err)
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
	if typeText == queryTypeTables || typeText == queryTypeIndex {
		return false
	}

	return true
}

func (b *baselilner) getActiveDatabaseName() (string, error) {
	if b.baselineInstruction.activeDatabaseSQL == "" {
		return "", nil
	}

	var dbName string
	err := b.db.QueryRow(b.baselineInstruction.activeDatabaseSQL).Scan(&dbName)
	if err != nil {
		log.Fatal(err)
	}
	return dbName, nil
}

func (b *baselilner) getListQuery(queryType string) (string, error) {
	if query, ok := b.baselineInstruction.listerQueries[queryType]; ok {
		return query, nil
	}

	return "", fmt.Errorf("getListQuery: query type %s not implemented", queryType)

}

func (b *baselilner) getSchemaQueryByType(queryType string, name string) (string, int, error) {
	if query, ok := b.baselineInstruction.schemaRetrievalQueries[queryType]; ok {
		if query.dbNameShouldBePassed {
			return fmt.Sprintf(query.query, b.databaseName, name), query.fieldPosition, nil
		}
		return fmt.Sprintf(query.query, name), query.fieldPosition, nil
	}

	return "", 0, fmt.Errorf("getSchemaQueryByType: query type %s not implemented", queryType)
}

func (b *baselilner) GetDb() *sql.DB {
	return b.db
}
