package testing

import (
	"github.com/gomodule/redigo/redis"
)

func TestCacheGet(tt *T) {
	t := NewTester(tt)

	conn := cache.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	_, err := conn.Do("PING")
	t.NoError(err)
}

func TestCacheConn_SetString(tt *T) {
	t := NewTester(tt)
	conn := cache.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	s := "ready go"

	_, err := conn.Do("SET", "MYSTERIOUS_KEY", s)
	t.NoError(err)

	sGot, err := redis.String(conn.Do("GET", "MYSTERIOUS_KEY"))
	t.NoError(err)
	t.Is(s, sGot)
}
