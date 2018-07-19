package service

import (
	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/account/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *account.User

func TestMain(m *testing.M) {
	name, password := uuid.New()[:31], uuid.New()
	user, _ = accountFactory.NewUser(name, password)
	infra.NewUserStorage(testing.DB()).Add(user)

	code := m.Run()
	os.Exit(code)
}
