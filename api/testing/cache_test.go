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

func TestCacheConn_ScanStruct(tt *T) {
	t := NewTester(tt)
	conn := cache.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	type Dummy struct {
		Message string `redis:"message"`
	}

	dummy := Dummy{Message: "dummy message"}

	_, err := conn.Do("HSET", redis.Args{}.Add("MYSTERIOUS_STRUCT").AddFlat(&dummy)...)
	t.NoError(err)

	reply, err := conn.Do("HGETALL", "MYSTERIOUS_STRUCT")
	t.NoError(err)

	values, err := redis.Values(reply, err)
	t.NoError(err)

	scanned := Dummy{}
	t.NoError(redis.ScanStruct(values, &scanned))
	t.Is(dummy, scanned)
}
