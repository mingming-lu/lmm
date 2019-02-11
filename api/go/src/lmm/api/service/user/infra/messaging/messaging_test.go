package messaging

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/messaging"
	userEvent "lmm/api/service/user/domain/event"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/service"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

var (
	builder        *factory.Factory
	encrypter      service.EncryptService
	mysql          db.DB
	testSubscriber *Subscriber
)

func TestMain(m *testing.M) {
	mysql = db.DefaultMySQL()

	encrypter = &service.BcryptService{}
	builder = factory.NewFactory(encrypter)
	testSubscriber = NewSubscriber(mysql)
	messaging.SyncBus().Subscribe(&userEvent.UserPasswordChanged{}, testSubscriber.OnUserPasswordChanged)
	messaging.SyncBus().Subscribe(&userEvent.UserRoleChanged{}, testSubscriber.OnUserRoleChanged)

	code := m.Run()

	mysql.Close()

	os.Exit(code)
}

func TestOnUserPasswordChanged(tt *testing.T) {
	c := context.Background()
	t := testing.NewTester(tt)

	user := newUser()
	newRawPassword := "c0mp|exP@s$worD"
	newPassword, err := builder.NewPassword(newRawPassword)
	if err != nil {
		tt.Fatal(err)
	}
	user.ChangePassword(newPassword)

	t.NoError(userEvent.PublishUserPasswordChanged(c, user.Name(), user.Password()))

	var (
		password string
		token    string
	)

	t.NoError(mysql.QueryRow(c, "select password, token from user where name = ?", user.Name()).Scan(&password, &token))
	t.True(encrypter.Verify(newRawPassword, password))
	t.Not(user.Token(), token)
}

func newUser() *model.User {
	testuser := testutil.NewUser(mysql)

	user, err := model.NewUser(testuser.Name(), testuser.Email(), testuser.EncryptedPassword(), testuser.RawToken(), testuser.Role(), testuser.CreatedAt())
	if err != nil {
		panic(err)
	}

	return user
}
