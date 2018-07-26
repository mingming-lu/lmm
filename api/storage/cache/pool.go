package cache

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	cacheConnIdleTimeOut = 300 * time.Second
)

type ConnPool struct {
	*redis.Pool
}

func NewPool() *ConnPool {
	pool := redis.Pool{
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
	return &ConnPool{Pool: &pool}
}

func (p *ConnPool) Get() Conn {
	c := p.Pool.Get()
	return &conn{Conn: c}
}

func (p *ConnPool) GetContext(ctx context.Context) (Conn, error) {
	c, err := p.Pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{Conn: c}, err
}
