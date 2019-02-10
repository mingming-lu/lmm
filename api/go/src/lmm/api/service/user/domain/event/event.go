package event

import (
	"lmm/api/event"

	"context"
)

func PublishUserRoleChanged(c context.Context, operator, user, targetRole string) error {
	return event.SyncBus().Publish(c, NewUserRoleChangedEvent(
		operator, user, targetRole,
	))
}

func PublishUserPasswordChanged(c context.Context, username string) error {
	return event.SyncBus().Publish(c, NewUserPasswordChangedEvent(username))
}
