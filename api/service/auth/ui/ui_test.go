package ui

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/http"
	"lmm/api/service/auth/infra/persistence"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	dbEngine db.DB
	router   *http.Router
)

func TestMain(m *testing.M) {
	dbEngine = db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	userRepo := persistence.NewUserStorage(dbEngine)
	ui := NewUI(userRepo)

	router = http.NewRouter()
	router.POST("/v1/auth/login", ui.Login)

	code := m.Run()
	os.Exit(code)
}
