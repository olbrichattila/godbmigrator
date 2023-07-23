package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlStore(user, dbname, password string) (*sql.DB, error) {
	// @todo add host
	connStr := fmt.Sprintf(
		"%s:%s@tcp(localhost:3306)/%s",
		user,
		password,
		dbname)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
