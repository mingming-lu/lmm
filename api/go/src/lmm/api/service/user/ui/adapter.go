package ui

type signUpRequestBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type assignRoleRequestBody struct {
	Role string `json:"role"`
}

type userView struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
