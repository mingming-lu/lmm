package testingutil

import (
	"context"
	"time"

	"lmm/api/service/user/domain/model"
	"lmm/api/storage/db"
	"lmm/api/util/stringutil"

	"github.com/google/uuid"
)

// NewUserUser create new user from user service
func NewUserUser(db db.DB, username, rawPassword string) (*model.User, error) {
	password, err := model.NewPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	user, err := model.NewUser(
		username,
		*password,
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if _, err := db.Exec(context.Background(), `
		insert into user (name, password, token, created_at) values (?, ?, ?, ?)
	`, user.Name(), user.Password(), user.Token(), now); err != nil {
		return nil, err
	}

	return user, nil
}
