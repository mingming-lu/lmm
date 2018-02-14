package user

import (
	model "lmm/api/domain/model/user"
	repo "lmm/api/domain/repository/user"
	"lmm/api/domain/service/base64"
	"lmm/api/domain/service/token"
	"lmm/api/domain/service/uuid"

	"github.com/akinaru-lu/errors"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
)

func SignUp(name, password string) (int64, error) {
	token := uuid.New()
	guid := uuid.New()
	password = base64.Encode([]byte(guid + password))
	return repo.Add(name, password, guid, token)
}

func SignIn(name, password string) (*model.Response, error) {
	user, err := repo.ByName(name)

	if err != nil {
		return nil, err
	}

	encoded := base64.Encode([]byte(user.GUID + password))
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

func Verify(encodedToken string) (*model.User, error) {
	encoded, err := base64.Decode(encodedToken)
	if err != nil {
		return nil, err
	}

	token, err := token.Decode(encoded)
	if err != nil {
		return nil, err
	}

	return repo.ByToken(token)
}
