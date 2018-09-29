package repository

import (
	"context"

	"lmm/api/service/user/domain/model"
)

// UserRepository interface
type UserRepository interface {
	Save(c context.Context, user *model.User) error
	FindByUserName(c context.Context, username string) (*model.User, error)
}
