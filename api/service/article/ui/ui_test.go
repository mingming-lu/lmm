package ui

import (
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/service/article/infra/fetcher"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/service/article/infra/service"
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence"
	authUI "lmm/api/service/auth/ui"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	mysql  db.DB
	router *testing.Router
)

func TestMain(m *testing.M) {
	mysql = db.DefaultMySQL()
	defer mysql.Close()

	auth := auth(mysql)
	ui := articleUI(mysql)

	router = testing.NewRouter()
	router.POST("/v1/articles", auth.BearerAuth(ui.PostNewArticle))
	router.PUT("/v1/articles/:articleID", auth.BearerAuth(ui.EditArticle))

	code := m.Run()
	os.Exit(code)
}

func articleUI(db db.DB) *UI {
	authorService := service.NewAuthorAdapter(db)
	articleFinder := fetcher.NewArticleFetcher(db)
	articleRepository := persistence.NewArticleStorage(db, authorService)
	return NewUI(articleFinder, articleRepository, authorService)
}

func auth(db db.DB) *authUI.UI {
	repo := authStorage.NewUserStorage(db)
	app := authApp.NewService(repo)
	return authUI.NewUI(app)
}
