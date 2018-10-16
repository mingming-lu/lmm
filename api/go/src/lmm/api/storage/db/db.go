package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

// Config defines config of database
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// Retry defines retry time if connection fails, 0 for no retry, < 0 for infinite retries
	Retry int
}

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
