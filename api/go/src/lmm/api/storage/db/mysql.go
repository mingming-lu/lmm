package db

import (
	"net/url"
	"os"
	"strconv"
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
		Address:  os.Getenv("MYSQL_ADDR"),
		Protocol: os.Getenv("MYSQL_PROTOCOL"),
		Database: os.Getenv("MYSQL_DBNAME"),
		Retry: func() int {
			if i, err := strconv.Atoi(os.Getenv("MYSQL_CONNECTION_RETRY")); err == nil {
				return i
			}
			return 1
		}(),
	}

	if c.User == "" {
		c.User = "root"
	}

	if c.Protocol == "" {
		c.Protocol = "tcp"
	}

	if c.Address == "" {
		c.Address = "127.0.0.1:3306"
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
