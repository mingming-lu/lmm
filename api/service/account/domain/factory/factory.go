package factory

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/account/domain/model"
	"lmm/api/strings"
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

	token := strings.ReplaceAll(uuid.New().String(), "-", "")

	return model.NewUser(userID, name, string(hashedPassword), token, time.Now()), nil
}
