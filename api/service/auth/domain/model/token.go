package model

import (
	"lmm/api/domain/model"
	"time"
)

// Token domain model
type Token struct {
	model.ValueObject
	raw    string
	expire time.Time
}

// Raw returns raw token in string
func (t *Token) Raw() string {
	return t.raw
}

// IsExpired reports is token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.expire)
}
