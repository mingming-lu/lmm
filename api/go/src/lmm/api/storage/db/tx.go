package db

import (
	"context"
	"database/sql"
)

// Tx interface likes a database/sql.Tx
type Tx interface {
	Commit() error

	Exec(c context.Context, query string, args ...interface{}) (sql.Result, error)

	Prepare(c context.Context, query string) Stmt

	Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error)

	QueryRow(c context.Context, query string, args ...interface{}) *sql.Row

	Rollback() error
}

type tx struct {
	*sql.Tx
}

func (tx *tx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *tx) Exec(c context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.Tx.ExecContext(c, query, args...)
}

func (tx *tx) Prepare(c context.Context, query string) Stmt {
	st, err := tx.Tx.PrepareContext(c, query)
	if err != nil {
		panic(err)
	}
	return &stmt{Stmt: st}
}

func (tx *tx) Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.Tx.QueryContext(c, query, args...)
}

func (tx *tx) QueryRow(c context.Context, query string, args ...interface{}) *sql.Row {
	return tx.Tx.QueryRowContext(c, query, args...)
}

func (tx *tx) Rollback() error {
	return tx.Tx.Rollback()
}
