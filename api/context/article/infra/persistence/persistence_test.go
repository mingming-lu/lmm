package persistence

import (
	"os"

	"github.com/google/uuid"

	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	accountModel "lmm/api/context/account/domain/model"
	accountService "lmm/api/context/account/domain/service"
	accountStorage "lmm/api/context/account/infra"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
	infraService "lmm/api/context/article/infra/service"
	"lmm/api/storage"
	"lmm/api/testing"
)

var (
	articleRepository repository.ArticleRepository
	articleService    *service.ArticleService
	user              *account.User
	authorService     service.AuthorService
)

func TestMain(m *testing.M) {
	db := storage.NewDB()
	defer db.CloseNow()

	articleRepository = NewArticleStorage(db, authorService)
	articleService = service.NewArticleService(articleRepository)
	authorService = infraService.NewAuthorAdapter(db)
	user = initUser()

	code := m.Run()
	os.Exit(code)
}

func initUser() *account.User {
	var err error

	name, password := uuid.New().String()[:5], uuid.New().String()
	user, err = accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	err = accountStorage.NewUserStorage(testing.DB()).Add(user)
	if err != nil {
		panic(err)
	}

	return accountModel.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		accountService.EncodeToken(user.Token()),
		user.CreatedAt(),
	)
}
