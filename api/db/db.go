package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/akinaru-lu/errors"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEmptyValues   = errors.New("empty values")
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

type Values struct {
	url.Values
}

func NewValues(values url.Values) *Values {
	return &Values{values}
}

func (values *Values) Where() string {
	s := ""
	for k := range values.Values {
		s += fmt.Sprintf(`%s="%v" AND `, k, values.Get(k))
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

func UseDefault() *DB {
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
		panic(err)
	}
	return db
}

func Init(name string) {
	defaultDatabaseName = name
	d := New().CreateDatabase(defaultDatabaseName).Use(defaultDatabaseName)
	defer d.Close()

	for _, query := range CreateSQL {
		_, err := d.Exec(query)
		if err != nil {
			log.Println(err)
		}
	}
}
