package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/application/query"
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

const maxCount uint = 100

func (s *Service) ViewAllUsersByOptions(c context.Context, query query.ViewAllUsers) ([]*model.UserDescriptor, error) {
	page, err := stringutil.ParseUint(query.Page)
	if err != nil || page == 0 {
		return nil, errors.Wrap(domain.ErrInvalidPage, query.Page)
	}

	count, err := stringutil.ParseUint(query.Count)
	if err != nil || count > maxCount {
		return nil, errors.Wrap(domain.ErrInvalidCount, query.Count)
	}

	order, err := s.mappingOrder(query.OrderBy, query.Order)
	if err != nil {
		return nil, errors.Wrap(domain.ErrInvalidViewOrder, query.Order)
	}
	return s.userRepository.DescribeAll(c, repository.DescribeAllOptions{
		Page:  page,
		Count: count,
		Order: order,
	})
}

func (s *Service) mappingOrder(orderBy, order string) (repository.DescribeAllOrder, error) {
	switch orderBy + "_" + order {
	case "name_asc":
		return repository.DescribeAllOrderByNameAsc, nil
	case "name_desc":
		return repository.DescribeAllOrderByNameDesc, nil
	case "registered_date_asc":
		return repository.DescribeAllOrderByRegisteredDateAsc, nil
	case "registered_date_desc":
		return repository.DescribeAllOrderByRegisteredDateDesc, nil
	case "role_asc":
		return repository.DescribeAllOrderByRoleAsc, nil
	case "role_desc":
		return repository.DescribeAllOrderByRoleDesc, nil
	default:
		return repository.DescribeAllOrder(-1), domain.ErrInvalidViewOrder
	}
}
