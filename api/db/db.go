package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/akinaru-lu/errors"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEmptyValues = errors.New("empty values")
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

func (values Values) String() string {
	cs := " ("
	vs := " VALUES ("
	for k, v := range values {
		cs += k + ","
		if _, ok := v.(string); ok {
			vs += fmt.Sprintf(`"%s",`, v.(string))
		} else {
			vs += fmt.Sprintf(`%v,`, v)
		}
	}
	if strings.HasSuffix(cs, ",") {
		cs = cs[:len(cs)-1]
	}
	if strings.HasSuffix(vs, ",") {
		vs = vs[:len(vs)-1]
	}
	cs += ")"
	vs += ")"
	return cs + vs
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

func (db *DB) Insert(table string, values ...Values) (sql.Result, error) {
	if len(values) == 0 {
		return nil, ErrEmptyValues
	} else if len(values) == 1 {
		return db.Exec("INSERT INTO " + table + values[0].String())
	} else {
		return nil, nil
	}
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
