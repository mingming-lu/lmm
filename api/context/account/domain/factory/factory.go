package factory

import (
	"errors"
	"lmm/api/context/account/domain/model"
	"lmm/api/domain/factory"
	"lmm/api/utils/uuid"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyUserName = errors.New("empty user name")
	ErrEmptyPassword = errors.New("empry password")
)

func NewUser(name, password string) (*model.User, error) {
	if name == "" {
		return nil, ErrEmptyUserName
	}

	if password == "" {
		return nil, ErrEmptyPassword
	}

	userID, err := factory.Default().GenerateID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return model.NewUser(userID, name, string(hashedPassword), uuid.New(), time.Now()), nil
}
