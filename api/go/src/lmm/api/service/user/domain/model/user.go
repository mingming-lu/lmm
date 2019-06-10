package model

import (
	"net/mail"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"lmm/api/service/user/domain"
	"lmm/api/util/uuidutil"
)

var (
	patternUserName = regexp.MustCompile(`^[a-zA-Z]{1}[0-9a-zA-Z_-]{2,17}$`)
)

// UserID type
type UserID int64

// UserDescriptor describes user's basic infomation
type UserDescriptor struct {
	id           UserID
	name         string
	email        string
	role         Role
	registeredAt time.Time
}

// NewUserDescriptor creates a new *UserDescriptor
func NewUserDescriptor(id UserID, name, email string, role Role, registeredAt time.Time) (*UserDescriptor, error) {
	user := UserDescriptor{
		id:           id,
		role:         role,
		registeredAt: registeredAt,
	}

	if err := user.setName(name); err != nil {
		return nil, err
	}

	if err := user.setEmail(email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *UserDescriptor) ID() UserID {
	return user.id
}

// Name gets user's name
func (user *UserDescriptor) Name() string {
	return user.name
}

// Email gets user's email address
func (user *UserDescriptor) Email() string {
	return user.email
}

// Is compares if two use are the same
func (user *UserDescriptor) Is(target *UserDescriptor) bool {
	return user.Name() == target.Name()
}

// Role gets user's Role model
func (user *UserDescriptor) Role() Role {
	return user.role
}

// RegisteredAt gets user's register date
func (user *UserDescriptor) RegisteredAt() time.Time {
	return user.registeredAt
}

func (user *UserDescriptor) setName(name string) error {
	if !patternUserName.MatchString(name) {
		return domain.ErrInvalidUserName
	}
	user.name = name
	return nil
}

func (user *UserDescriptor) setRole(role Role) error {
	switch role {
	case Admin, Guest, Ordinary:
		user.role = role
		return nil
	default:
		return domain.ErrNoSuchRole
	}
}

func (user *UserDescriptor) setEmail(address string) error {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return errors.Wrap(domain.ErrInvalidEmail, err.Error())
	}
	user.email = addr.Address
	return nil
}

// User domain model
type User struct {
	UserDescriptor
	password string
	token    string
}

// NewUser creates a new user domain model
func NewUser(id UserID, name, email, password, token string, role Role, registeredDate time.Time) (*User, error) {
	descriptor, err := NewUserDescriptor(id, name, email, role, registeredDate)
	if err != nil {
		return nil, err
	}

	user := User{UserDescriptor: *descriptor}

	if err := user.setPassword(password); err != nil {
		return nil, err
	}

	if err := user.setToken(token); err != nil {
		return nil, err
	}

	return &user, nil
}

// Password gets user's encrypted password
func (user *User) Password() string {
	return user.password
}

// ChangePassword changes password to given newPassword
func (user *User) ChangePassword(newPassword string) error {
	return user.setPassword(newPassword)
}

func (user *User) setPassword(password string) error {
	user.password = password
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

func (user *User) ChangeToken(newToken string) error {
	return user.setToken(newToken)
}

// ChangeRole changes user's role
func (user *User) ChangeRole(role Role) error {
	return user.setRole(role)
}

// ChangeEmail changes user's email
func (user *User) ChangeEmail(newEmailAddress string) error {
	return user.setEmail(newEmailAddress)
}

// Is compares if two users are the same
func (user *User) Is(other *User) bool {
	return user.UserDescriptor.Is(&other.UserDescriptor)
}
