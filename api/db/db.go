package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEmptyValues   = errors.New("empty values")
	ErrNoChange      = errors.New("no change")
	ErrNoRows        = errors.New(sql.ErrNoRows.Error())
)

var (
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
)

type DB struct {
	*sql.DB
}

func New() *DB {
	conn, err := sql.Open("mysql", dbSrcName+dbName+"?"+connParams)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		panic(err)
	}
	return &DB{DB: conn}
}

func (db *DB) Close() error {
	return errors.New(`
		DB.Close is unexpected to be called after every use.\n
		See https://golang.org/pkg/database/sql/#DB.Close
	`)
}

func (db *DB) CloseNow() error {
	return db.DB.Close()
}

func (db *DB) MustPrepare(query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

func (db *DB) MustPreparef(format string, args ...interface{}) *sql.Stmt {
	query := fmt.Sprintf(format, args...)
	return db.MustPrepare(query)
}
