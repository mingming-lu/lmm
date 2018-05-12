package user

import (
	"encoding/json"
	"io"
	model "lmm/api/domain/model/user"
	repo "lmm/api/domain/repository/user"
	"lmm/api/domain/service/base64"
	"lmm/api/domain/service/sha256"
	"lmm/api/domain/service/token"

	"github.com/akinaru-lu/errors"
)

var (
	ErrInvalidInput      = errors.New("Invalid input")
	ErrDuplicateUserName = errors.New("User name duplicated")
	ErrIncorrectPassword = errors.New("Incorrect password")
)

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Usecase struct {
	repo *repo.Repository
}

func New(repo *repo.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) SignUp(requestBody io.ReadCloser) (uint64, error) {
	auth := &Auth{}
	err := json.NewDecoder(requestBody).Decode(auth)
	if err != nil {
		return 0, ErrInvalidInput
	}
	m := model.New(auth.Name, auth.Password)
	user, err := uc.repo.Save(m)
	if err != nil {
		key, _, ok := uc.repo.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "name" {
			return 0, ErrDuplicateUserName
		}
		return 0, err
	}
	return user.ID, nil
}

func SignIn(name, password string) (*model.Response, error) {
	user, err := repo.ByName(name)

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
