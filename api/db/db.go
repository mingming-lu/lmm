package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/akinaru-lu/errors"
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
	return &DB{DB: conn}
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
