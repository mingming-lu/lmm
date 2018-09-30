package application

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/auth/domain"
	"lmm/api/service/auth/domain/model"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
)

var (
	userRepo                     *InMemoryUserRepository
	appService                   *Service
	registeredUserName           = "someone"
	registeredUserPassword       = "whatever"
	registeredUserHashedPassword = ""
	registeredUserToken          = stringutil.ReplaceAll(uuid.New().String(), "-", "")
)

type InMemoryUserRepository struct {
	memory map[string]*model.User
}

func (repo *InMemoryUserRepository) FindByName(c context.Context, name string) (*model.User, error) {
	if name == registeredUserName {
		return model.NewUser(registeredUserName, registeredUserHashedPassword, registeredUserToken), nil
	}
	return nil, domain.ErrNoSuchUser
}

func (repo *InMemoryUserRepository) FindByToken(c context.Context, token *model.Token) (*model.User, error) {
	if token.Raw() == registeredUserToken {
		return model.NewUser(registeredUserName, registeredUserHashedPassword, registeredUserToken), nil
	}
	return nil, domain.ErrNoSuchUser
}

func TestMain(m *testing.M) {
	userRepo = &InMemoryUserRepository{memory: make(map[string]*model.User)}
	appService = NewService(userRepo)

	b, _ := bcrypt.GenerateFromPassword([]byte(registeredUserPassword), bcrypt.DefaultCost)
	registeredUserHashedPassword = string(b)

	code := m.Run()
	os.Exit(code)
}

func TestBasicAuth(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	t.Run("Success", func(_ *testing.T) {
		auth := basicAuth{
			UserName: registeredUserName,
			Password: registeredUserPassword,
		}
		b, err := json.Marshal(auth)
		if !t.NoError(err) {
			t.FailNow()
		}

		token, err := appService.BasicAuth(c, "Basic "+base64.URLEncoding.EncodeToString(b))
		t.NoError(err)
		t.NotNil(token)
		t.Is(token.Raw(), registeredUserToken)
	})

	t.Run("InvalidBasicAuthFormat", func(_ *testing.T) {
		token, err := appService.BasicAuth(c, "whatever")
		t.IsError(domain.ErrInvalidBasicAuthFormat, errors.Cause(err))
		t.Nil(token)
	})

	t.Run("Fail", func(_ *testing.T) {
		cases := map[string]struct {
			UserName string
			Password string
			Err      error
		}{
			"NoSuchUser":    {"noone", "whatever", domain.ErrNoSuchUser},
			"WrongPassword": {registeredUserName, "somethingwrong", domain.ErrPasswordNotMatched},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(_ *testing.T) {
				auth := basicAuth{
					UserName: testCase.UserName,
					Password: testCase.Password,
				}
				b, err := json.Marshal(auth)
				if !t.NoError(err) {
					t.FailNow()
				}
				token, err := appService.BasicAuth(c, "Basic "+base64.URLEncoding.EncodeToString(b))
				t.IsError(testCase.Err, errors.Cause(err), testName)
				t.Nil(token, testName)
			})
		}
	})
}

func TestBearerAuth(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	token, err := appService.tokenService.Encode(registeredUserToken)
	if !t.NoError(err) {
		t.FailNow()
	}

	t.Run("Success", func(_ *testing.T) {
		user, err := appService.BearerAuth(c, "Bearer "+token.Hashed())
		t.NoError(err)
		t.Is(token.Raw(), user.RawToken())
	})

	t.Run("InvalidTokenFormat", func(_ *testing.T) {
		user, err := appService.BearerAuth(c, token.Hashed())
		t.IsError(domain.ErrInvalidBearerAuthFormat, errors.Cause(err))
		t.Nil(user)
	})

	t.Run("NoSuchUser", func(_ *testing.T) {
		token, err := appService.tokenService.Encode(uuid.New().String())
		if !t.NoError(err) {
			t.FailNow()
		}
		user, err := appService.BearerAuth(c, "Bearer "+token.Hashed())
		t.IsError(domain.ErrNoSuchUser, errors.Cause(err))
		t.Nil(user)
	})
}
