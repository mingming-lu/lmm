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
	"lmm/api/service/user/domain/repository"
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
	if user, _ := repo.FindByName(c, user.Name()); user != nil {
		return domain.ErrUserNameAlreadyUsed
	}
	repo.memory = append(repo.memory, user)
	return nil
}

func (repo *InmemoryUserRepository) FindByName(c context.Context, username string) (*model.User, error) {
	for _, user := range repo.memory {
		if user.Name() == username {
			return user, nil
		}
	}
	return nil, domain.ErrNoSuchUser
}

func (repo *InmemoryUserRepository) DescribeByName(c context.Context, username string) (*model.UserDescriptor, error) {
	panic("not implemented")
}

func (repo *InmemoryUserRepository) DescribeAll(context.Context, repository.DescribeAllOptions) ([]*model.UserDescriptor, uint, error) {
	panic("not implemented")
}

func TestMain(m *testing.M) {
	testAppService = NewService(&InmemoryUserRepository{memory: make([]*model.User, 0)})
	code := m.Run()
	os.Exit(code)
}

func TestRegisterNewUser(tt *testing.T) {
	c := context.Background()

	tt.Run("Success", func(tt *testing.T) {
		t := testing.NewTester(tt)
		username, password := "username", "~!@#$%^&*()-_=+{[}]|\\:;\"'<,>.?/"
		nameGot, err := testAppService.RegisterNewUser(c, username, password)
		t.NoError(err)
		t.Is(username, nameGot)

		user, err := testAppService.userRepository.FindByName(c, "username")
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

	tt.Run("Fail", func(tt *testing.T) {
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
			tt.Run(testName, func(tt *testing.T) {
				t := testing.NewTester(tt)
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
