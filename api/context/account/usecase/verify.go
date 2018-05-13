package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/utils/base64"
	"lmm/api/utils/token"
)

func Verify(encodedToken string) (*model.User, error) {
	encoded, err := base64.Decode(encodedToken)
	if err != nil {
		return nil, err
	}

	token, err := token.Decode(encoded)
	if err != nil {
		return nil, err
	}

	return repository.ByToken(token)
}
