package testing

func TestCacheGetConn(tt *T) {
	t := NewTester(tt)

	conn := cacheEngine.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	t.NoError(conn.Ping())
}

func TestCacheConn_String(tt *T) {
	t := NewTester(tt)
	conn := cacheEngine.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	var (
		key = "MYSTERIOUS_STR"
		str = "ready go !"
	)

	t.NoError(conn.Set(key, str))

	got, err := conn.GetString(key)
	t.NoError(err)
	t.Is(str, got)
}

func TestCacheConn_Struct(tt *T) {
	t := NewTester(tt)
	conn := cacheEngine.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	type Dummy struct {
		Message string `redis:"message"`
	}

	var (
		key   = "MYSTERIOUS_STRUCT"
		dummy = Dummy{Message: "dummy message"}
	)

	t.NoError(conn.SetStruct(key, &dummy))

	container := Dummy{}
	t.NoError(conn.GetStruct(key, &container))
	t.Is(dummy, container)
}

func TestCacheCoon_Bytes(tt *T) {
	t := NewTester(tt)
	conn := cacheEngine.Get()
	defer func() {
		t.NoError(conn.Close())
	}()

	var (
		key   = "MYSTERIOUS_BYTES"
		bytes = []byte("this is a mysterious string")
	)

	t.NoError(conn.Set(key, bytes))

	got, err := conn.GetBytes(key)
	t.NoError(err)
	t.Is(bytes, got)
}
