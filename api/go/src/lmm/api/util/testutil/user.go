package testutil

import (
	"context"
	"time"

	"lmm/api/clock"
	"lmm/api/service/user/domain/model"
	"lmm/api/storage/db"
	"lmm/api/util/stringutil"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User model for testing
type User struct {
	id                int64
	name              string
	email             string
	rawPassword       string
	encryptedPassword string
	rawToken          string
	accessToken       string
	role              model.Role
	createdAt         time.Time
}

// NewAdmin creates a new admin user for testing
func NewAdmin(db db.DB) User {
	user := NewUser(db)
	_, err := db.Exec(context.TODO(), `update user set role = ? where name = ?`, "admin", user.Name())
	if err != nil {
		panic(err)
	}
	return user
}

// NewUser creates a user for testing
func NewUser(db db.DB) User {
	username := "U" + uuid.New().String()[:8]
	email := username + "@lmm.local"
	rawPassword := uuid.New().String()

	password, err := model.NewPassword(rawPassword)
	if err != nil {
		panic(err)
	}

	b, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user, err := model.NewUser(
		0, // auto-generate allocate new id by datastore
		username,
		email,
		string(b),
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
		model.Ordinary,
		clock.Now(),
	)
	if err != nil {
		panic(err)
	}

	res, err := db.Exec(context.TODO(),
		`insert into user (name, email, password, token, role, created_at) values (?, ?, ?, ?, ?, ?)`,
		user.Name(), user.Email(), user.Password(), user.Token(), user.Role().Name(), user.RegisteredAt(),
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
		email:             email,
		rawPassword:       rawPassword,
		encryptedPassword: user.Password(),
		rawToken:          user.Token(),
		accessToken:       EncodeToken(user.Token()).Hashed(),
		role:              user.Role(),
		createdAt:         user.RegisteredAt(),
	}
}

func (u User) ID() int64 {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) Email() string {
	return u.email
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

func (u User) Role() model.Role {
	return u.role
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}
