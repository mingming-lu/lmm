package application

import (
	"context"

	"lmm/api/service/user/application/command"
	"lmm/api/service/user/application/query"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"
	"lmm/api/transaction"
	"lmm/api/util/stringutil"

	"github.com/pkg/errors"
)

// Service is a application service
type Service struct {
	encrypter      service.EncryptService
	factory        *factory.Factory
	userRepository repository.UserRepository
	txManager      transaction.Manager
}

// NewService creates a new Service pointer
func NewService(txManager transaction.Manager, userRepository repository.UserRepository) *Service {
	encrypter := &service.BcryptService{}
	return &Service{
		encrypter:      encrypter,
		factory:        factory.NewFactory(encrypter),
		userRepository: userRepository,
		txManager:      txManager,
	}
}

// RegisterNewUser registers new user
func (s *Service) RegisterNewUser(c context.Context, cmd command.Register) error {
	return s.txManager.RunInTransaction(c, func(c context.Context) error {
		user, err := s.factory.NewUser(cmd.UserName, cmd.EmailAddress, cmd.Password)
		if err != nil {
			return err
		}

		if _, err := s.userRepository.FindByName(c, cmd.UserName); err == nil {
			return domain.ErrUserNameAlreadyUsed
		}

		if err := s.userRepository.Save(c, user); err != nil {
			return err
		}

		return nil
	})
}

// AssignRole handles command which operator assign user to role
func (s *Service) AssignRole(c context.Context, cmd command.AssignRole) error {
	return s.txManager.RunInTransaction(c, func(c context.Context) error {
		operator, err := s.userRepository.FindByName(c, cmd.OperatorUser)
		if err != nil {
			return errors.Wrap(domain.ErrNoSuchUser, err.Error())
		}

		user, err := s.userRepository.FindByName(c, cmd.TargetUser)
		if err != nil {
			return errors.Wrap(domain.ErrNoSuchUser, err.Error())
		}

		role := model.NewRole(cmd.TargetRole)
		if role == model.Guest {
			return domain.ErrNoSuchRole
		}

		if err := operator.AssignRole(user, role); err != nil {
			return err
		}

		return s.userRepository.Save(c, user)
	})
}

const maxCount uint = 100

func (s *Service) ViewAllUsersByOptions(c context.Context, query query.ViewAllUsers) ([]*model.UserDescriptor, uint, error) {
	page, err := stringutil.ParseUint(query.Page)
	if err != nil || page == 0 {
		return nil, 0, errors.Wrap(domain.ErrInvalidPage, query.Page)
	}

	count, err := stringutil.ParseUint(query.Count)
	if err != nil || count > maxCount {
		return nil, 0, errors.Wrap(domain.ErrInvalidCount, query.Count)
	}

	order, err := s.mappingOrder(query.OrderBy, query.Order)
	if err != nil {
		return nil, 0, errors.Wrap(domain.ErrInvalidViewOrder, query.Order)
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

// UserChangePassword supports a application to chagne user's password
func (s *Service) UserChangePassword(c context.Context, cmd command.ChangePassword) error {
	hashedPassword, err := s.factory.NewPassword(cmd.NewPassword)
	if err != nil {
		return err
	}

	return s.txManager.RunInTransaction(c, func(c context.Context) error {
		user, err := s.userRepository.FindByName(c, cmd.User)
		if err != nil {
			return errors.Wrap(domain.ErrNoSuchUser, err.Error())
		}

		if !s.encrypter.Verify(cmd.OldPassword, user.Password()) {
			return domain.ErrUserPassword
		}

		if err := user.ChangePassword(hashedPassword); err != nil {
			return err
		}

		if err := user.ChangeToken(s.factory.NewToken()); err != nil {
			return err
		}

		return s.userRepository.Save(c, user)
	})
}
