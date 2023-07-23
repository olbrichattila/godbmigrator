package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	PgSslModeDisable    = "disable"
	PgSslModeRequire    = "require"
	PgSslModeVerifyCa   = "verify-ca"
	PgSslModeVerifyFull = "verify-full"
	PgSslModePrefer     = "prefer"
	PgSslModeAllow      = "allow"
)

func NewPostgresStore(
	host string,
	port int,
	user,
	password,
	dbname,
	sslmode string,
) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
