package user

import (
	"lmm/api/domain/service/sha256"
	"lmm/api/testing"
)

func TestNew(t *testing.T) {
	tester := testing.NewTester(t)
	m := New("foobar", "1234")
	tester.Is(m.ID, uint64(0))
	tester.Is(m.Name, "foobar")
	tester.Is(m.Password, sha256.Hex([]byte(m.GUID+"1234")))
}
