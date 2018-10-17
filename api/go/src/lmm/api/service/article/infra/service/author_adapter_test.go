package service

import (
	"context"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/service"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/testutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var (
	mysql         db.DB
	authorAdapter service.AuthorService
)

func TestMain(m *testing.M) {
	testing.NewTestRunner(m).Setup(func() {
		mysql = db.DefaultMySQL()
		authorAdapter = NewAuthorAdapter(mysql)
	}).Teardown(func() {
		mysql.Close()
	}).Run()
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
