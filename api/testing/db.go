package testing

import (
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
