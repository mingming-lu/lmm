package ui

type loginResponse struct {
	AccessToken string `json:"accessToken"`
}

type loginRequestBody struct {
	GrantType string `json:"grantType"`
}
