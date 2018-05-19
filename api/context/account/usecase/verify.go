package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
)

func (uc *Usecase) VerifyToken(encodedToken string) (user *model.User, err error) {
	token, err := service.DecodeToken(encodedToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err = uc.repo.FindByToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return user, nil
}
