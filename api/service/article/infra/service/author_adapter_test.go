package service

import (
	"context"
	"math/rand"
	"os"

	"github.com/google/uuid"

	accountFactory "lmm/api/service/account/domain/factory"
	accountRepository "lmm/api/service/account/domain/repository"
	accountStorage "lmm/api/service/account/infra"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/service"
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

	author, err := authorAdapter.AuthorFromUserID(context.Background(), user.ID())
	t.NoError(err)
	t.Is(user.Name(), author.Name())
}

func TestAuthorFromUserID_NoSuchUser(tt *testing.T) {
	t := testing.NewTester(tt)

	author, err := authorAdapter.AuthorFromUserID(context.Background(), uint64(rand.Intn(10000)))
	t.IsError(domain.ErrNoSuchUser, err)
	t.Nil(author)
}
