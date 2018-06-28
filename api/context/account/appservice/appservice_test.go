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

func randomUserNameAndPassword() Auth {
	return Auth{
		Name:     uuid.New()[:31],
		Password: uuid.New(),
	}
}
