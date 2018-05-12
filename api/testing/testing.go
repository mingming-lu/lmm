package testing

import "lmm/api/db"

func init() {
	db.SetDefaultDatabaseName("lmm_test")
}
