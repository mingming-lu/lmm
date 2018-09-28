package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	accountFactory "lmm/api/service/account/domain/factory"
	accountRepository "lmm/api/service/account/domain/repository"
	accountStorage "lmm/api/service/account/infra"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/service"
	"lmm/api/storage"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	authorAdapter  service.AuthorService
	userRepository accountRepository.UserRepository
)

func TestMain(m *testing.M) {
	d := storage.NewDB()
	defer d.Close()

	mysql := db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	defer mysql.Close()

	authorAdapter = NewAuthorAdapter(mysql)
	userRepository = accountStorage.NewUserStorage(d)
	code := m.Run()
	os.Exit(code)
}

func TestAuthorFromUserID_OK(tt *testing.T) {
	t := testing.NewTester(tt)

	name, password := uuid.New().String()[:5], uuid.New().String()
	user, err := accountFactory.NewUser(name, password)
	t.NoError(err)
	t.NoError(userRepository.Add(user))

	author, err := authorAdapter.AuthorFromUserID(context.Background(), user.ID())
	t.NoError(err)
	t.Is(user.Name(), author.Name())
}

func TestAuthorFromUserID_NoSuchUser(tt *testing.T) {
	t := testing.NewTester(tt)

	author, err := authorAdapter.AuthorFromUserID(context.Background(), uint64(rand.Intn(10000)))
	t.IsError(domain.ErrNoSuchUser, err)
	t.Nil(author)
}
