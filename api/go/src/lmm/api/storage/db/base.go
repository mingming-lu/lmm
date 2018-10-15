package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDonw   = sql.ErrTxDone
)

type base struct {
	src *sql.DB
}

func newBase(driver, dsn string) DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	return &base{src: db}
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
	return db.src.ExecContext(c, query, args...)
}

func (db *base) Ping(c context.Context) error {
	return db.src.PingContext(c)
}

func (db *base) Prepare(c context.Context, query string) Stmt {
	st, err := db.src.PrepareContext(c, query)
	if err != nil {
		panic(err)
	}
	return &stmt{Stmt: st}
}

func (db *base) Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.src.QueryContext(c, query, args...)
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
