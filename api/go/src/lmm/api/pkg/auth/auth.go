package auth

import (
	"context"

	"lmm/api/service/user/port/adapter/util"
)

type Auth = util.Auth

func NewContext(c context.Context, auth *Auth) context.Context {
	return util.NewContext(c, auth)
}

func FromContext(c context.Context) (*Auth, bool) {
	return util.FromContext(c)
}
