package command

// LoginCommand contains parameters needed for login
type LoginCommand struct {
	AccessToken string
	BasicAuth   string
	GrantType   string
}
