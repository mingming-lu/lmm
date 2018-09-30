package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/auth/domain/model"
	"lmm/api/service/auth/domain/repository"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
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

func newUser() *model.User {
	rawPassword := uuid.New().String()
	b, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	encryptedPassword := string(b)

	user := model.NewUser(
		uuid.New().String()[:8],
		encryptedPassword,
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
	)

	now := time.Now()

	if _, err := dbEngine.Exec(context.Background(), `
		insert into user (name, password, token, created_at) values (?, ?, ?, ?)
	`, user.Name(), encryptedPassword, user.RawToken(), now); err != nil {
		panic(err)
	}

	return user
}
