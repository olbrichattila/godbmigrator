package migrator

import (
	"database/sql"
	"fmt"
)

func NewPostgresBaseliner(db *sql.DB) Baseliner {
	return &blPostgres{db: db}
}

type blPostgres struct {
	db *sql.DB
}

func (*blPostgres) Save(migrationFilePath string) error {
	return fmt.Errorf("not implemented")
}

func (b *blPostgres) Load(migrationFilePath string) error {
	return fmt.Errorf("not yet implemented")
}