package rand

import "lmm/api/testing"

func TestDefault_Success(tt *testing.T) {
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
