# Golang Database Migrator

This package is designed to add database migration functionality to your project. It is not intended to be used independently.
If you would like to use it as a command-line tool, please use it with the command-line wrapper available here:

https://github.com/olbrichattila/godbmigrator_cmd

This package is lightweight and does not include any database drivers.
You can create migration SQL files in a designated folder.

---

## Migration File Structure
Follow this naming convention for migration files:
[date-time]-migrate-[custom-content].sql

Files will be processed in ascending order, so it's important to create a unique identifier, such as a timestamp.

--- 

### Example:
```
2024-05-27_19_49_38-migrate.sql
2024-05-27_19_49_38-rollback.sql
2024-05-27_19_50_04-migrate.sql
2024-05-27_19_50_04-rollback.sql
```
This ensures that migrations are executed in the correct sequence.

---

### Adding to Your Code
#### Import the module:

```
migrator "github.com/olbrichattila/godbmigrator"
```
#### Requirements
You need to have:
- A database connection (*sql.DB)
- Migration table prefix, or it can be empty string as well


### Table Prefix
The prefix parameter sets a database table prefix. If you set it to xyz, the migration table will be named xyz_migration to store applied migrations.

---

### Migration Operations
#### Example: Running Migrations
```
migrationFilePath := "./migration"
m := migrator.New(db, "prefix", migrationFilePath)
err := m.Migrate(count)
if err != nil {
    panic("Error: " + err.Error())
}
```
#### Example: Rolling Back Migrations
```
migrationFilePath := "./migration"
m := migrator.New(db, "prefix", migrationFilePath)
err := m.Rollback(count)
if err != nil {
    panic("Error: " + err.Error())
}
```
#### Example: Refreshing Migrations
A refresh rolls back all migrations and applies them from scratch.
```
migrationFilePath := "./migration"
m := migrator.New(db, "prefix", migrationFilePath)
err := m.Refresh()
if err != nil {
    panic("Error: " + err.Error())
}
```
#### Example: Creating a New Migration File
```
migrationFilePath := "./migration"
m := migrator.New(db, "prefix", migrationFilePath)
err := m.AddNewMigrationFiles("custom-text-or-empty")
if err != nil {
    panic("Error: " + err.Error())
}
```

---

### Checksum Validator
You can validate whether any migration file has changed since it was applied.
```
m := migrator.New(db, "prefix", migrationFilePath)
err := m.Migrate(3)
if err != nil {
    panic("Error: " + err.Error())
}

errors := m.ChecksumValidation(db, tablePrefix, testChecksumFixtureFolder)
// 'errors' contains a list of error strings ([]string). If empty, there are no validation errors.
```

---
### Baseline Operations
#### Create a Baseline of Your Existing Database Structure
```
m := migrator.New(db, "prefix", migrationFilePath)
err := m.SaveBaseline()
if err != nil {
    panic("Error: " + err.Error())
}
```
#### Restore Baseline
```
m := migrator.New(db, "prefix", migrationFilePath)
err := m.LoadBaseline()
if err != nil {
    panic("Error: " + err.Error())
}
```

---

### Migration Report
The application stores a migration audit report, where you can track applied migrations, rollbacks, and any errors encountered during migration.
#### Fetching the migration report as a readable string:
```
migrationFilePath := "./migration"
m := migrator.New(db, "prefix", migrationFilePath)
report, err := m.Report()
if err != nil {
    panic("Error: " + err.Error())
}

fmt.Println(report)
```

---

## Subscribing to Messages:

You can subscribe to messages using a callback function. When a message is received, the callback ```func(int, string)``` will be invoked with the message type ID and an optional message text.

### Example
```
m := migrator.New(db, "prefix", migrationFilePath)
m.SubscribeToMessages(func(et int, msg string) {
    fmt.Println(et, msg)
})
```

---

### Available Make Targets
You can use make commands for quick setup and testing:
```
make run
make install
make run-test
```

---
### Supported Database Drivers
- SQLite
- MySQL
- PostgreSQL
- Firebird / InterBase (except for baseline operations)

## About me:
- Learn more about me on my personal website. https://attilaolbrich.co.uk/menu/my-story
- Check out my latest blog blog at my personal page. https://attilaolbrich.co.uk/blog/1/single
