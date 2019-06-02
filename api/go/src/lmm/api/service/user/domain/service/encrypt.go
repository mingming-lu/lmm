package service

import (
	"lmm/api/service/user/domain/model"
)

type EncryptService interface {
	Encrypt(password *model.Password) (encryptedText string, err error)
	Verify(raw, hashed string) bool
}
