module dodbtester/migrator_test

go 1.18

replace github.com/olbrichattila/godbmigrator => ../

require (
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/olbrichattila/godbmigrator v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
