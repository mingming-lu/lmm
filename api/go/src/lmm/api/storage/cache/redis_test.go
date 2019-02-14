package cache

import (
	"context"

	"lmm/api/testing"

	"github.com/gomodule/redigo/redis"
)

func TestRedis(tt *testing.T) {
	c := context.Background()
	t := testing.NewTester(tt)

	client, err := NewRedisClient()
	t.NoError(err)

	conn, err := client.GetContext(c)
	t.NoError(err)

	tt.Run("FLUSHALL", func(tt *testing.T) {
		t := testing.NewTester(tt)
		_, err := conn.Do("FLUSHALL")
		t.NoError(err)
	})

	testKey := "test"
	testVal := "testtest"

	tt.Run("GETBeforeSet", func(tt *testing.T) {
		t := testing.NewTester(tt)

		reply, err := conn.Do("GET", testKey)
		t.NoError(err)
		t.Nil(reply)

	})

	tt.Run("SET", func(tt *testing.T) {
		t := testing.NewTester(tt)

		_, err := conn.Do("SET", testKey, testVal)
		t.NoError(err)

	})

	tt.Run("GETAfterSet", func(tt *testing.T) {
		t := testing.NewTester(tt)

		s, err := redis.String(conn.Do("GET", testKey))
		t.NoError(err)
		t.Is(testVal, s)
	})

	t.NoError(conn.Close(), "close redis connection")
}
