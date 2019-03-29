package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Rows = sql.Rows

// Config defines config of database
type Config struct {
	// auhthorization
	User     string
	Password string

	// base dsn
	Protocol string
	Address  string
	Database string
	Options  url.Values // dsn query parameters

	// Retry defines retry time if connection fails, 0 for no retry, < 0 for infinite retries
	Retry int
}

// DSN converts c to dsn string
// A DSN in its fullest form: username:password@protocol(address)/dbname?param=value
// See examples for MySQL: https://github.com/go-sql-driver/mysql
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?%s",
		c.User, c.Password, c.Protocol, c.Address, c.Database, c.Options.Encode())
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

func RollbackWithError(tx Tx, err error) error {
	if err == nil {
		return tx.Rollback()
	}

	if e := tx.Rollback(); e != nil {
		return errors.Wrap(err, e.Error())
	}

	return err
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

var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// called with stringBuilderPool.Put(sb)
func mustStringBuilder() *strings.Builder {
	sb, ok := stringBuilderPool.Get().(*strings.Builder)
	if !ok {
		panic("expected a *strings.Builder")
	}
	sb.Reset()
	return sb
}

// Masks builds a string in a (?,?,?) like format
func Masks(length uint) string {
	switch length {
	case 0:
		return ""
	case 1:
		return "(?)"
	default:
		sb := mustStringBuilder()
		defer stringBuilderPool.Put(sb)

		sb.WriteString("(?")
		for i := uint(1); i < length; i++ {
			sb.WriteString(",?")
		}
		sb.WriteString(")")

		return sb.String()
	}
}

// FieldMasks builds a (field,?,?,?) like string
func FieldMasks(field string, length uint) string {
	switch length {
	case 0:
		return ""
	default:
		sb := mustStringBuilder()
		defer stringBuilderPool.Put(sb)

		sb.WriteString("(")
		sb.WriteString(field)
		for i := uint(0); i < length; i++ {
			sb.WriteString(",?")
		}
		sb.WriteString(")")

		return sb.String()
	}
}
