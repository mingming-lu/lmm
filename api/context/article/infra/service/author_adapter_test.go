package service

import (
	"os"

	"github.com/google/uuid"

	accountFactory "lmm/api/context/account/domain/factory"
	accountRepository "lmm/api/context/account/domain/repository"
	accountStorage "lmm/api/context/account/infra"
	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/service"
	"lmm/api/storage"
	"lmm/api/testing"
)

var (
	authorAdapter  service.AuthorService
	userRepository accountRepository.UserRepository
)

func TestMain(m *testing.M) {
	db := storage.NewDB()
	defer db.Close()

	authorAdapter = NewAuthorAdapter(db)
	userRepository = accountStorage.NewUserStorage(db)
	code := m.Run()
	os.Exit(code)
}

func TestAuthorFromUserID_OK(tt *testing.T) {
	t := testing.NewTester(tt)

	name, password := uuid.New().String()[:5], uuid.New().String()
	user, err := accountFactory.NewUser(name, password)
	t.NoError(err)
	t.NoError(userRepository.Add(user))

	author, err := authorAdapter.AuthorFromUserID(user.ID())
	t.NoError(err)
	t.Is(user.Name(), author.Name())
}

func TestAuthorFromUserID_NoSuchUser(tt *testing.T) {
	t := testing.NewTester(tt)

	author, err := authorAdapter.AuthorFromUserID(testing.GenerateID())
	t.IsError(domain.ErrNoSuchUser, err)
	t.Nil(author)
}
