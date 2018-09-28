package db

// MySQL is a DB implementation
type MySQL struct {
	*base
}

// NewMySQL returns a MySQL DB implementation
func NewMySQL(dsn string) DB {
	return newBase("mysql", dsn)
}
