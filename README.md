# Golang Database migrator

This is a package, not to be used individually, but serves the purpose to add migration to your project.
If you would like to use this as a command line tool, please use this package with the command line wrapper
available here:

The package is a lightweight installation, does not contain any driver.

https://github.com/olbrichattila/godbmigrator_cmd

## Create migration SQL files into a folderfmysql

## What is the provider?

Currently it supports two type of migration provider, json and database.
This is the way the migrator knows which migration was executed and when.

If the json provider is used, then a json file will be saved next to the migration files:
```./migrations/migrations.json```

If the db provider is user, then a migrations table will be created in the same database where you are migrating to.

## Migration file structure

Follow the structure:
[id]-migrate-[custom-content].sql

The files will be processed in ascending order, therefore it is important to create an id as follows:
For example:
```
0001-migrate.sql
0001-rollback.sql
0002-migrate.sql
0002-rollback.sql
0003-migrate.sql
0003-rollback.sql
0004-migrate.sql
0004-rollback.sql
0005-migrate-new.sql
0005-rollback-new.sql
0006-migrate-new.sql
0006-rollback-new.sql
```

## Adding to your code.

Import the module:

```migrator "github.com/olbrichattila/godbmigrator"```

You need to have a DB connection, and a migration provider.

The migration provider stores the migration status to:
- json
- database
- (others to come)

## prefix is the database table prefix, if you set xyz, then it will create xyz_migration table for storing migrations
Currently the command line utility supports only SqLite, the build in solution should work, but not tested with other databases

## Example migrate: (where the db is your *sql.DB)

```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("json", "prefix")
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Migrate(db, MigrationProvider, migrationFilePath, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Example rollback: (where the db is your *sql.DB)
```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("json", "prefix", nil)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Rollback(db, MigrationProvider, migrationFilePath, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Example refresh: (where the db is your *sql.DB)
Refresh is when everything rolled back and migrated from scratch
```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("json", "prefix", nil)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Refresh(db, MigrationProvider, migrationFilePath)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Migrate with database provider
```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("db", "prefix", db)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Migrate(db, MigrationProvider, migrationFilePath, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Rollback with database provider
```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("db", "prefix", db)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Rollback(db, MigrationProvider, migrationFilePath, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Refresh with database provider
Refresh is when everything rolled back and migrated from scratch
```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("db", "prefix", db)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Refresh(db, MigrationProvider, migrationFilePath)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Example, create new migration file:
```
migrationFilePath := "./migration"
err := migrator.AddNewMigrationFiles(migrationFilePath, "custom-text-or-empty")
if err != nil {
    panic("Error: " + err.Error())
}
```


## Migration report
The application stores a migration audit report, where you can track rollback, migrations, errors thrown during migration.
Fetching the migration report to a readable string:

```
migrationFilePath := "./migration"
MigrationProvider, err := migrator.NewMigrationProvider("db", "prefix", db)
if err != nil {
    panic("Error: " + err.Error())
}

report, err := migrator.Report(db, MigrationProvider, migrationFilePath)
if err != nil {
    panic("Error: " + err.Error())
}

fmt.Println(report)
```

## Available make targets:
```
make run
make install
make run-test
```

## Currently supported database drivers:

- SqLite
- MySql
- PostgresQl
- Firebird / Interbase

## To be expected in next version

- Store what was migrated as well even when migration error happened and rolled back
- Report function, which can create a short report (migrations and rollback in order) or a large report with the queries migrator executed (or tried to execute)
