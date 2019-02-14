package cache

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	redisServer           = "lmm-api-redis:6379"
	redisConnFetchTimeout = 10 * time.Second
)

// RedisClient is a Redis client using redigo
type RedisClient struct {
	redis.Pool
}

// NewRedisClient creates a new RedisClient
func NewRedisClient() (*RedisClient, error) {
	client := &RedisClient{
		Pool: redis.Pool{
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisServer) },
			IdleTimeout: 5 * time.Minute,
			MaxIdle:     1,
			MaxActive:   10,
			Wait:        true,

			TestOnBorrow: func(c redis.Conn, elapsed time.Time) error {
				if time.Since(elapsed) < redisConnFetchTimeout {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), redisConnFetchTimeout)
	var errChan chan error
	go func() {
		conn, err := client.GetContext(ctx)
		if err != nil {
			errChan <- err
			return
		}
		if _, err := conn.Do("PING"); err != nil {
			errChan <- err
			return
		}
		if err := conn.Close(); err != nil {
			errChan <- err
			return
		}
		cancel()
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ctx.Err()
		}
	case err := <-errChan:
		return nil, err
	}

	return client, nil
}
