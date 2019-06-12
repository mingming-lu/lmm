package model

import (
	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/service/user/domain"
	"lmm/api/util/uuidutil"
)

type Factory struct {
	encrypter      EncryptService
	userRepository UserRepository
}

func NewFactory(encrypter EncryptService, userRepository UserRepository) *Factory {
	return &Factory{
		encrypter:      encrypter,
		userRepository: userRepository,
	}
}

func (f *Factory) NewUser(tx transaction.Transaction, username, email, password string) (*User, error) {
	hashedPassword, err := f.NewPassword(password)
	if err != nil {
		return nil, err
	}

	token := f.NewToken()

	newID, err := f.userRepository.NextID(tx)
	if err != nil {
		return nil, err
	}

	return NewUser(newID, username, email, hashedPassword, token, Ordinary, clock.Now())
}

func (f *Factory) NewPassword(plainText string) (string, error) {
	pw, err := NewPassword(plainText)
	if err != nil {
		return "", err
	}

	if pw.IsWeak() {
		return "", domain.ErrUserPasswordTooWeak
	}

	hashedPassword, err := f.encrypter.Encrypt(pw)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func (f *Factory) NewToken() string {
	return uuidutil.NewUUID()
}
