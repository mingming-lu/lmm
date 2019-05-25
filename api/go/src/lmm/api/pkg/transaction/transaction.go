package transaction

import "context"

// Transaction abstraction
type Transaction interface {
	context.Context
	Rollback() error
	Commit() error
}

type IsolationLevel int

type Option struct {
	IsolationLevel IsolationLevel
	ReadOnly       bool
}

type FuncRunInTransaction = func(Transaction) error

// Manager is a transaction manager abstraction
type Manager interface {
	Begin(context.Context, *Option) (Transaction, error)
	RunInTransaction(context.Context, FuncRunInTransaction, *Option) error
}
