package testutil

import (
	"context"
	"time"

	"lmm/api/clock"
	"lmm/api/service/user/domain/model"
	"lmm/api/storage/db"
	"lmm/api/util/stringutil"

	"github.com/google/uuid"
)

// User model for testing
type User struct {
	id                int64
	name              string
	rawPassword       string
	encryptedPassword string
	rawToken          string
	accessToken       string
	createdAt         time.Time
}

// NewUser creates a user for testing
func NewUser(db db.DB) User {
	username := "U" + uuid.New().String()[:8]
	rawPassword := uuid.New().String()

	password, err := model.NewPassword(rawPassword)
	if err != nil {
		panic(err)
	}

	user, err := model.NewUser(
		username,
		*password,
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
	)
	if err != nil {
		panic(err)
	}

	now := clock.Now()

	res, err := db.Exec(context.Background(),
		`insert into user (name, password, token, created_at) values (?, ?, ?, ?)`,
		user.Name(), user.Password(), user.Token(), now,
	)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return User{
		id:                id,
		name:              username,
		rawPassword:       rawPassword,
		encryptedPassword: user.Password(),
		rawToken:          user.Token(),
		accessToken:       EncodeToken(user.Token()).Hashed(),
		createdAt:         now,
	}
}

func (u User) ID() int64 {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) RawPassword() string {
	return u.rawPassword
}

func (u User) EncryptedPassword() string {
	return u.encryptedPassword
}

func (u User) RawToken() string {
	return u.rawToken
}

func (u User) AccessToken() string {
	return u.accessToken
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}
