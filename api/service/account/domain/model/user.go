package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"lmm/api/domain/model"
)

type User struct {
	model.Entity
	id        uint64
	name      string
	password  string
	token     string
	createdAt time.Time
}

func NewUser(id uint64, name, password, token string, createdAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		password:  password,
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

func (u *User) VerifyPassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.password), []byte(rawPassword))
}

func (u *User) UpdatePassword(rawPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hashedPassword)
	return nil
}

func (u *User) Token() string {
	return u.token
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
