package migrator

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Baseliner interface {
	Save(migrationFilePath string) error
	Load(migrationFilePath string) error
}

type SpecificBaseliner interface {
	GetSchemaData(callback func(string, bool) error) error
	GetDb() *sql.DB
}

func GetBaseliner(db *sql.DB) (Baseliner, error) {
	driverType := reflect.TypeOf(db.Driver()).String()

	if strings.Contains(driverType, "mysql") {
		return newBaseliner(NewMySQLBaseliner(db)), nil
	}

	if strings.Contains(driverType, "pq") || strings.Contains(driverType, "postgres") {
		return newBaseliner(NewPostgresBaseliner(db)), nil
	}

	if strings.Contains(driverType, "sqlite") {
		return newBaseliner(NewSQLiteBaseliner(db)), nil
	}

	if strings.Contains(driverType, "firebirdsql") {
		return newBaseliner(NewFirebirdBaseliner(db)), nil
	}

	return nil, fmt.Errorf("the driver used %s does not match any known driver by the application", driverType)
}

