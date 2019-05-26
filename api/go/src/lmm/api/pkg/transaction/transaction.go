package transaction

import (
	"context"
	"time"
)

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

func Nop() Transaction {
	return nop
}

var nop Transaction = &nopTx{}

type nopTx struct{}

func (tx *nopTx) Deadline() (time.Time, bool) {
	return context.Background().Deadline()
}

func (tx *nopTx) Done() <-chan struct{} {
	return context.Background().Done()
}

func (tx *nopTx) Err() error {
	return context.Background().Err()
}

func (tx *nopTx) Value(key interface{}) interface{} {
	return context.Background().Value(key)
}

func (tx *nopTx) Commit() error {
	return nil
}

func (tx *nopTx) Rollback() error {
	return nil
}
