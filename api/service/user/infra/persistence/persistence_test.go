package persistence

import (
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/service/user/domain/repository"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	dbEngine db.DB
	userRepo repository.UserRepository
)

func TestMain(m *testing.M) {
	dbEngine = db.DefaultMySQL()
	defer dbEngine.Close()

	userRepo = NewUserStorage(dbEngine)

	code := m.Run()
	os.Exit(code)
}
