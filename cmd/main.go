// main package is for testing only locally, rename it to use
package main

// rename cmd.go to bak, then rename this to .go and test your changes
// don't forget to go mod tidy before pushing back to git
import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/nakagami/firebirdsql"
	migrator "github.com/olbrichattila/godbmigrator"
)

func main() {
	migrationFilePath := "./migrations"
	dbType, provider, function, count, add := params()

	fmt.Printf("Running with %s, %s, %s, %d %s\n", dbType, provider, function, count, add)

	if add != "" {
		migrator.AddNewMigrationFiles(add, "")
	}

	db, err := getConnection(dbType)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	MigrationProvider, err := migrator.NewMigrationProvider(provider, "", db)
	if err != nil {
		panic(err.Error())
	}

	switch function {
	case "migrate":
		err = migrator.Migrate(db, MigrationProvider, migrationFilePath, count)
	case "rollback":
		err = migrator.Rollback(db, MigrationProvider, migrationFilePath, count)
	case "refresh":
		err = migrator.Refresh(db, MigrationProvider, migrationFilePath)
	case "report":
		report, err := migrator.Report(db, MigrationProvider, migrationFilePath)
		if err != nil {
			panic("Error: " + err.Error())
		}
		fmt.Print(report)
	case "baseline-save":
		err := migrator.SaveBaseline(db, migrationFilePath)
		if err != nil {
			panic("Error: " + err.Error())
		}
	case "baseline-load":
		err := migrator.LoadBaseline(db, migrationFilePath)
		if err != nil {
			panic("Error: " + err.Error())
		}
	default:
		panic("Unknown function " + function)
	}

	if err != nil {
		panic("Error: " + err.Error())
	}
}

func params() (string, string, string, int, string) {
	db := flag.String("db", "sqlite", "Database driver name")
	function := flag.String("function", "migrate", "function=migrate/rollback/report/baseline-save/baseline-load")
	count := flag.Int("count", 0, "count=1")
	provider := flag.String("provider", "db", "provider=db/json")
	add := flag.String("add", "", "--add=filename")
	flag.Parse()

	return *db, *provider, *function, *count, *add
}

func connectToFirebaseDatabase() (*sql.DB, error) {
	dataSourceName := "SYSDBA:masterkey@localhost:3050/firebird/data/employee.fdb"

	return sql.Open("firebirdsql", dataSourceName)
}

func connectToPostresql() (*sql.DB, error) {
	connectionString := "user=root password=root dbname=root sslmode=disable"

	return sql.Open("postgres", connectionString)
}

func connectToSqLite() (*sql.DB, error) {
	return sql.Open("sqlite3", "./data/data.db")
}

func getConnection(dbType string) (*sql.DB, error) {
	var conn *sql.DB
	var err error

	switch dbType {
	case "sqlite":
		conn, err = connectToSqLite()
	case "firebase":
		conn, err = connectToFirebaseDatabase()
	case "postgresql":
		conn, err = connectToPostresql()
	default:
		return nil, fmt.Errorf("invalid dbtype " + dbType)
	}

	return conn, err
}
