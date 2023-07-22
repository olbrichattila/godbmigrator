package main

import (
	"fmt"
	"os"
	"strconv"

	migrator "github.com/olbrichattila/go-database-migrator"
)

func main() {
	db, err := migrator.NewSqliteStore("./data/database.sqlite")
	if err != nil {
		panic("Error: " + err.Error())
	}

	commandLineParameter, count, err := commandLineParameters()
	if err != nil {
		printUsage()
		return
	}

	// migrationProvider, err := migrator.NewMigrationProvider("json", nil)
	migrationProvider, err := migrator.NewMigrationProvider("db", db)
	if err != nil {
		panic("Error: " + err.Error())
	}

	switch commandLineParameter {
	case "migrate":
		err = migrator.Migrate(db, migrationProvider, count)
		if err != nil {
			panic("Error: " + err.Error())
		}
	case "rollback":
		err = migrator.Rollback(db, migrationProvider, count)
		if err != nil {
			panic("Error: " + err.Error())
		}
	default:
		printUsage()
	}
}

func commandLineParameters() (string, int, error) {
	if len(os.Args) == 2 {
		return os.Args[1], 0, nil
	}

	if len(os.Args) == 3 {
		migrationCount, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return "", 0, err
		}

		return os.Args[1], migrationCount, nil
	}

	return "", 0, fmt.Errorf("Invalid command line paramteres")
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
