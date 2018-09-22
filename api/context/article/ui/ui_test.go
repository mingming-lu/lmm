package ui

import (
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
	router            *testing.Router
	user              *auth.User
	articleRepository repository.ArticleRepository
)

func TestMain(m *testing.M) {
	user = testing.NewUser()
	userRepo := authInfra.NewUserStorage(testing.DB())
	auth := authUI.New(userRepo)

	authorService := service.NewAuthorAdapter(testing.DB())
	articleRepository = persistence.NewArticleStorage(testing.DB(), authorService)

	ui := NewUI(articleRepository, authorService)

	router = testing.NewRouter()
	router.POST("/v1/articles", auth.BearerAuth(ui.PostArticle))

	code := m.Run()
	os.Exit(code)
}
