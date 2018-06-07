package appservice

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func Main(m *testing.M) {
	m.Run()
}

func randomUserNameAndPassword() (string, string) {
	return uuid.New()[:32], uuid.New()
}
