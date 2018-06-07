package model

import (
	"lmm/api/domain/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	model.Entity
	id        uint64
	name      string
	password  string
	guid      string
	token     string
	createdAt time.Time
}

func NewUser(id uint64, name, password, guid, token string, createdAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		password:  password,
		guid:      guid,
		token:     token,
		createdAt: createdAt,
	}
}

func (u *User) ID() uint64 {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Password() string {
	return u.password
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hashedPassword)
	return nil
}

func (u *User) GUID() string {
	return u.guid
}

func (u *User) Token() string {
	return u.token
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
