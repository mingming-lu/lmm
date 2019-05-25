package factory

import (
	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"
	"lmm/api/util/uuidutil"
)

type Factory struct {
	encrypter      service.EncryptService
	userRepository repository.UserRepository
}

func NewFactory(encrypter service.EncryptService, userRepository repository.UserRepository) *Factory {
	return &Factory{
		encrypter:      encrypter,
		userRepository: userRepository,
	}
}

func (f *Factory) NewUser(tx transaction.Transaction, username, email, password string) (*model.User, error) {
	hashedPassword, err := f.NewPassword(password)
	if err != nil {
		return nil, err
	}

	token := uuidutil.NewUUID()

	newID, err := f.userRepository.NextID(tx)
	if err != nil {
		return nil, err
	}

	return model.NewUser(newID, username, email, hashedPassword, token, model.Ordinary, clock.Now())
}

func (f *Factory) NewPassword(plainText string) (string, error) {
	pw, err := model.NewPassword(plainText)
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
