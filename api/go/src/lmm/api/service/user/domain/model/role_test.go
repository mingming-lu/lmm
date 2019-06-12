package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleAdmin(t *testing.T) {
	assert.Equal(t, "admin", Admin.Name())
	assert.True(t, Admin.HasPermission(PermissionAssignToAdmin))
	assert.True(t, Admin.HasPermission(PermissionAssignToOrdinary))
}
