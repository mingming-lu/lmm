package cache

import "github.com/gomodule/redigo/redis"

type Conn interface {
	redis.Conn

	Ping() error

	GetString(key string) (string, error)
	SetString(key, value string) error

	GetStruct(key string, dest interface{}) error
	SetStruct(key string, value interface{}) error
}

type conn struct {
	redis.Conn
}

func (c *conn) Ping() error {
	_, err := c.Do("PING")
	return err
}

func (c *conn) GetString(key string) (string, error) {
	return redis.String(c.Do("GET", key))
}

func (c *conn) SetString(key, value string) error {
	_, err := c.Do("SET", key, value)
	return err
}

func (c *conn) GetStruct(key string, dest interface{}) error {
	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return redis.ScanStruct(values, dest)
}

func (c *conn) SetStruct(key string, value interface{}) error {
	_, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	return err
}
