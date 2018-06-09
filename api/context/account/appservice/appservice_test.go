package appservice

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

func Main(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func randomUserNameAndPassword() (string, string) {
	return uuid.New()[:31], uuid.New()
}
