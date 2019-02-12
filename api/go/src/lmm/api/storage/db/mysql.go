package db

import (
	"net/url"
	"os"
	"time"
)

// MySQL is a DB implementation
type MySQL struct {
	*base
}

// DefaultMySQLConfig returns default MySQL Config
func DefaultMySQLConfig() Config {
	c := Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASS"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_NAME"),
		Retry:    -1,
	}

	if c.User == "" {
		c.User = "root"
	}

	if c.Protocol == "" {
		c.Protocol = "tcp"
	}

	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	if c.Port == "" {
		c.Port = "3306"
	}

	c.Options = url.Values{}
	c.Options.Add("parseTime", "true")
	c.Options.Add("autocommit", "true")
	c.Options.Add("tx_isolation", "'READ-COMMITTED'")
	c.Options.Add("charset", "utf8mb4")

	return c
}

// DefaultMySQL returns a new DB with default dsn
func DefaultMySQL() DB {
	c := DefaultMySQLConfig()

	mysql := NewMySQL(c)
	mysql.SetConnMaxLifetime(time.Hour)
	mysql.SetMaxIdleConns(10)
	mysql.SetMaxOpenConns(100)

	return mysql
}

// NewMySQL returns a MySQL DB implementation
func NewMySQL(config Config) DB {
	return open("mysql", config)
}
