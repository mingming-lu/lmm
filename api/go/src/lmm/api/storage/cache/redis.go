package cache

import (
	"context"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"

	"lmm/api/util/stringutil"
)

var (
	redisServer           string
	redisConnFetchTimeout time.Duration
)

func init() {
	redisServer = os.Getenv("REDIS_SERVER")

	second, err := stringutil.ParseUint(os.Getenv("REDIS_CONN_TIMEOUT_SECOND"))
	if err != nil {
		panic(err)
	}
	redisConnFetchTimeout = time.Duration(second) * time.Second
}

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
