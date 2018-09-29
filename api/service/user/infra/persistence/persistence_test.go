package persistence

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/service/user/domain/repository"
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
	userRepo repository.UserRepository
)

func TestMain(m *testing.M) {
	dbEngine = db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	userRepo = NewUserStorage(dbEngine)

	code := m.Run()
	os.Exit(code)
}
