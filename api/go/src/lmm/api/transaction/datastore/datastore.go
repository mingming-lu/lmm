package datastore

import (
	"context"

	"github.com/pkg/errors"

	"cloud.google.com/go/datastore"
)

type contextKey string

const (
	txKey contextKey = "datastoreTxkey"
)

type Transaction struct {
	*datastore.Transaction
}

func (tx *Transaction) Commit() error {
	_, err := tx.Transaction.Commit()
	return err
}

func (tx *Transaction) Rollback() error {
	return tx.Transaction.Rollback()
}

type TransactionManager struct {
	datastore *datastore.Client
}

func NewTransactionManager(client *datastore.Client) *TransactionManager {
	return &TransactionManager{datastore: client}
}

func FromContext(c context.Context) (*Transaction, error) {
	tx, ok := c.Value(txKey).(*datastore.Transaction)
	if !ok {
		return nil, errors.New("transaction not found")
	}
	return &Transaction{Transaction: tx}, nil
}

func (tm *TransactionManager) New(c context.Context) (context.Context, error) {
	return tm.new(c)
}

func (tm *TransactionManager) NewReadOnly(c context.Context) (context.Context, error) {
	return tm.new(c, datastore.ReadOnly)
}

func (tm *TransactionManager) new(c context.Context, opts ...datastore.TransactionOption) (context.Context, error) {
	tx, err := tm.datastore.NewTransaction(c, opts...)
	if err != nil {
		return nil, err
	}
	return context.WithValue(c, txKey, tx), nil
}

func (tm *TransactionManager) RunInTransaction(c context.Context, f func(c context.Context) error) error {
	return tm.runInTransaction(c, f)
}

func (tm *TransactionManager) RunInReadOnly(c context.Context, f func(c context.Context) error) error {
	return tm.runInTransaction(c, f, datastore.ReadOnly)
}

func (tm *TransactionManager) runInTransaction(c context.Context, f func(c context.Context) error, opts ...datastore.TransactionOption) error {
	txCtx, err := tm.new(c, opts...)
	if err != nil {
		return err
	}

	tx, err := FromContext(txCtx)
	if err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}

	if err := f(txCtx); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return tx.Commit()
}
