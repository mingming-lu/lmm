package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/service"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
)

func TestSaveUser(tt *testing.T) {
	c := context.Background()

	builder := factory.NewFactory(&service.BcryptService{})

	username := "U" + stringutil.ReplaceAll(uuid.New().String(), "-", "")[:8]
	email := "admin@lmm.local"
	password := "notweakpassword123!?"
	user, err := builder.NewUser(username, email, password)

	if err != nil {
		tt.Fatal(err)
	}

	tt.Run("Success", func(tt *testing.T) {
		t := testing.NewTester(tt)

		t.NoError(userRepo.Save(c, user))

		var (
			nameFound     string
			emailFound    string
			passwordFound string
			tokenFound    string
		)
		t.NoError(
			dbEngine.QueryRow(c,
				`select name, email, password, token from user where name = ?`, username,
			).Scan(&nameFound, &emailFound, &passwordFound, &tokenFound),
		)
		t.Is(username, nameFound)
		t.Is(email, emailFound)
		t.NoError(bcrypt.CompareHashAndPassword([]byte(passwordFound), []byte(password)))
	})

	tt.Run("DuplicateUserName", func(tt *testing.T) {
		t := testing.NewTester(tt)

		t.IsError(
			domain.ErrUserNameAlreadyUsed,
			errors.Cause(userRepo.Save(c, user)),
		)
	})
}
