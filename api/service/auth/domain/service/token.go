package service

import "lmm/api/service/auth/domain/model"

// TokenService provides interfaces to encode/decode token
type TokenService interface {
	Encode(rawToken string) (*model.Token, error)
	Decode(hashedToken string) (*model.Token, error)
}
