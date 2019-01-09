package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"lmm/api/util"
)

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDonw   = sql.ErrTxDone
)

type sqlError struct {
	query string
	args  []interface{}
}

func (e *sqlError) Error() string {
	return e.String()
}

func (e sqlError) String() string {
	if len(e.args) == 0 {
		return fmt.Sprintf("query: " + e.query)
	}
	return fmt.Sprintf("query: "+e.query+", args: %#v", e.args)
}

type base struct {
	src *sql.DB
}

func open(driver string, config Config) DB {
	var (
		db  DB
		err error
	)

	err = util.Retry(config.Retry, func() error {
		db, err = tryToOpenDB(driver, config.DSN())
		if err != nil {
			fmt.Printf(
				"retry connecting to MySQL... error: %s, host: %s, port: %s, db: %s.",
				err.Error(), config.Host, config.Port, config.Database,
			)
			<-time.After(5 * time.Second)
		}
		return err
	})

	if err != nil {
		log.Print(err.Error())
	}

	return db
}

func tryToOpenDB(driver, dsn string) (DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &base{src: db}, nil
}

func (db *base) Begin(c context.Context, opts *sql.TxOptions) (Tx, error) {
	txn, err := db.src.BeginTx(c, opts)
	if err != nil {
		return nil, err
	}
	return &tx{Tx: txn}, nil
}

func (db *base) Conn(c context.Context) (*sql.Conn, error) {
	return db.src.Conn(c)
}

func (db *base) Close() error {
	return db.src.Close()
}

func (db *base) Driver() driver.Driver {
	return db.src.Driver()
}

func (db *base) Exec(c context.Context, query string, args ...interface{}) (sql.Result, error) {
	r, err := db.src.ExecContext(c, query, args...)
	if err != nil {
		return r, errors.Wrap(&sqlError{query: query, args: args}, err.Error())
	}
	return r, err
}

func (db *base) Ping(c context.Context) error {
	return db.src.PingContext(c)
}

func (db *base) Prepare(c context.Context, query string) Stmt {
	st, err := db.src.PrepareContext(c, query)
	if err != nil {
		panic(errors.Wrap(&sqlError{query: query}, err.Error()))
	}
	return &stmt{Stmt: st}
}

func (db *base) Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	r, err := db.src.QueryContext(c, query, args...)
	if err != nil {
		return r, errors.Wrap(&sqlError{query: query, args: args}, err.Error())
	}
	return r, err
}

func (db *base) QueryRow(c context.Context, query string, args ...interface{}) *sql.Row {
	return db.src.QueryRowContext(c, query, args...)
}

func (db *base) SetConnMaxLifetime(d time.Duration) {
	db.src.SetConnMaxLifetime(d)
}

func (db *base) SetMaxIdleConns(n int) {
	db.src.SetMaxIdleConns(n)
}

func (db *base) SetMaxOpenConns(n int) {
	db.src.SetMaxIdleConns(n)
}

func (db *base) Stats() sql.DBStats {
	return db.src.Stats()
}
