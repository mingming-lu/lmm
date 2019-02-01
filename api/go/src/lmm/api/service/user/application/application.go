package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"
	"lmm/api/util/stringutil"
)

// Service is a application service
type Service struct {
	userRepository repository.UserRepository
}

// NewService creates a new Service pointer
func NewService(userRepository repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

// RegisterNewUser registers new user
func (s *Service) RegisterNewUser(c context.Context, name, password string) (string, error) {
	token := uuid.New().String()
	token = stringutil.ReplaceAll(token, "-", "")

	pw, err := model.NewPassword(password)
	if err != nil {
		return "", errors.Wrap(err, "invalid input password")
	}

	if pw.IsWeak() {
		return "", domain.ErrUserPasswordTooWeak
	}

	user, err := model.NewUser(name, *pw, token, model.Ordinary)
	if err != nil {
		return "", errors.Wrap(err, "failed to register new user")
	}

	if err := s.userRepository.Save(c, user); err != nil {
		return "", errors.Wrap(err, "failed to save user")
	}

	return name, nil
}

// AssignRole handles command which operator assign user to role
func (s *Service) AssignRole(c context.Context, cmd command.AssignRole) error {
	operator, err := s.userRepository.DescribeByName(c, cmd.OperatorUser)
	if err != nil {
		http.Log().Panic(c, errors.Wrapf(err, "operator not found: %s", cmd.OperatorUser).Error())
	}

	user, err := s.userRepository.DescribeByName(c, cmd.TargetUser)
	if err != nil {
		return errors.Wrap(domain.ErrNoSuchUser, cmd.TargetUser)
	}

	role := service.RoleAdapter(cmd.TargetRole)

	return service.AssignUserRole(c, operator, user, role)
}
