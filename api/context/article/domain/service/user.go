package service

import (
	"lmm/api/context/article/domain/model"
)

// UserService is a user adapter interface
type UserService interface {
	ToWriter(userID uint64) model.ArticleWriter
}
