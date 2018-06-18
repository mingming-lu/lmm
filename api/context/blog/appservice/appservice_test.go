package appservice

import (
	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *account.User

func TestMain(m *testing.M) {
	name, password := uuid.New()[:31], uuid.New()
	user, _ = accountFactory.NewUser(name, password)
	accountRepository.New().Add(user)

	code := m.Run()
	os.Exit(code)
}
