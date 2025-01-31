package migrator

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func NewSQLiteBaseliner(db *sql.DB) Baseliner {
	return &blSQLite{db: db}
}

type blSQLite struct {
	db *sql.DB
}

func (b *blSQLite) Save(migrationFilePath string) error {
	filename := migrationFilePath + "/baseline.sql"
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
		
	}
	defer file.Close()

	err = b.getSchemaData(func(schemaDef string) error {
		_, err := file.WriteString(schemaDef + ";\n")
		if err != nil {
			return fmt.Errorf("cannot save baseline %v", err)
		}

		return nil
	});

	if err != nil {
		return err;
	}

	return nil;

}

func (b *blSQLite) Load(migrationFilePath string) error {
	filename := migrationFilePath + "/baseline.sql"

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file opening error %s Error:%v", filename, err)
		
	}
	defer file.Close()

	var statementBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "--") || line == "" {
			continue 
		}

		statementBuilder.WriteString(line + "\n")
		if strings.HasSuffix(line, ";") {
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

func (b *blSQLite) getSchemaData(callback func(string) error) error {
	sql := "SELECT sql FROM sqlite_master WHERE type IN ('table', 'index', 'view', 'trigger')"
	rows, err := b.db.Query(sql)
	if err != nil {
		return fmt.Errorf("cannot get schema definition from SQLite %v", err)
	}

	defer rows.Close()
	var sqlText string
	for rows.Next() {
		err := rows.Scan(&sqlText)
		if err != nil {
			return fmt.Errorf("cannot get schema row definition from SQLite %v", err)
		}
		err = callback(sqlText)
		if err != nil {
			return err
		}	
	}

	return nil;
}
