package model

// Role is user' role
type Role struct {
	name      string
	permssion Permission
}

// Name returns r's name
func (r Role) Name() string {
	return r.name
}

// HasPermission returns ture if r has perm
func (r Role) HasPermission(permission Permission) bool {
	return r.permssion&permission == permission
}

var (
	// Admin role
	Admin = Role{
		name:      "admin",
		permssion: PermissionAssignToAdmin | PermissionAssignToOrdinary,
	}

	// Guest role
	Guest = Role{
		name:      "guest",
		permssion: NoPermission,
	}

	// Ordinary role
	Ordinary = Role{
		name:      "ordinary",
		permssion: NoPermission,
	}
)
