package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	assert.Equal(t, Permission(0), NoPermission)
	assert.Equal(t, Permission(1), PermissionAssignToOrdinary)
	assert.Equal(t, Permission(2), PermissionAssignToAdmin)
	assert.Equal(t, Permission(3), PermissionAssignToAdmin|PermissionAssignToOrdinary)
}
