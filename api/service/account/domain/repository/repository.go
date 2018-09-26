package repository

import (
	"lmm/api/service/account/domain/model"
)

type UserRepository interface {
	Add(*model.User) error
	FindByName(string) (*model.User, error)
	FindByToken(string) (*model.User, error)
}
