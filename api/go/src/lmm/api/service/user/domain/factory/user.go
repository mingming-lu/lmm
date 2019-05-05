package factory

import (
	"lmm/api/clock"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/service"
	"lmm/api/util/uuidutil"
)

type Factory struct {
	encrypter service.EncryptService
}

func NewFactory(encrypter service.EncryptService) *Factory {
	return &Factory{
		encrypter: encrypter,
	}
}

func (f *Factory) NewUser(username, email, password string) (*model.User, error) {
	hashedPassword, err := f.NewPassword(password)
	if err != nil {
		return nil, err
	}

	token := f.NewToken()

	return model.NewUser(username, email, hashedPassword, token, model.Ordinary, clock.Now())
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

func (f *Factory) NewToken() string {
	return uuidutil.NewUUID()
}
