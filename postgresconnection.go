package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresStore(user, dbname, password string) (*sql.DB, error) {
	// @TODO pramaeterize sslmode, add host and port
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s sslmode=disable",
		user,
		dbname,
		password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
