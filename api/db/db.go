package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/akinaru-lu/errors"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEmptyValues   = errors.New("empty values")
	ErrNoChange      = errors.New("no change")
	ErrNoRows        = errors.New(sql.ErrNoRows.Error())
)

var defaultDatabaseName = ""
var mux sync.Mutex

func SetDefaultDatabaseName(name string) {
	mux.Lock()
	defaultDatabaseName = name
	mux.Unlock()
}

type DB struct {
	*sql.DB
}

type Values map[string]interface{}

func NewValues() Values {
	return make(Values)
}

func NewValuesFromURL(values url.Values) Values {
	ret := NewValues()
	for k := range values {
		ret[k] = values.Get(k)
	}
	return ret
}

func (values Values) Where() string {
	s := ""
	for k, v := range values {
		s += fmt.Sprintf(`%s="%v" AND `, k, v)
	}
	if s != "" {
		s = "WHERE " + strings.TrimSuffix(s, " AND ")
	}
	return s
}

func New() *DB {
	super, err := sql.Open("mysql", "root:@/")
	if err != nil {
		panic(err)
	}
	return &DB{DB: super}
}

func Default() *DB {
	if defaultDatabaseName == "" {
		panic("Default database has not been set")
	}
	return New().Use(defaultDatabaseName)
}

func (db *DB) CreateDatabase(name string) *DB {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}
	return db
}

func (db *DB) DropDatabase(name string) *DB {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		panic(err)
	}
	return db
}

func (db *DB) Use(database string) *DB {
	_, err := db.Exec("USE " + database)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (db *DB) Must(query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

func (db *DB) Mustf(format string, args ...interface{}) *sql.Stmt {
	query := fmt.Sprintf(format, args...)
	return db.Must(query)
}

func (db *DB) Exists(query string, args ...interface{}) (bool, error) {
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)

	var exists bool
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}
