# Golang Database migrator

!!Under development...

## Command line usage:

# Without building the app

Migrate:
```go run cmd/cmd.go migrate```

Rollback:
```go run cmd/cmd.go rollback```

### Migrate or rollback specified amount of migrations (like 2)

Migrate:
```go run cmd/cmd.go migrate 2```

Rollback:
```go run cmd/cmd.go rollback 2```

### When building the application.

```make install```
The build folder will contain the migrator executable.

Usage is the same but using the application:

```
migrator migrate
migrator rollback

migrator migrate 2
migrator rollback 2
```

The number of rollbacks and migrates are not mandatory.
If it is set, for rollbacks it only apply for the last rollback batch

## Adding to your code.

Import the module:

```migrator "github.com/olbrichattila/go-database-migrator"```

You need to have a DB connection, and a migration provider.

The migration provider stores the migration status to:
- json
- database (under implementation)
- (others to come)

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
migrationProvider, err := migrator.NewMigrationProvider("json")
if err != nil {
    panic("Error: " + err.Error())
}

err = migrator.Rollback(db, migrationProvider, count)
if err != nil {
    panic("Error: " + err.Error())
}
```

## Available make targets:

```
mage migrate
make rollback
make install
```
## Coming soon

.env where you can define the database connection. migration file paths and migration provider type therefore it can be used as a full featured command line migrator.

