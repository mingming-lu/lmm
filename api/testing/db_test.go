package testing

import (
	"lmm/api/db"
	"testing"
)

func TestDBName(t *testing.T) {
	tester := NewTester(t)
	db := db.New()

	var msg string
	err := db.QueryRow("SELECT DATABASE()").Scan(&msg)
	tester.NoError(err)
	tester.Is("lmm_test", msg)
}
