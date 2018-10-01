package service

import (
	"context"
	"fmt"
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
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	mysql         db.DB
	authorAdapter service.AuthorService
)

func TestMain(m *testing.M) {
	mysql = db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	defer mysql.Close()

	authorAdapter = NewAuthorAdapter(mysql)
	code := m.Run()
	os.Exit(code)
}

func TestAuthorFromUserID_OK(tt *testing.T) {
	t := testing.NewTester(tt)

	name, password := "U"+uuid.New().String()[:5], uuid.New().String()
	user := testutil.NewUserUser(mysql, name, password)

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
