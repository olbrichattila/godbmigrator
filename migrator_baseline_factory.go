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

func GetBaseliner(db *sql.DB) (Baseliner, error) {
	driverType := reflect.TypeOf(db.Driver()).String()

	if strings.Contains(driverType, "mysql") {
		return NewSQLiteBaseliner(db), nil;
	}

	if strings.Contains(driverType, "pq") || strings.Contains(driverType, "postgres") {
		return NewPostgresBaseliner(db), nil
	}

	if strings.Contains(driverType, "sqlite") {
		return NewSQLiteBaseliner(db), nil
	}

	if strings.Contains(driverType, "firebirdsql") {
		return NewFirebirdBaseliner(db), nil
	}

	return nil, fmt.Errorf("the driver used %s does not match any known driver by the application", driverType)
}

