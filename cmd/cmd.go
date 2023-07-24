package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
	migrator "github.com/olbrichattila/godbmigrator"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	db, err := getDatabase()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	commandLineParameter, secondParameter, err := commandLineParameters()
	if err != nil {
		printUsage()
		return
	}

	migrationProvider, err := getMigrationProvider(db)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	switch commandLineParameter {
	case "migrate":
		count, err := getCount(secondParameter)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		err = migrator.Migrate(db, migrationProvider, count)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	case "rollback":
		count, err := getCount(secondParameter)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		err = migrator.Rollback(db, migrationProvider, count)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	case "add":
		err := migrator.AddNewMigrationFiles(secondParameter)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	default:
		printUsage()
	}
}

func getCount(param string) (int, error) {
	if param == "" {
		return 0, nil
	}

	count, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func commandLineParameters() (string, string, error) {
	if len(os.Args) == 2 {
		return os.Args[1], "", nil
	}

	if len(os.Args) == 3 {
		return os.Args[1], os.Args[2], nil
	}

	return "", "0", fmt.Errorf("Invalid command line paramteres")
}

func printUsage() {
	fmt.Printf(`Usage:
	migrator migrate
	migrator rollback
	migrator migrate 2
	migrator rollback 2

The number of rollbacks and migrates are not mandatory.
If it is set, for rollbacks it only apply for the last rollback batch
`)
}

func getDatabase() (*sql.DB, error) {
	dbConnection := os.Getenv("DB_CONNECTION")
	switch dbConnection {
	case "sqlite":
		db, err := migrator.NewSqliteStore(os.Getenv("DB_DATABASE"))
		return db, err
	case "pgsql":
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			return nil, err
		}
		db, err := migrator.NewPostgresStore(
			os.Getenv("DB_HOST"),
			port,
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
			migrator.PgSslModeDisable,
		)
		return db, err
	case "mysql":
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			return nil, err
		}
		db, err := migrator.NewMysqlStore(
			os.Getenv("DB_HOST"),
			port,
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
		)
		return db, err
	default:
		return nil, fmt.Errorf("Invalid DB_CONNECTION %s", dbConnection)
	}
}

func getMigrationProvider(db *sql.DB) (migrator.MigrationProvider, error) {
	migrationProvider := os.Getenv("MIGRATOR_MIGRATION_PROVIDER")

	switch migrationProvider {
	case "db", "":
		return migrator.NewMigrationProvider("db", db)
	case "json":
		return migrator.NewMigrationProvider("json", nil)
	default:
		return nil, fmt.Errorf("Migration provider for type %s does not exists", migrationProvider)
	}
}
