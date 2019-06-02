package datastore

import (
	"context"
	"lmm/api/pkg/transaction"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

type txImpl struct {
	context.Context
	*datastore.Transaction
}

var (
	testTxImpl    transaction.Transaction = &txImpl{}
	testTxManager transaction.Manager     = &TransactionManager{}
)

func mustTxImpl(i interface{}) *txImpl {
	tx, ok := i.(*txImpl)
	if !ok {
		panic("not a *txImpl")
	}

	return tx
}

func MustContext(i interface{}) context.Context {
	return mustTxImpl(i).Context
}

func MustTransaction(i interface{}) *datastore.Transaction {
	return mustTxImpl(i).Transaction
}

func (tx *txImpl) Commit() error {
	_, err := tx.Transaction.Commit()
	return err
}

type TransactionManager struct {
	source *datastore.Client
}

func NewTransactionManager(source *datastore.Client) *TransactionManager {
	return &TransactionManager{
		source: source,
	}
}

func (tm *TransactionManager) Begin(c context.Context, opts *transaction.Option) (transaction.Transaction, error) {
	if opts != nil && opts.ReadOnly {
		tx, err := tm.source.NewTransaction(c, datastore.ReadOnly)
		if err != nil {
			return nil, err
		}
		return &txImpl{c, tx}, nil
	}

	tx, err := tm.source.NewTransaction(c)
	if err != nil {
		return nil, err
	}
	return &txImpl{c, tx}, nil
}

func (tm *TransactionManager) RunInTransaction(c context.Context, f func(tx transaction.Transaction) error, opts *transaction.Option) error {
	tx, err := tm.Begin(c, opts)
	if err != nil {
		return err
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			panic(recovered)
		}
	}()

	if err := f(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return err
	}

	return tx.Commit()
}
