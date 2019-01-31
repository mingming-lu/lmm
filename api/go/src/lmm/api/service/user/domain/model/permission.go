package model

type Permission uint64

const (
	// NoPermission means no any permission
	NoPermission Permission = 0
)

const (
	// PermissionAssignToOrdinary means permission to assign user to ordinary role
	PermissionAssignToOrdinary Permission = 1 << iota

	// PermissionAssignToAdmin means permission to assign user to admin role
	PermissionAssignToAdmin
)

// PermissionAssignToRole returns permission demanded to assign to role
func PermissionAssignToRole(role Role) Permission {
	switch role {
	case Admin:
		return PermissionAssignToAdmin
	case Ordinary:
		return PermissionAssignToOrdinary
	default:
		return NoPermission
	}
}
