package persistence

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	accountFactory "lmm/api/service/account/domain/factory"
	account "lmm/api/service/account/domain/model"
	accountModel "lmm/api/service/account/domain/model"
	accountService "lmm/api/service/account/domain/service"
	accountStorage "lmm/api/service/account/infra"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/service"
	infraService "lmm/api/service/article/infra/service"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	articleRepository repository.ArticleRepository
	articleService    *service.ArticleService
	user              *account.User
	authorService     service.AuthorService
)

func TestMain(m *testing.M) {
	mysql := db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))

	authorService = infraService.NewAuthorAdapter(mysql)
	articleRepository = NewArticleStorage(mysql, authorService)
	articleService = service.NewArticleService(articleRepository)
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
