package application

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"lmm/api/clock"
	"lmm/api/pkg/pubsub/pubsubtest"
	testUtil "lmm/api/pkg/testing"
	"lmm/api/pkg/transaction"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/port/adapter/messaging"
	"lmm/api/service/user/port/adapter/service"
	"lmm/api/util/uuidutil"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	testAppService *Service
)

type InmemoryUserRepository struct {
	sync.RWMutex
	memory map[model.UserID]*model.User
}

func (repo *InmemoryUserRepository) NextID(tx transaction.Transaction) (model.UserID, error) {
	repo.Lock()
	defer repo.Unlock()

	return model.UserID(int64(len(repo.memory) + 1)), nil
}

func (repo *InmemoryUserRepository) Save(tx transaction.Transaction, user *model.User) error {
	repo.Lock()
	defer repo.Unlock()

	repo.memory[user.ID()] = user
	return nil
}

func (repo *InmemoryUserRepository) FindByName(tx transaction.Transaction, username string) (*model.User, error) {
	repo.RLock()
	defer repo.RUnlock()

	for _, user := range repo.memory {
		if user.Name() == username {
			return user, nil
		}
	}
	return nil, domain.ErrNoSuchUser
}

func (repo *InmemoryUserRepository) FindByToken(tx transaction.Transaction, token string) (*model.User, error) {
	repo.RLock()
	defer repo.RUnlock()

	for _, user := range repo.memory {
		if user.Token() == token {
			return user, nil
		}
	}
	return nil, domain.ErrNoSuchUser
}

func (repo *InmemoryUserRepository) Begin(c context.Context, opts *transaction.Option) (transaction.Transaction, error) {
	return transaction.Nop(), nil
}

func (repo *InmemoryUserRepository) RunInTransaction(c context.Context, f func(tx transaction.Transaction) error, opts *transaction.Option) error {
	tx, err := repo.Begin(c, opts)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}
	defer tx.Commit()

	return f(tx)
}

func TestMain(m *testing.M) {
	repo := &InmemoryUserRepository{memory: make(map[model.UserID]*model.User)}
	pubsubClient := pubsubtest.NewClient()
	pub := messaging.NewUserEventPublisher(pubsubClient)
	testAppService = NewService(
		&service.BcryptService{},
		testUtil.TokenService,
		repo, repo, pub)
	code := m.Run()
	pubsubClient.Close()
	os.Exit(code)
}

func TestRegisterNewUser(t *testing.T) {
	c := context.Background()

	t.Run("Success", func(t *testing.T) {
		username, password := "username", "~!@#$%^&*()-_=+{[}]|\\:;\"'<,>.?/"
		userID, err := testAppService.RegisterNewUser(c, command.Register{
			UserName:     username,
			EmailAddress: username + "@lmm.local",
			Password:     password,
		})
		assert.NoError(t, err)
		assert.NotZero(t, userID)

		user, err := testAppService.userRepository.FindByName(nil, "username")
		assert.NoError(t, err)
		assert.Equal(t, username, user.Name())
		assert.NoError(t, bcrypt.CompareHashAndPassword(
			[]byte(user.Password()),
			[]byte(password),
		))
		assert.NotPanics(t, func() {
			uuid.Must(uuidutil.ParseString(user.Token()))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		cases := map[string]struct {
			UserName string
			Email    string
			Password string
			Err      error
		}{
			"UserNameTooShort": {
				"ur", "ur@lmm.local", "password1234", domain.ErrInvalidUserName,
			},
			"UserNameStartsWithoutLetter": {
				"1username", "1username@lmm.local", "password1234", domain.ErrInvalidUserName,
			},
			"EmptyPassword": {
				"username", "username@lmm.local", "", domain.ErrUserPasswordEmpty,
			},
			"PasswordIsTooShort": {
				"username", "username@lmm.local", "passwor", domain.ErrUserPasswordTooShort,
			},
			"PasswordIsTooWeak": {
				"username", "username@lmm.local", "password", domain.ErrUserPasswordTooWeak,
			},
			"PasswordIsTooLong": {
				"username", "username@lmm.local", strings.Repeat("s", 251), domain.ErrUserPasswordTooLong,
			},
			"DuplicateUserName": {
				"username", "username@lmm.local", "password1234", domain.ErrUserNameAlreadyUsed,
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(t *testing.T) {
				userIDGot, err := testAppService.RegisterNewUser(c, command.Register{
					UserName:     testCase.UserName,
					EmailAddress: testCase.Email,
					Password:     testCase.Password,
				})
				assert.Error(t, testCase.Err, errors.Cause(err), testName)
				assert.Equal(t, int64(0), userIDGot, testName)
			})
		}
	})
}

func TestUserChangePassword(t *testing.T) {
	c := context.Background()

	username, password := "U"+uuidutil.NewUUID()[:8], "U$ErP@ssw0rD"
	userID, err := testAppService.RegisterNewUser(c, command.Register{
		UserName:     username,
		EmailAddress: username + "@lmm.local",
		Password:     password,
	})
	if !assert.NoError(t, err) || !assert.NotZero(t, userID) {
		t.Fatal("failed to create new user")
	}

	userBeforePasswordChanging, err := testAppService.userRepository.FindByName(nil, username)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	// record value since it's changed by pointer
	oldToken := userBeforePasswordChanging.Token()

	newPassword := uuidutil.NewUUID() + uuidutil.NewUUID()

	assert.NoError(t, testAppService.UserChangePassword(c, command.ChangePassword{
		User:        username,
		OldPassword: password,
		NewPassword: newPassword,
	}))

	userAfterPasswordChanging, err := testAppService.userRepository.FindByName(nil, username)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	assert.True(t, testAppService.encrypter.Verify(newPassword, userAfterPasswordChanging.Password()))
	assert.NotEqual(t, oldToken, userAfterPasswordChanging.Token())
}

func newAdmin() *model.User {
	return newUserWithRole(model.Admin)
}

func newOrdinary() *model.User {
	return newUserWithRole(model.Ordinary)
}

func newUserWithRole(role model.Role) *model.User {
	randomUserName := "u" + uuid.New().String()[:7]
	email := randomUserName + "@lmm.local"
	password := uuid.New().String()
	token := uuid.New().String()

	user, err := model.NewUser(0, randomUserName, email, password, token, role, clock.Now())
	if err != nil {
		panic(err)
	}

	return user
}
