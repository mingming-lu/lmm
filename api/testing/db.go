package testing

import "go.uber.org/zap"

func InitTableAll() {
	r, err := dbEngine.Query("SHOW TABLES")
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	for r.Next() {
		var table string
		r.Scan(&table)
		dbEngine.Query("TRUNCATE TABLE " + table)
	}
}

func InitTable(name string) {
	_, err := dbEngine.Query("TRUNCATE TABLE " + name)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}
