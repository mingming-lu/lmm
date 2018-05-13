package model

import (
	"lmm/api/testing"
	"lmm/api/utils/sha256"
)

func TestNew(t *testing.T) {
	tester := testing.NewTester(t)
	m := New("foobar", "1234")
	tester.Is(uint64(0), m.ID)
	tester.Is("foobar", m.Name)
	tester.Is(sha256.Hex([]byte(m.GUID+"1234")), m.Password)
}
