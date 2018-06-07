package factory

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/domain/factory"
	"lmm/api/utils/uuid"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func NewUser(name, password string) (*model.User, error) {
	userID, err := factory.Default().GenerateID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return model.NewUser(userID, name, string(hashedPassword), uuid.New(), uuid.New(), time.Now()), nil
}
