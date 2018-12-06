package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

type Rows = sql.Rows

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

// SQLOptions has sql options
type SQLOptions struct {
	Where   string
	OrderBy string
	Limit   string
}

// Count counts the number of columns from table
func Count(c context.Context, db DB, table, column string, opts SQLOptions) (uint, error) {
	var count uint
	sql := "SELECT COUNT(" + column + ") FROM " + table
	if opts.Where != "" {
		sql += " " + opts.Where
	}

	stmt := db.Prepare(c, sql)
	defer stmt.Close()

	row := stmt.QueryRow(c)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Masks build a string in a (?,?,?) like format
func Masks(length uint) string {
	switch length {
	case 0:
		return ""
	case 1:
		return "(?)"
	default:
		var (
			i uint = 1
			s      = "(?"
		)
		for ; i < length; i++ {
			s += ",?"
		}
		return s + ")"
	}
}
