package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/utils/base64"
	"lmm/api/utils/sha256"
	"lmm/api/utils/token"
)

func SignIn(name, password string) (*model.Response, error) {
	user, err := repository.ByName(name)

	if err != nil {
		return nil, err
	}

	encoded := sha256.Hex([]byte(user.GUID + password))
	if encoded != user.Password {
		return nil, ErrIncorrectPassword
	}

	token, err := token.Encode(user.Token)
	if err != nil {
		return nil, err
	}

	encodedToken := base64.Encode(token)

	res := model.Response{
		ID:    user.ID,
		Name:  user.Name,
		Token: encodedToken,
	}

	return &res, nil
}
