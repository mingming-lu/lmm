package ui

import (
	"os"

	auth "lmm/api/service/account/domain/model"
	authInfra "lmm/api/service/account/infra"
	authUI "lmm/api/service/account/ui"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/infra/fetcher"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/service/article/infra/service"
	"lmm/api/testing"
)

var (
	articleRepository repository.ArticleRepository
	router            *testing.Router
	ui                *UI
	user              *auth.User
)

func TestMain(m *testing.M) {
	user = testing.NewUser()
	userRepo := authInfra.NewUserStorage(testing.DB())
	auth := authUI.New(userRepo)

	authorService := service.NewAuthorAdapter(testing.DB())
	articleFinder := fetcher.NewArticleFetcher(testing.DB())
	articleRepository = persistence.NewArticleStorage(testing.DB(), authorService)

	ui = NewUI(articleFinder, articleRepository, authorService)

	router = testing.NewRouter()
	router.POST("/v1/articles", auth.BearerAuth(ui.PostArticle))
	router.PUT("/v1/articles/:articleID", auth.BearerAuth(ui.EditArticleText))

	code := m.Run()
	os.Exit(code)
}
