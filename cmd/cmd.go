package main

import (
	"fmt"
	"os"
	"strconv"

	migrator "github.com/olbrichattila/go-database-migrator"
)

func main() {
	// db, err := migrator.NewSqliteStore("./data/database.sqlite")
	// db, err := migrator.NewPostgresStore("localhost", 5432, "postgres", "postgres", "postgres", migrator.PgSslModeDisable)
	db, err := migrator.NewMysqlStore("localhost", 3306, "root", "H8E7kU8Y", "migrator")
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	commandLineParameter, secondParameter, err := commandLineParameters()
	if err != nil {
		printUsage()
		return
	}

	// migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	migrationProvider, err := migrator.NewMigrationProvider("db", db)
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
