package service

import (
	"context"

	"lmm/api/service/article/domain/model"
)

// AuthorService is a user adapter interface
type AuthorService interface {
	AuthorFromUserID(c context.Context, userID uint64) (*model.Author, error)
}
