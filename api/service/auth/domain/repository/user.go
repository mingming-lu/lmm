package repository

import (
	"context"

	"lmm/api/service/auth/domain/model"
)

// UserRepository interface
type UserRepository interface {
	FindByName(c context.Context, name string) (*model.User, error)
	FindByToken(c context.Context, token *model.Token) (*model.User, error)
}
