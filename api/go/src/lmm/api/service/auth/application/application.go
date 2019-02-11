package application

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"regexp"

	"github.com/pkg/errors"

	"lmm/api/service/auth/application/command"
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

// Login try to login with basic auth or access token and returns a new access token
func (s *Service) Login(c context.Context, cmd command.LoginCommand) (*model.Token, error) {
	switch cmd.GrantType {
	case domain.GrantTypeBasicAuth:
		return s.BasicAuth(c, cmd.BasicAuth)
	case domain.GrantTypeRefreshToken:
		return s.RefreshToken(c, cmd.AccessToken)
	default:
		return nil, errors.Wrap(domain.ErrInvalidGrantType, cmd.GrantType)
	}
}

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
		return nil, errors.Wrap(domain.ErrNoSuchUser, err.Error())
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

	user, err := s.userRepository.FindByToken(c, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrInvalidAuthToken
		}
		return nil, err
	}

	return user, nil
}

// RefreshToken refreshes access token in bearer auth format
func (s *Service) RefreshToken(c context.Context, auth string) (*model.Token, error) {
	user, err := s.BearerAuth(c, auth)
	if err != nil {
		return nil, err
	}

	return s.tokenService.Encode(user.RawToken())
}
