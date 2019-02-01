package service

import (
	"context"

	eventBus "lmm/api/domain/event"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/event"
	"lmm/api/service/user/domain/model"

	"github.com/pkg/errors"
)

// AssignUserRole allows operator assign targetUser to role
func AssignUserRole(c context.Context, operator *model.User, targetUser *model.User, role model.Role) error {
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

	return eventBus.SyncBus().Publish(c, event.NewUserRoleChangedEvent(
		operator.Name(), targetUser.Name(), targetUser.Role().Name(),
	))
}
