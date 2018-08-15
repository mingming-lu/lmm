package rand

import (
	"lmm/api/testing"
	"time"
)

func TestDefault_Output(tt *testing.T) {
	t := testing.NewTester(tt)
	t.Log("=========base36=========")
	t.Log(Base36(8))
	t.Log(Base36(8))
	t.Log(Base36(8))
	t.Log()

	t.Log("=========base62=========")
	t.Log(Base62(8))
	t.Log(Base62(8))
	t.Log(Base62(8))
	t.Log()

	t.Log("=========base64=========")
	t.Log(Base64(8))
	t.Log(Base64(8))
	t.Log(Base64(8))
	t.Log()
}

func TestNew(tt *testing.T) {
	t := testing.NewTester(tt)

	r := New(0)

	Seed(0)
	defer Seed(time.Now().UnixNano())

	t.Is(Base36(5), r.Base36(5))
	t.Is(Base62(9), r.Base62(9))
	t.Is(Base36(6), r.Base36(6))
}

func TestBase(tt *testing.T) {
	t := testing.NewTester(tt)

	b := base(defaultRand, 6, 62)
	t.Is(6, len(b))
	t.Is(6, cap(b))
}
