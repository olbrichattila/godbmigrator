package migrator

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	queryTypeTables = "table"
	queryTypeViews = "view"
	queryTypeProcedures = "procedure"
	queryTypeFunctions = "function"
	queryTypeTriggers = "trigger"
)

func NewMySQLBaseliner(db *sql.DB) SpecificBaseliner {
	return &blMySQL{db: db}
}

type blMySQL struct {
	db *sql.DB
	databaseName string
}

func (b *blMySQL) GetSchemaData(callback func(string, bool) error) error {
	processList := []string{queryTypeTables, queryTypeViews, queryTypeProcedures, queryTypeFunctions, queryTypeTriggers}
	
	databaseName, err := b.getActiveDatabaseName()
	if err != nil {
		return err
	}
	b.databaseName = databaseName

	for _, pType := range processList {
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

func (b *blMySQL) getInformationSchemaList(queryType string) ([]string, error) {
	sql, err := b.getListQuery(queryType)
	if err != nil {
		return nil, err
	}
	
	rows, err := b.db.Query(sql, b.databaseName)
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

func (b *blMySQL) getSchemaSql(queryType string, tableName string) (string, error) {
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

	if val, ok := pointers[fieldIndex].([]byte); ok {
		return string(val) , nil
	}

	return "", fmt.Errorf("index %d, in SQL %s returns no data, possibly not sufficient database privilege", fieldIndex, sql)
}

func (b *blMySQL) useDelimiter(typeText string) bool {
	if typeText == queryTypeProcedures || typeText == queryTypeFunctions || typeText == queryTypeTriggers {
		return true
	}

	return false
}

func (b *blMySQL) getActiveDatabaseName() (string, error) {
	var dbName string
	err := b.db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		log.Fatal(err)
	}
	return dbName, nil
}

func (b *blMySQL) getListQuery(queryType string) (string, error) {
	switch (queryType) {
		case queryTypeTables:
			return "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA = ?", nil
		case queryTypeViews:
			return "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.VIEWS WHERE TABLE_SCHEMA = ?", nil
		case queryTypeProcedures: 
			return "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'PROCEDURE' AND ROUTINE_SCHEMA = ?", nil
		case queryTypeFunctions:
			return "SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_TYPE = 'FUNCTION' AND ROUTINE_SCHEMA = ?", nil
		case queryTypeTriggers:
			return "SELECT TRIGGER_NAME FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_SCHEMA = ?", nil
		default:
			return "", fmt.Errorf("query type %s not implemented", queryType)
	}
}

func (b *blMySQL) getSchemaQueryByType(queryType string, name string) (string, int, error) {
	switch (queryType) {
		case queryTypeTables:
			return fmt.Sprintf("SHOW CREATE TABLE `%s`", name), 1, nil
		case queryTypeViews:
			return fmt.Sprintf("SHOW CREATE VIEW `%s`", name), 1, nil
		case queryTypeProcedures:
			return fmt.Sprintf("SHOW CREATE PROCEDURE `%s`", name), 2, nil
		case queryTypeFunctions:
			 return fmt.Sprintf("SHOW CREATE FUNCTION `%s`", name), 2, nil
		case queryTypeTriggers:
			 return fmt.Sprintf("SHOW CREATE TRIGGER `%s`", name), 2, nil
		default:
			return "", 0, fmt.Errorf("query type %s not implemented", queryType)
	}
}

func (b *blMySQL) GetDb() *sql.DB {
	return b.db
}