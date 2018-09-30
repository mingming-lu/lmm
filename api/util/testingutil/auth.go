package testingutil

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/auth/domain/model"
	"lmm/api/storage/db"
	"lmm/api/util/stringutil"
)

// NewAuthUser creates new user from auth service
func NewAuthUser(db db.DB) (*model.User, error) {
	rawPassword := uuid.New().String()
	b, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	encryptedPassword := string(b)

	user := model.NewUser(
		uuid.New().String()[:8],
		encryptedPassword,
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
	)

	now := time.Now()

	if _, err := db.Exec(context.Background(), `
		insert into user (name, password, token, created_at) values (?, ?, ?, ?)
	`, user.Name(), encryptedPassword, user.RawToken(), now); err != nil {
		panic(err)
	}

	return user, nil
}
