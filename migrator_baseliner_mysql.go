package migrator

import (
	"database/sql"
	"fmt"
)

func NewMySQLBaseliner(db *sql.DB) Baseliner {
	return &blMySQL{db: db}
}

type blMySQL struct {
	db *sql.DB
}

func (*blMySQL) Save(migrationFilePath string) error {
	return fmt.Errorf("not implemented")
}

func (b *blMySQL) Load(migrationFilePath string) error {
	return fmt.Errorf("not yet implemented")
}