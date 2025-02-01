package migrator

import (
	"database/sql"
	"fmt"
)

func NewPostgresBaseliner(db *sql.DB) SpecificBaseliner {
	return &blPostgres{db: db}
}

type blPostgres struct {
	db *sql.DB
}

func (*blPostgres) GetSchemaData(callback func(string, bool) error) error {
	return fmt.Errorf("not implemented")
}

func (b *blPostgres) GetDb() *sql.DB {
	return b.db
}