package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/google/uuid"

	accountFactory "lmm/api/service/account/domain/factory"
	account "lmm/api/service/account/domain/model"
	accountModel "lmm/api/service/account/domain/model"
	accountService "lmm/api/service/account/domain/service"
	accountStorage "lmm/api/service/account/infra"
)

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
