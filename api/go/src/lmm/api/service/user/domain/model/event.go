package model

import "context"

type UserEventPublisher interface {
	NotifyUserRegistered(context.Context, UserID) error
	NotifyUserPasswordChanged(context.Context, UserID) error
}
