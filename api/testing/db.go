package testing

import (
	"log"
)

func InitTableAll() {
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
	_, err := db.Query("TRUNCATE TABLE " + name)
	if err != nil {
		log.Fatal(err)
	}
}
