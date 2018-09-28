package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

// DB is a database abstraction
type DB interface {
	Begin(c context.Context, opts *sql.TxOptions) (Tx, error)

	Conn(c context.Context) (*sql.Conn, error)

	Close() error

	Driver() driver.Driver

	Exec(c context.Context, query string, args ...interface{}) (sql.Result, error)

	Prepare(c context.Context, query string) Stmt

	Ping(c context.Context) error

	Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error)

	QueryRow(c context.Context, query string, args ...interface{}) *sql.Row

	SetConnMaxLifetime(d time.Duration)

	SetMaxIdleConns(n int)

	SetMaxOpenConns(n int)

	Stats() sql.DBStats
}
