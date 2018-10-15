package application

import (
	"context"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"
	"lmm/api/util/uuidutil"
)

var (
	testAppService *Service
)

type InmemoryUserRepository struct {
	memory []*model.User
}

func (repo *InmemoryUserRepository) Save(c context.Context, user *model.User) error {
	if user, _ := repo.FindByUserName(c, user.Name()); user != nil {
		return domain.ErrUserNameAlreadyUsed
	}
	repo.memory = append(repo.memory, user)
	return nil
}

func (repo *InmemoryUserRepository) FindByUserName(c context.Context, username string) (*model.User, error) {
	for _, user := range repo.memory {
		if user.Name() == username {
			return user, nil
		}
	}
	return nil, domain.ErrNoSuchUser
}

func TestMain(m *testing.M) {
	testAppService = NewService(&InmemoryUserRepository{memory: make([]*model.User, 0)})
	code := m.Run()
	os.Exit(code)
}

func TestRegisterNewUser(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	t.Run("Success", func(_ *testing.T) {
		username, password := "username", "~!@#$%^&*()-_=+{[}]|\\:;\"'<,>.?/"
		nameGot, err := testAppService.RegisterNewUser(c, username, password)
		t.NoError(err)
		t.Is(username, nameGot)

		user, err := testAppService.userRepository.FindByUserName(c, "username")
		t.NoError(err)
		t.Is(username, user.Name())
		t.NoError(bcrypt.CompareHashAndPassword(
			[]byte(user.Password()),
			[]byte(password),
		))
		t.NotPanic(func() {
			uuid.Must(uuidutil.ParseString(user.Token()))
		})
	})

	t.Run("Fail", func(_ *testing.T) {
		cases := map[string]struct {
			UserName string
			Password string
			Err      error
		}{
			"UserNameTooShort": {
				"ur", "password1234", domain.ErrInvalidUserName,
			},
			"UserNameStartsWithoutLetter": {
				"1username", "password1234", domain.ErrInvalidUserName,
			},
			"EmptyPassword": {
				"username", "", domain.ErrUserPasswordEmpty,
			},
			"PasswordIsTooShort": {
				"username", "passwor", domain.ErrUserPasswordTooShort,
			},
			"PasswordIsTooWeak": {
				"username", "password", domain.ErrUserPasswordTooWeak,
			},
			"PasswordIsTooLong": {
				"username", strings.Repeat("s", 251), domain.ErrUserPasswordTooLong,
			},
			"DuplicateUserName": {
				"username", "password1234", domain.ErrUserNameAlreadyUsed,
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(_ *testing.T) {
				nameGot, err := testAppService.RegisterNewUser(c,
					testCase.UserName,
					testCase.Password,
				)
				t.IsError(testCase.Err, errors.Cause(err), testName)
				t.Is("", nameGot, testName)
			})
		}
	})
}
