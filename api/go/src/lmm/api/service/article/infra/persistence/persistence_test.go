package persistence

import (
	"os"

	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/service"
	infraService "lmm/api/service/article/infra/service"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	articleRepository repository.ArticleRepository
	articleService    *service.ArticleService
	authorService     service.AuthorService
	mysql             db.DB
)

func TestMain(m *testing.M) {
	mysql = db.DefaultMySQL()
	authorService = infraService.NewAuthorAdapter(mysql)
	articleRepository = NewArticleStorage(mysql, authorService)
	articleService = service.NewArticleService(articleRepository)

	code := m.Run()

	if err := mysql.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
