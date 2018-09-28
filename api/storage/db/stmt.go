package db

import (
	"context"
	"database/sql"
)

// Stmt interface likes a database/sql.Stmt
type Stmt interface {
	Close() error

	Exec(c context.Context, args ...interface{}) (sql.Result, error)

	Query(c context.Context, args ...interface{}) (*sql.Rows, error)

	QueryRow(c context.Context, args ...interface{}) *sql.Row
}

type stmt struct {
	*sql.Stmt
}

func (s *stmt) Close() error {
	return s.Stmt.Close()
}

func (s *stmt) Exec(c context.Context, args ...interface{}) (sql.Result, error) {
	return s.Stmt.ExecContext(c, args...)
}

func (s *stmt) Query(c context.Context, args ...interface{}) (*sql.Rows, error) {
	return s.Stmt.QueryContext(c, args...)
}

func (s *stmt) QueryRow(c context.Context, args ...interface{}) *sql.Row {
	return s.Stmt.QueryRowContext(c, args...)
}
