package model

import "time"

type AccessToken struct {
	raw    string
	hashed string
	expire time.Time
}

func (token AccessToken) Raw() string {
	return token.raw
}

func (token AccessToken) Hashed() string {
	return token.hashed
}

func (token AccessToken) Expired() bool {
	return token.expire.Before(time.Now())
}

func NewAccessToken(raw, hashed string, expire time.Time) *AccessToken {
	return &AccessToken{
		raw:    raw,
		hashed: hashed,
		expire: expire,
	}
}

type TokenService interface {
	Encrypt(string) (*AccessToken, error)
	Decrypt(string) (*AccessToken, error)
}
