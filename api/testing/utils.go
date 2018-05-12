package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"lmm/api/db"
	"log"
)

func InitTableAll() {
	db := db.Default()
	defer db.Close()
	r, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	for r.Next() {
		var table string
		r.Scan(&table)
		db.Query("TRUNCATE TABLE " + table)
	}
}

func InitTable(name string) {
	db := db.Default()
	defer db.Close()

	_, err := db.Query("TRUNCATE TABLE " + name)
	if err != nil {
		log.Fatal(err)
	}
}

func StructToRequestBody(o interface{}) io.ReadCloser {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}
