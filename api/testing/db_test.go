package testing

import (
	"testing"
)

func TestDBName(t *testing.T) {
	tester := NewTester(t)

	var msg string
	err := db.QueryRow("SELECT DATABASE()").Scan(&msg)
	tester.NoError(err)
	tester.Is("lmm_test", msg)
}
