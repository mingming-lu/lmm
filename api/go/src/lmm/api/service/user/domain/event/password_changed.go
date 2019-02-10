package event

import (
	"lmm/api/clock"

	"time"
)

const (
	TopicUserPasswordChanged = "user.password.changed"
)

type UserPasswordChanged struct {
	username  string
	password  string
	changedAt time.Time
}

func NewUserPasswordChangedEvent(username, password string) *UserPasswordChanged {
	return &UserPasswordChanged{
		username:  username,
		password:  password,
		changedAt: clock.Now(),
	}
}

func (event *UserPasswordChanged) Topic() string {
	return TopicUserPasswordChanged
}

func (event *UserPasswordChanged) UserName() string {
	return event.username
}

func (event *UserPasswordChanged) Password() string {
	return event.password
}

func (event *UserPasswordChanged) OccurredAt() time.Time {
	return event.changedAt
}
