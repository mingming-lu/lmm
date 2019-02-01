package service

import "lmm/api/service/user/domain/model"

func RoleAdapter(name string) model.Role {
	switch name {
	case "admin":
		return model.Admin
	case "ordinary":
		return model.Ordinary
	default:
		return model.Guest
	}
}
