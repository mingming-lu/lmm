package persistence

import (
	"lmm/api/service/user/domain/repository"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbEngine db.DB
	userRepo repository.UserRepository
)

func TestMain(m *testing.M) {
	testing.NewTestRunner(m).Setup(func() {
		dbEngine = db.DefaultMySQL()
		userRepo = NewUserStorage(dbEngine)
	}).Teardown(func() {
		dbEngine.Close()
	}).Run()
}
