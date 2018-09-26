package factory

import (
	"errors"
	"time"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/context/account/domain/model"
	"lmm/api/utils/uuid"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

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

	userID, err := idGenerator.NextID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return model.NewUser(userID, name, string(hashedPassword), uuid.New(), time.Now()), nil
}
