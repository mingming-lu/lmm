package persistence

import (
	"os"

	"lmm/api/service/auth/domain/repository"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbEngine db.DB
	userRepo repository.UserRepository
)

func TestMain(m *testing.M) {
	dbEngine = db.DefaultMySQL()
	userRepo = NewUserStorage(dbEngine)

	code := m.Run()

	if err := dbEngine.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
