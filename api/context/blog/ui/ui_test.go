package ui

import (
	accountFactory "lmm/api/context/account/domain/factory"
	accountModel "lmm/api/context/account/domain/model"
	accountService "lmm/api/context/account/domain/service"
	accountStorage "lmm/api/context/account/infra"
	account "lmm/api/context/account/ui"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *accountModel.User
var ui *UI
var accountUI *account.UI

func TestMain(m *testing.M) {
	initUser()
	initUI()

	accountUI = account.New(accountStorage.NewUserStorage(testing.DB()))

	code := m.Run()
	os.Exit(code)
}

func initUser() {
	var err error

	name, password := uuid.New()[:31], uuid.New()
	user, err = accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	err = accountStorage.NewUserStorage(testing.DB()).Add(user)
	if err != nil {
		panic(err)
	}

	user = accountModel.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		accountService.EncodeToken(user.Token()),
		user.CreatedAt(),
	)
}

func initUI() {
	db := testing.DB()
	blogRepo := infra.NewBlogStorage(db)
	categoryRepo := infra.NewCategoryStorage(db)
	tagRepo := infra.NewTagStorage(db)
	ui = New(blogRepo, categoryRepo, tagRepo)
}

func randomBlog() *model.Blog {
	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	if err != nil {
		panic(err)
	}

	repo := infra.NewBlogStorage(testing.DB())
	if err := repo.Add(blog); err != nil {
		panic(err)
	}

	return blog
}
