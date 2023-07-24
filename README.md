# Golang Database migrator

!!Under development...

This is a package, not to be used individually, but serves the purpose to add migration to your project.
If you would like to use this as a command line tool, please use this package with the command line wrapper
available here:

git@github.com:olbrichattila/godbmigrator_cmd.git

## Create migration SQL files into the the folder migrations (will be configurable)

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

Currently the command line utility supports only SqLite, the build in solution shoud work, but not tested with oher databases

## Example migrate: (where the db is your *sql.DB)

```
migrationProvider, err := migrator.NewMigrationProvider("json")
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Migrate(db, migrationProvider, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Example rollback: (where the db is your *sql.DB)
```
migrationProvider, err := migrator.NewMigrationProvider("json", nil)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Rollback(db, migrationProvider, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Migrate With database provider
```
db, err := migrator.NewSqliteStore("./data/database.sqlite")
if err != nil {
    panic("Error: " + err.Error())
}

migrationProvider, err := migrator.NewMigrationProvider("db", db)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Migrate(db, migrationProvider, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Rollback With database provider
```
db, err := migrator.NewSqliteStore("./data/database.sqlite")
if err != nil {
    panic("Error: " + err.Error())
}

migrationProvider, err := migrator.NewMigrationProvider("db", db)
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Rollback(db, migrationProvider, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Example, create new migration file:
```
err := migrator.AddNewMigrationFiles(secondParameter)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Available make targets:

```
make install
```
