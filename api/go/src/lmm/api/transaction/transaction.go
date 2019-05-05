package transaction

import "context"

type Transaction interface {
	Commit() error
	Rollback() error
}

type Manager interface {
	New(c context.Context) (context.Context, error)
	NewReadOnly(c context.Context) (context.Context, error)

	RunInTransaction(c context.Context, f func(c context.Context) error) error
	RunInReadOnly(c context.Context, f func(c context.Context) error) error
}

type NopTx struct{}

func (tx *NopTx) Commit() error {
	return nil
}

func (tx *NopTx) Rollback() error {
	return nil
}

type NopTxManager struct{}

func (tx *NopTxManager) New(c context.Context) (context.Context, error) {
	return c, nil
}

func (tx *NopTxManager) NewReadOnly(c context.Context) (context.Context, error) {
	return c, nil
}

func (tx *NopTxManager) RunInTransaction(c context.Context, f func(c context.Context) error) error {
	return f(c)
}

func (tx *NopTxManager) RunInReadOnly(c context.Context, f func(c context.Context) error) error {
	return f(c)
}
