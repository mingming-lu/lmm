package appservice

import (
	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/context/blog/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var (
	app  *AppService
	user *account.User
)

func TestMain(m *testing.M) {
	// init app
	db := testing.DB()
	blogRepo := infra.NewBlogStorage(db)
	categoryRepo := infra.NewCategoryStorage(db)
	app = New(blogRepo, categoryRepo)

	// init user
	name, password := uuid.New()[:31], uuid.New()
	user, _ = accountFactory.NewUser(name, password)
	accountRepository.New(testing.DB()).Add(user)

	code := m.Run()
	os.Exit(code)
}
