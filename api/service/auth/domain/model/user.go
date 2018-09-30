package model

import (
	"golang.org/x/crypto/bcrypt"

	"lmm/api/domain/model"
)

// User domain model
type User struct {
	model.Entity
	name     string
	password string
	rawToken string
}

// Name gets user name
func (user *User) Name() string {
	return user.name
}

// ComparePassword compares encrypted password with given raw password
// and returns true if matched
func (user *User) ComparePassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.password),
		[]byte(rawPassword),
	)
	if err == bcrypt.ErrHashTooShort {
		panic(bcrypt.ErrHashTooShort.Error())
	}
	return bcrypt.ErrMismatchedHashAndPassword != err
}

// RawToken gets user's raw token
func (user *User) RawToken() string {
	return user.rawToken
}
