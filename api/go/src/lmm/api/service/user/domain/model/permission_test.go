package model

import "lmm/api/testing"

func TestPermission(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Is(Permission(0), NoPermission)
	t.Is(Permission(1), PermissionAssignToOrdinary)
	t.Is(Permission(2), PermissionAssignToAdmin)
	t.Is(Permission(3), PermissionAssignToAdmin|PermissionAssignToOrdinary)
}
