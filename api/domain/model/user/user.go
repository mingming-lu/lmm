package user

type Minimal struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type User struct {
	ID        int64
	Name      string
	Password  string
	GUID      string
	Token     string
	CreatedAt string
}
