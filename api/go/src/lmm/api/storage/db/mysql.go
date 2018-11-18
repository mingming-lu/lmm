package db

import (
	"os"
	"time"
)

// MySQL is a DB implementation
type MySQL struct {
	*base
}

// DefaultMySQL returns a new DB with default dsn
func DefaultMySQL() DB {
	config := Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASS"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_NAME"),
		Retry:    -1,
	}

	mysql := NewMySQL(config)
	mysql.SetConnMaxLifetime(time.Hour)
	mysql.SetMaxIdleConns(10)
	mysql.SetMaxOpenConns(100)

	return mysql
}

// NewMySQL returns a MySQL DB implementation
func NewMySQL(config Config) DB {
	return newBase("mysql", config)
}
