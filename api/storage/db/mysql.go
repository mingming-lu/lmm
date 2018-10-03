package db

import (
	"fmt"
	"os"
)

// MySQL is a DB implementation
type MySQL struct {
	*base
}

// DefaultMySQL returns a new DB with default dsn
func DefaultMySQL() DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return NewMySQL(dsn)
}

// NewMySQL returns a MySQL DB implementation
func NewMySQL(dsn string) DB {
	return newBase("mysql", dsn)
}
