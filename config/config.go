// Package config contains publicly available configuration values
package config

// Migration response message types
const (
	MigratedItems = iota
	NothingToRollback
	RolledBack
	RunningMigrations
	SkipRollback
	RunningRollback
	MigrationFileCreated
)
