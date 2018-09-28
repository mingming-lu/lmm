package service

import (
	"context"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/storage"
	"lmm/api/storage/db"
)

// AuthorAdapter is an implementation of AuthorService
type AuthorAdapter struct {
	db db.DB
}

// NewAuthorAdapter is a construct of AuthorAdapter
func NewAuthorAdapter(db db.DB) *AuthorAdapter {
	return &AuthorAdapter{db: db}
}

// AuthorFromUserID implements AuthorAdapter.AuthorFromUserID
func (a *AuthorAdapter) AuthorFromUserID(c context.Context, userID uint64) (*model.Author, error) {
	stmt := a.db.Prepare(c, `SELECT id, name FROM user WHERE id = ?`)
	defer stmt.Close()

	var (
		authorID   uint64
		authorName string
	)
	if err := stmt.QueryRow(userID).Scan(&authorID, &authorName); err != nil {
		if err == storage.ErrNoRows {
			return nil, domain.ErrNoSuchUser
		}
		return nil, err
	}
	return model.NewAuthor(int64(authorID), authorName), nil
}
