package service

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/service"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

var (
	mysql         db.DB
	authorAdapter service.AuthorService
)

func TestMain(m *testing.M) {
	mysql = db.DefaultMySQL()
	defer mysql.Close()

	authorAdapter = NewAuthorAdapter(mysql)
	code := m.Run()
	os.Exit(code)
}

func TestAuthorFromUserID_OK(tt *testing.T) {
	t := testing.NewTester(tt)

	user := testutil.NewUser(mysql)

	author, err := authorAdapter.AuthorFromUserName(context.Background(), user.Name())
	t.NoError(err)
	t.Is(user.Name(), author.Name())
}

func TestAuthorFromUserID_NoSuchUser(tt *testing.T) {
	t := testing.NewTester(tt)

	author, err := authorAdapter.AuthorFromUserName(context.Background(), "U"+uuid.New().String()[:8])
	t.IsError(domain.ErrNoSuchUser, err)
	t.Nil(author)
}
