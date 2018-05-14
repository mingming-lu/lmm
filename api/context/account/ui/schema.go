package ui

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
