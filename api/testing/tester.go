package testing

import (
	"lmm/api/db"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	db.SetDefaultDatabaseName("lmm_test")
}

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

type T = testing.T

type Tester struct {
	*T
}

func NewTester(t *T) *Tester {
	return &Tester{T: t}
}

func (t *Tester) Is(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Equal(t, expected, actual, msgAndArgs...)
}

func (t *Tester) Not(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotEqual(t, expected, actual, msgAndArgs...)
}

func (t *Tester) Error(err error, msgAndArgs ...interface{}) bool {
	return assert.Error(t, err, msgAndArgs...)
}

func (t *Tester) NoError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(t, err, msgAndArgs...)
}
