package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/google/uuid"

	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	accountModel "lmm/api/context/account/domain/model"
	accountService "lmm/api/context/account/domain/service"
	accountStorage "lmm/api/context/account/infra"
	"lmm/api/domain/factory"
)

func GenerateID() uint64 {
	id, err := factory.Default().GenerateID()
	if err != nil {
		panic(err)
	}
	return id
}

func StructToRequestBody(o interface{}) io.ReadCloser {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}

func NewUser() *account.User {
	var err error

	name, password := uuid.New().String()[:10], uuid.New().String()
	user, err := accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	err = accountStorage.NewUserStorage(DB()).Add(user)
	if err != nil {
		panic(err)
	}

	return accountModel.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		accountService.EncodeToken(user.Token()),
		user.CreatedAt(),
	)
}
