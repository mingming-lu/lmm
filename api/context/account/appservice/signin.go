package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
	"lmm/api/db"
	"lmm/api/utils/sha256"
)

// SignIn is a usecase which users sign in with a account
func (uc *Usecase) SignIn(name, password string) (*model.User, error) {
	if name == "" || password == "" {
		return nil, ErrEmptyUserNameOrPassword
	}

	user, err := uc.repo.FindByName(name)
	if err != nil {
		if err.Error() == db.ErrNoRows.Error() {
			return nil, ErrInvalidUserNameOrPassword
		}
		return nil, err
	}

	// validate password
	encoded := sha256.Hex([]byte(user.GUID + password))
	if encoded != user.Password {
		return nil, ErrInvalidUserNameOrPassword
	}

	encodedToken := service.EncodeToken(user.Token)
	user.Token = encodedToken

	return user, nil
}
