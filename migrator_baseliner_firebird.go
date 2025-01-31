package migrator

import (
	"database/sql"
	"fmt"
)

func NewFirebirdBaseliner(db *sql.DB) Baseliner {
	return &blFirebird{db: db}
}

type blFirebird struct {
	db *sql.DB
}

func (*blFirebird) Save(migrationFilePath string) error {
	return fmt.Errorf("not implemented")
}

func (b *blFirebird) Load(migrationFilePath string) error {
	return fmt.Errorf("not yet implemented")
}