package application

import (
	"context"

	"lmm/api/pkg/transaction"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/application/query"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"

	"github.com/pkg/errors"
)

// Service is a application service
type Service struct {
	encrypter          service.EncryptService
	factory            *factory.Factory
	tokenService       model.TokenService
	transactionManager transaction.Manager
	userRepository     repository.UserRepository
}

// NewService creates a new Service pointer
func NewService(
	encrypter service.EncryptService,
	tokenService model.TokenService,
	txManager transaction.Manager,
	userRepository repository.UserRepository,
) *Service {
	return &Service{
		encrypter:          encrypter,
		factory:            factory.NewFactory(encrypter, userRepository),
		tokenService:       tokenService,
		transactionManager: txManager,
		userRepository:     userRepository,
	}
}

// RegisterNewUser registers new user
func (s *Service) RegisterNewUser(c context.Context, cmd command.Register) (int64, error) {
	var userID int64

	err := s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		if user, err := s.userRepository.FindByName(tx, cmd.UserName); err != domain.ErrNoSuchUser {
			if user != nil {
				return domain.ErrUserNameAlreadyUsed
			}
			return errors.Wrap(err, "error occurred when checking user name duplication")
		}

		user, err := s.factory.NewUser(tx, cmd.UserName, cmd.EmailAddress, cmd.Password)
		if err != nil {
			return errors.Wrap(err, "invalid user")
		}

		if err := s.userRepository.Save(tx, user); err != nil {
			return errors.Wrap(err, "failed to save user")
		}

		userID = int64(user.ID())

		return nil
	}, nil)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// BasicAuth authenticate user by basic auth
func (s *Service) BasicAuth(c context.Context, cmd command.Login) (user *model.User, err error) {
	err = s.transactionManager.RunInTransaction(c,
		func(tx transaction.Transaction) error {
			user, err = s.login(tx, cmd.UserName, cmd.Password)
			if err != nil {
				return errors.Wrap(err, "failed to login")
			}
			return nil
		},
		&transaction.Option{ReadOnly: true},
	)
	return
}

// RefreshAccessToken refreshes a valid oldAccessToken into a valid newAccessToken
func (s *Service) RefreshAccessToken(c context.Context, hashed string) (newAccessToken *model.AccessToken, err error) {
	token, err := s.tokenService.Decrypt(hashed)
	if err != nil {
		return nil, errors.Wrap(err, "invalid access token")
	}

	if token.Expired() {
		return nil, errors.New("access token expired")
	}

	err = s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		user, err := s.userRepository.FindByToken(tx, token.Raw())
		if err != nil {
			return err
		}
		newAccessToken, err = s.tokenService.Encrypt(user.Token())
		if err != nil {
			panic(errors.Wrap(err, "internal error"))
		}
		return nil
	}, &transaction.Option{ReadOnly: true})

	return
}

func (s *Service) login(tx transaction.Transaction, username, password string) (*model.User, error) {
	user, err := s.userRepository.FindByName(tx, username)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchUser, err.Error())
	}

	if !s.encrypter.Verify(password, user.Password()) {
		return nil, domain.ErrUserPassword
	}

	return user, nil
}

// AssignRole handles command which operator assign user to role
func (s *Service) AssignRole(c context.Context, cmd command.AssignRole) error {
	panic("not implemented")
}

const maxCount uint = 100

func (s *Service) ViewAllUsersByOptions(c context.Context, query query.ViewAllUsers) ([]*model.UserDescriptor, uint, error) {
	panic("not implemented")
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
		return errors.Wrap(err, "invalid password")
	}

	return s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		user, err := s.login(tx, cmd.User, cmd.OldPassword)
		if err != nil {
			return errors.Wrap(err, "failed to login")
		}

		if err := user.ChangePassword(hashedPassword); err != nil {
			return errors.Wrap(err, "failed to change password")
		}

		if err := user.ChangeToken(s.factory.NewToken()); err != nil {
			return errors.Wrap(err, "failed to change token")
		}

		if err := s.userRepository.Save(tx, user); err != nil {
			return errors.Wrap(err, "failed to save user after password and token changed")
		}

		return nil
	}, nil)
}
