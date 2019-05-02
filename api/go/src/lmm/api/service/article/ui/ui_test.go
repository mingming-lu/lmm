package ui

import (
	"os"
	"sync"

	"lmm/api/http"
	"lmm/api/service/article/infra/fetcher"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/service/article/infra/service"
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence/mysql"
	authUI "lmm/api/service/auth/ui"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysql  db.DB
	router *http.Router
	lock   sync.Mutex
)

func TestMain(m *testing.M) {
	mysql = db.DefaultMySQL()
	auth := auth(mysql)
	ui := articleUI(mysql)
	router = http.NewRouter()
	router.POST("/v1/articles", auth.BearerAuth(ui.PostNewArticle))
	router.PUT("/v1/articles/:articleID", auth.BearerAuth(ui.EditArticle))
	router.GET("/v1/articles", ui.ListArticles)
	router.GET("/v1/articles/:articleID", ui.GetArticle)

	code := m.Run()

	if err := mysql.Close(); err != nil {
		panic(err)
	}

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

func getArticleByID(articleID string) *testing.Response {
	request := testing.GET("/v1/articles/"+articleID, nil)
	return testing.DoRequest(request, router)
}
