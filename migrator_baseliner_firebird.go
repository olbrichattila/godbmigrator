package migrator

import (
	"database/sql"
	"fmt"
)

func NewFirebirdBaseliner(db *sql.DB) SpecificBaseliner {
	return &blFirebird{db: db}
}

type blFirebird struct {
	db *sql.DB
}

func (*blFirebird) GetSchemaData(callback func(string, bool) error) error {
	return fmt.Errorf("not implemented")
}

func (b *blFirebird) GetDb() *sql.DB {
	return b.db
}
