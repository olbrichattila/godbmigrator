// Package messager receives a message and calls back a callback
package messager

// Message types
const (
	MigratedItems = iota
	NothingToRollback
	RolledBack
	RunningMigrations
	SkipRollback
	RunningRollback
	MigrationFileCreated
)

// CallbackFunc Function type
type CallbackFunc func(int, string)

// New messenger instance
func New() Messager {
	return &message{
		callbacks: make([]CallbackFunc, 0),
	}
}

// Messager abstracts message register and dispatch
type Messager interface {
	Dispatch(eventType int, message string)
	Register(CallbackFunc)
}

type message struct {
	callbacks []CallbackFunc
}

// Dispatch will receive and dispatch the message to the callback
func (m *message) Dispatch(eventType int, message string) {
	for _, cb := range m.callbacks {
		cb(eventType, message)
	}
}

// Register callback
func (m *message) Register(callback CallbackFunc) {
	m.callbacks = append(m.callbacks, callback)
}
