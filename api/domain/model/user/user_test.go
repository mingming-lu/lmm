package user

import (
	"lmm/api/domain/service/sha256"
	"lmm/api/testing"
)

func TestNew(t *testing.T) {
	tester := testing.NewTester(t)
	m := New("foobar", "1234")
	tester.Is(uint64(0), m.ID)
	tester.Is("foobar", m.Name)
	tester.Is(sha256.Hex([]byte(m.GUID+"1234")), m.Password)
}
