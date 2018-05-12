package tester

import (
	"lmm/api/db"
	"testing"

	"github.com/akinaru-lu/errors"
)

func TestInit(t *testing.T) {
	tester := New(t)
	db := db.Default()
	_, err := db.Query("SHOW TABLES")
	tester.NoError(err)
}

func TestIs(t *testing.T) {
	tester := New(t)
	tester.Is(1, 1)
	tester.Is(true, true)
	tester.Is("abc", "abc")
}

func TestNot(t *testing.T) {
	tester := New(t)
	tester.Not(1, 2)
	tester.Not(true, false)
	tester.Not(1, "1")
	tester.Not(tester, t)
}

func TestError(t *testing.T) {
	tester := New(t)
	tester.Error(errors.New(""))
	tester.Error(errors.New("msg"))
}

func TestNoError(t *testing.T) {
	tester := New(t)
	tester.NoError(nil)
}
