package ui

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	accountFactory "lmm/api/service/account/domain/factory"
	auth "lmm/api/service/account/domain/model"
	authRepo "lmm/api/service/account/domain/repository"
	authService "lmm/api/service/account/domain/service"
	authInfra "lmm/api/service/account/infra"
	authUI "lmm/api/service/account/ui"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/infra/fetcher"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/service/article/infra/service"
	"lmm/api/storage"
	featureDB "lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	articleRepository repository.ArticleRepository
	router            *testing.Router
	ui                *UI
	user              *auth.User
	userRepo          authRepo.UserRepository
)

func TestMain(m *testing.M) {
	db := storage.NewDB()
	defer db.Close()
	userRepo = authInfra.NewUserStorage(db)
	auth := authUI.New(userRepo)
	user = NewUser()

	mysql := featureDB.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	defer mysql.Close()

	authorService := service.NewAuthorAdapter(mysql)
	articleFinder := fetcher.NewArticleFetcher(mysql)
	articleRepository = persistence.NewArticleStorage(mysql, authorService)

	ui = NewUI(articleFinder, articleRepository, authorService)

	router = testing.NewRouter()
	router.POST("/v1/articles", auth.BearerAuth(ui.PostArticle))
	router.PUT("/v1/articles/:articleID", auth.BearerAuth(ui.EditArticleText))

	code := m.Run()
	os.Exit(code)
}

func NewUser() *auth.User {
	name, password := uuid.New().String()[:5], uuid.New().String()
	user, err := accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	if err := userRepo.Add(user); err != nil {
		panic(err)
	}

	return auth.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		authService.EncodeToken(user.Token()),
		user.CreatedAt(),
	)
}
