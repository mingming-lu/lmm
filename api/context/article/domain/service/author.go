package service

import (
	"lmm/api/context/article/domain/model"
)

// AuthorService is a user adapter interface
type AuthorService interface {
	AuthorFromUserID(userID uint64) (*model.Author, error)
}
