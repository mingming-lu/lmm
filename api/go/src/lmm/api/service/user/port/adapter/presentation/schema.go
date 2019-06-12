package presentation

import "encoding/json"

type signUpRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type basicAuth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type assignRoleRequestBody struct {
	Role string `json:"role"`
}

type changePasswordRequestBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type accessTokenView struct {
	AccessToken string `json:"access_token"`
}

type userView struct {
	Name           string `json:"name"`
	Role           string `json:"role"`
	RegisteredDate int64  `json:"registered_date,string"`
}

type usersView struct {
	Users  []userView  `json:"users"`
	Page   json.Number `json:"page"`
	Count  json.Number `json:"count"`
	Total  uint        `json:"total"`
	SortBy string      `json:"sort_by"`
	Sort   string      `json:"sort"`
}
