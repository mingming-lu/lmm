package persistence

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

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
	authorService     service.AuthorService
	mysql             db.DB
)

func TestMain(m *testing.M) {
	mysql = db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	defer mysql.Close()

	authorService = infraService.NewAuthorAdapter(mysql)
	articleRepository = NewArticleStorage(mysql, authorService)
	articleService = service.NewArticleService(articleRepository)

	code := m.Run()
	os.Exit(code)
}
