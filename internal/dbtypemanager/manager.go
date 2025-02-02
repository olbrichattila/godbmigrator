// Package dbtypemanager contains constants and shared functions to determine specifics for a db type
package dbtypemanager

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Specific Database constants
const (
	DbTypeSqlite   = "sqlite"
	DbTypePostgres = "pg"
	DbTypeMySQL    = "mysql"
	DbTypeFirebird = "firebird"
)

// GetDiverType returns driver type according to the passed DB
func GetDiverType(db *sql.DB) (string, error) {
	driverType := reflect.TypeOf(db.Driver()).String()

	if strings.Contains(driverType, "mysql") {
		return DbTypeMySQL, nil
	}

	if strings.Contains(driverType, "pq") || strings.Contains(driverType, "postgres") {
		return DbTypePostgres, nil
	}

	if strings.Contains(driverType, "sqlite") {
		return DbTypeSqlite, nil
	}

	if strings.Contains(driverType, "firebirdsql") {
		return DbTypeFirebird, nil
	}

	return "", fmt.Errorf("DB manager: the driver used %s does not match any known driver by the application", driverType)
}
