package service

import (
	"lmm/api/service/user/domain/model"

	"golang.org/x/crypto/bcrypt"
)

// BcryptService implements service.EncryptService
type BcryptService struct{}

// Encrypt encrypts password into hashed one
func (s *BcryptService) Encrypt(password *model.Password) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Verify compares raw and hashed by bcrypt algorithm
func (s *BcryptService) Verify(raw, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)); err != nil {
		return false
	}
	return true
}
