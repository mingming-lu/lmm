package service

import (
	"context"

	"github.com/pkg/errors"

	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/model"
	"lmm/api/service/asset/domain/service"
	"lmm/api/storage/db"
)

// UserAdapter implements UploaderService
type UserAdapter struct {
	db db.DB
}

// NewUserAdapter creates a UploaderService implementation
func NewUserAdapter(db db.DB) service.UploaderService {
	return &UserAdapter{db: db}
}

// FromUserName implementation
func (adapter *UserAdapter) FromUserName(c context.Context, name string) (*model.Uploader, error) {
	stmt := adapter.db.Prepare(c, `select id from user where name = ?`)
	defer stmt.Close()

	var (
		userID int64
	)
	if err := stmt.QueryRow(c, name).Scan(&userID); err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchUser, err.Error())
	}

	return model.NewUploader(userID), nil
}
