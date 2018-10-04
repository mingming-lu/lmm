package application

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"regexp"

	"github.com/pkg/errors"

	"lmm/api/service/auth/domain"
	"lmm/api/service/auth/domain/model"
	"lmm/api/service/auth/domain/repository"
	"lmm/api/service/auth/domain/service"
)

// Service struct
type Service struct {
	tokenService   service.TokenService
	userRepository repository.UserRepository
}

// NewService creates a new Service pointer
func NewService(
	userRepository repository.UserRepository,
) *Service {
	return &Service{
		tokenService:   service.NewTokenService(),
		userRepository: userRepository,
	}
}

type basicAuth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var (
	patternBasicAuthorization  = regexp.MustCompile(`^Basic (.+)$`)
	patternBearerAuthorization = regexp.MustCompile(`^Bearer (.+)$`)
)

// BasicAuth authorizes given authorization
func (s *Service) BasicAuth(c context.Context, authorization string) (*model.Token, error) {
	matched := patternBasicAuthorization.FindStringSubmatch(authorization)
	if len(matched) != 2 {
		return nil, domain.ErrInvalidBasicAuthFormat
	}

	b, err := base64.URLEncoding.DecodeString(matched[1])
	if err != nil {
		return nil, errors.Wrap(domain.ErrInvalidBasicAuthFormat, err.Error())
	}

	auth := basicAuth{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&auth); err != nil {
		return nil, domain.ErrInvalidBasicAuthFormat
	}

	user, err := s.userRepository.FindByName(c, auth.UserName)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find user named '%s'", auth.UserName)
	}

	if !user.ComparePassword(auth.Password) {
		return nil, domain.ErrPasswordNotMatched
	}

	return s.tokenService.Encode(user.RawToken())
}

// BearerAuth authorized given authorization
func (s *Service) BearerAuth(c context.Context, authorization string) (*model.User, error) {
	matched := patternBearerAuthorization.FindStringSubmatch(authorization)
	if len(matched) != 2 {
		return nil, domain.ErrInvalidBearerAuthFormat
	}

	token, err := s.tokenService.Decode(matched[1])
	if err != nil {
		return nil, errors.Wrap(domain.ErrInvalidBearerAuthFormat, err.Error())
	}

	if token.IsExpired() {
		return nil, domain.ErrTokenExpired
	}

	return s.userRepository.FindByToken(c, token)
}
