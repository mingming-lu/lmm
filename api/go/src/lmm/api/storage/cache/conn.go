package cache

import "github.com/gomodule/redigo/redis"

type Conn struct {
	redis.Conn
}

func (c *Conn) Ping() error {
	_, err := c.Do("PING")
	return err
}

func (c *Conn) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(c.Do("GET", key))
}

func (c *Conn) GetString(key string) (string, error) {
	return redis.String(c.Do("GET", key))
}

func (c *Conn) GetStruct(key string, dest interface{}) error {
	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return redis.ScanStruct(values, dest)
}

func (c *Conn) Set(key string, value interface{}) error {
	_, err := c.Do("SET", key, value)
	return err
}

func (c *Conn) SetTTL(key string, value interface{}, seconds uint) error {
	_, err := c.Do("SETEX", seconds, value)
	return err
}

func (c *Conn) SetStruct(key string, value interface{}) error {
	_, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	return err
}

func (c *Conn) SetStructTTL(key string, value interface{}, seconds uint) error {
	if err := c.Send("HMSET", redis.Args{}.Add(key).AddFlat(value)...); err != nil {
		return err
	}
	if err := c.Send("EXPIRE", key, seconds); err != nil {
		return err
	}
	return c.Flush()
}
