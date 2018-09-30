package model

import (
	"lmm/api/domain/model"
	"time"
)

// Token domain model
type Token struct {
	model.ValueObject
	raw    string
	hashed string
	expire time.Time
}

// NewToken creates a new token
func NewToken(raw, hashed string, expire time.Time) *Token {
	return &Token{raw: raw, hashed: hashed, expire: expire}
}

// Raw returns raw token string
func (t *Token) Raw() string {
	return t.raw
}

// Hashed returns hashed token in string
func (t *Token) Hashed() string {
	return t.hashed
}

// IsExpired reports is token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.expire)
}
