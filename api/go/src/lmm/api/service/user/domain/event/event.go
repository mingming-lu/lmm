package event

import (
	"context"

	"lmm/api/messaging"
)

func PublishUserRoleChanged(c context.Context, operator, user, targetRole string) error {
	return messaging.SyncBus().Publish(c, NewUserRoleChangedEvent(
		operator, user, targetRole,
	))
}

func PublishUserPasswordChanged(c context.Context, username, newPassword string) error {
	return messaging.SyncBus().Publish(c, NewUserPasswordChangedEvent(username, newPassword))
}
