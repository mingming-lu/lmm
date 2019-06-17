package util

import (
	"context"
	"lmm/api/service/user/domain/model"
)

type contextKey string

const (
	authKey contextKey = "auth"
)

type Auth struct {
	ID    int64
	Name  string
	Token string
	Role  string
}

func NewContext(c context.Context, auth *Auth) context.Context {
	return context.WithValue(c, authKey, auth)
}

func FromContext(c context.Context) (*Auth, bool) {
	auth, ok := c.Value(authKey).(*Auth)
	if ok {
		return auth, true
	}

	return nil, false
}

func (auth *Auth) IsAdmin() bool {
	return auth.Role == model.Admin.Name()
}
