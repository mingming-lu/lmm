package db

import (
	"os"
)

// MySQL is a DB implementation
type MySQL struct {
	*base
}

// DefaultMySQL returns a new DB with default dsn
func DefaultMySQL() DB {
	config := Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Retry:    -1,
	}
	return NewMySQL(config)
}

// NewMySQL returns a MySQL DB implementation
func NewMySQL(config Config) DB {
	return newBase("mysql", config)
}
