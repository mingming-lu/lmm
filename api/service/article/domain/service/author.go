package service

import (
	"context"

	"lmm/api/service/article/domain/model"
)

// AuthorService is a user adapter interface
type AuthorService interface {
	AuthorFromUserName(c context.Context, userName string) (*model.Author, error)
}
