package model

import "lmm/api/testing"

func TestRoleAdmin(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Is("admin", Admin.Name())
	t.True(Admin.HasPermission(PermissionAssignToAdmin))
	t.True(Admin.HasPermission(PermissionAssignToOrdinary))
}
