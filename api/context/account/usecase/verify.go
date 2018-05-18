package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
	"log"
)

func VerifyToken(encodedToken string) (user *model.User, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v", r)
			err = ErrInvalidToken
		}
	}()
	token, err := service.DecodeToken(encodedToken)
	if err != nil {
		return nil, err
	}

	user, err = repository.New().FindByToken(token)
	return user, err
}
