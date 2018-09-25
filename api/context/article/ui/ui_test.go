package ui

import (
	"lmm/api/context/article/infra/fetcher"
	"os"

	auth "lmm/api/context/account/domain/model"
	authInfra "lmm/api/context/account/infra"
	authUI "lmm/api/context/account/ui"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/infra/persistence"
	"lmm/api/context/article/infra/service"
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
