package service

import (
	"context"

	"lmm/api/event"
	"lmm/api/service/user/domain"
	userEvent "lmm/api/service/user/domain/event"
	"lmm/api/service/user/domain/model"

	"github.com/pkg/errors"
)

// AssignUserRole allows operator assign targetUser to role
func AssignUserRole(c context.Context, operator *model.UserDescriptor, targetUser *model.UserDescriptor, role model.Role) error {
	if operator.Is(targetUser) {
		return domain.ErrCannotAssignSelfRole
	}

	perm := model.PermissionAssignToRole(role)
	if perm == model.NoPermission {
		return errors.Wrap(domain.ErrNoSuchRole, role.Name())
	}

	if !operator.Role().HasPermission(perm) {
		return domain.ErrNoPermission
	}

	if err := targetUser.AssignRole(role); err != nil {
		return err
	}

	return event.SyncBus().Publish(c, userEvent.NewUserRoleChangedEvent(
		operator.Name(), targetUser.Name(), targetUser.Role().Name(),
	))
}
