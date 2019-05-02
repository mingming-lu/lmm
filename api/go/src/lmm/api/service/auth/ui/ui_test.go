package ui

import (
	"os"

	"lmm/api/http"
	"lmm/api/service/auth/application"
	"lmm/api/service/auth/infra/persistence/mysql"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbEngine db.DB
	router   *http.Router
)

func TestMain(m *testing.M) {
	dbEngine = db.DefaultMySQL()
	repo := mysql.NewUserStorage(dbEngine)
	app := application.NewService(repo)
	ui := NewUI(app)
	router = http.NewRouter()
	router.POST("/v1/auth/login", ui.Login)
	router.GET("/v1/dummy", ui.BearerAuth(func(c http.Context) {
		c.String(http.StatusOK, http.StatusText(http.StatusOK))
	}))

	code := m.Run()

	if err := dbEngine.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
