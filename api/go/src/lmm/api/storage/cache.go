package storage

import (
	"context"

	"lmm/api/storage/cache"
)

type Cache struct {
	pool *cache.ConnPool
}

func NewCacheEngine() *Cache {
	return &Cache{pool: cache.NewPool()}
}

func (c *Cache) Get() *cache.Conn {
	return c.pool.Get()
}

func (c *Cache) WithContext(ctx context.Context) (*cache.Conn, error) {
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
