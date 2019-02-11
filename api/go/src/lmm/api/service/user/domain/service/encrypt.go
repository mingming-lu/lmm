package service

import (
	"lmm/api/service/user/domain/model"

	"golang.org/x/crypto/bcrypt"
)

type EncryptService interface {
	Encrypt(password *model.Password) (encryptedText string, err error)
	Verify(raw, hashed string) bool
}

type BcryptService struct{}

func (s *BcryptService) Encrypt(password *model.Password) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *BcryptService) Verify(raw, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)); err != nil {
		return false
	}
	return true
}
