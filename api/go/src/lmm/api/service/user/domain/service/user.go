package service

import (
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"

	"github.com/pkg/errors"
)

// AssignUserRole allows operator assign targetUser to role
func AssignUserRole(operator *model.User, targetUser *model.User, role model.Role) error {
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

	return nil
}
