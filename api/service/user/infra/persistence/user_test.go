package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
)

func TestSaveUser(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	token := stringutil.ReplaceAll(uuid.New().String(), "-", "")

	pw, err := model.NewPassword("notweakpassword123!?")
	if !t.NoError(err) {
		t.FailNow()
	}

	username := "U" + stringutil.ReplaceAll(uuid.New().String(), "-", "")[:8]

	user, err := model.NewUser(username, *pw, token)
	if !t.NoError(err) {
		t.FailNow()
	}

	t.Run("Success", func(_ *testing.T) {
		t.NoError(userRepo.Save(c, user))

		var (
			nameFound     string
			passwordFound string
			tokenFound    string
		)
		t.NoError(
			dbEngine.QueryRow(c,
				`select name, password, token from user where name = ?`, username,
			).Scan(&nameFound, &passwordFound, &tokenFound),
		)
		t.Is(username, nameFound)
		t.NoError(bcrypt.CompareHashAndPassword([]byte(passwordFound), []byte(pw.String())))
		t.Is(token, tokenFound)
	})

	t.Run("DuplicateUserName", func(_ *testing.T) {
		t.IsError(
			domain.ErrUserNameAlreadyUsed,
			errors.Cause(userRepo.Save(c, user)),
		)
	})
}
