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

// AuthorFromUserName implements AuthorAdapter.AuthorFromUserID
func (a *AuthorAdapter) AuthorFromUserName(c context.Context, userName string) (*model.Author, error) {
	stmt := a.db.Prepare(c, `select name from user where name = ?`)
	defer stmt.Close()

	var authorName string
	if err := stmt.QueryRow(c, userName).Scan(&authorName); err != nil {
		if err == storage.ErrNoRows {
			return nil, domain.ErrNoSuchUser
		}
		return nil, err
	}
	return model.NewAuthor(authorName), nil
}
