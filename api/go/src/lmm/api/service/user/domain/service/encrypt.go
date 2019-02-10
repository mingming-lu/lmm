package service

import (
	"lmm/api/service/user/domain/model"

	"golang.org/x/crypto/bcrypt"
)

type EncryptService interface {
	Encrypt(password *model.Password) (encryptedText string, err error)
}

type BcryptService struct{}

func (s *BcryptService) Encrypt(password *model.Password) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
