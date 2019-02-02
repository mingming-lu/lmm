package model

import (
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/model"
	"lmm/api/service/user/domain"
	"lmm/api/util/uuidutil"
)

var (
	patternUserName = regexp.MustCompile(`^[a-zA-Z]{1}[0-9a-zA-Z_-]{2,17}$`)
)

type UserDescriptor struct {
	model.Entity
	name string
	role Role
}

func NewUserDescriptor(name string, role Role) (*UserDescriptor, error) {
	user := UserDescriptor{}

	if err := user.setName(name); err != nil {
		return nil, err
	}

	user.role = role

	return &user, nil
}

// User domain model
type User struct {
	UserDescriptor
	password string
	token    string
}

// NewUser creates a new user domain model
func NewUser(name string, password Password, token string, role Role) (*User, error) {
	user := User{}

	if err := user.setName(name); err != nil {
		return nil, err
	}

	if err := user.setPassword(password.String()); err != nil {
		return nil, err
	}

	if err := user.setToken(token); err != nil {
		return nil, err
	}

	if err := user.AssignRole(role); err != nil {
		return nil, err
	}

	return &user, nil
}

// Name gets user's name
func (user *UserDescriptor) Name() string {
	return user.name
}

func (user *UserDescriptor) Is(target *UserDescriptor) bool {
	return user.Name() == target.Name()
}

func (user *UserDescriptor) Role() Role {
	return user.role
}

func (user *UserDescriptor) AssignRole(role Role) error {
	switch role {
	case Admin, Guest, Ordinary:
		user.role = role
		return nil
	default:
		return domain.ErrNoSuchRole
	}
}

func (user *UserDescriptor) setName(name string) error {
	if !patternUserName.MatchString(name) {
		return domain.ErrInvalidUserName
	}
	user.name = name
	return nil
}

// Password gets user's encrypted password
func (user *User) Password() string {
	return user.password
}

// ChangePassword changes password to given raw password
func (user *User) ChangePassword(newPassword Password) error {
	return user.setPassword(newPassword.String())
}

func (user *User) setPassword(password string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt password")
	}

	user.password = string(b)

	return nil
}

// Token gets user's token
func (user *User) Token() string {
	return user.token
}

func (user *User) setToken(token string) error {
	uuid, err := uuidutil.ParseString(token)
	if err != nil {
		return errors.Wrap(domain.ErrInvalidUserToken, err.Error())
	}
	if v := uuid.Version().String(); v != "VERSION_4" {
		return errors.Wrap(domain.ErrInvalidUserToken, "unexpected uuid version: "+v)
	}
	user.token = token
	return nil
}
