package storage

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

func NewCacheEngine() *Cache {
	return &Cache{pool: newPool()}
}

func (c *Cache) Get() redis.Conn {
	return c.pool.Get()
}

func (c *Cache) WithContext(ctx context.Context) (redis.Conn, error) {
	return c.pool.GetContext(ctx)
}

func (c *Cache) ActiveCount() int {
	return c.pool.ActiveCount()
}

func (c *Cache) IdelCount() int {
	return c.pool.IdleCount()
}

func (c *Cache) Close() error {
	return c.pool.Close()
}

const (
	cacheConnIdleTimeOut = 300 * time.Second
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: cacheConnIdleTimeOut,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "lmm-redis:6379")
			if err != nil {
				panic(err.Error())
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < cacheConnIdleTimeOut {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}
