package migrator

import (
	"database/sql"
	"fmt"
)

func NewSQLiteBaseliner(db *sql.DB) SpecificBaseliner {
	return &blSQLite{db: db}
}

type blSQLite struct {
	db *sql.DB
}

func (b *blSQLite) GetSchemaData(callback func(string, bool) error) error {
	sql := "SELECT sql, type FROM sqlite_master WHERE type IN ('table', 'index', 'view', 'trigger')"
	rows, err := b.db.Query(sql)
	if err != nil {
		return fmt.Errorf("cannot get schema definition from SQLite %v", err)
	}

	defer rows.Close()
	var sqlText string
	var typeText string
	for rows.Next() {
		err := rows.Scan(&sqlText, &typeText)
		if err != nil {
			return fmt.Errorf("cannot get schema row definition from SQLite %v", err)
		}
		
		err = callback(sqlText, b.useDelimiter(typeText))
		if err != nil {
			return err
		}	
	}

	return nil;
}

func (b *blSQLite) useDelimiter(typeText string) bool {
	return typeText == "trigger"
}

func (b *blSQLite) GetDb() *sql.DB {
	return b.db
}
