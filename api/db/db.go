package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

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

func (db *DB) Create(database string) *DB {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + database)
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

func (db *DB) Insert(table string, values Values) (sql.Result, error) {
	query := "INSERT INTO " + table + values.String()
	return db.Exec(query)
}
