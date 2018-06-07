package appservice

import "lmm/api/testing"

func Main(m *testing.M) {
	testing.InitTableAll()
	m.Run()
}
