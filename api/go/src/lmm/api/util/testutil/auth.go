package testutil

import (
	"encoding/json"
	"strings"

	"lmm/api/service/auth/domain/model"
	"lmm/api/service/auth/domain/service"
)

// EncodeToken calls lmm/api/service/auth/domain/service.TokenService.Encode
func EncodeToken(rawToken string) *model.Token {
	token, err := service.NewTokenService().Encode(rawToken)
	if err != nil {
		panic(err)
	}
	return token
}

// ExtractAccessToken tries to extract access token from given string
func ExtractAccessToken(s string) string {
	// avoid cycle import, see lmm/api/service/auth/ui/adapter.go
	type loginResponse struct {
		AccessToken string `json:"accessToken"`
	}

	schema := loginResponse{}

	if err := json.NewDecoder(strings.NewReader(s)).Decode(&schema); err != nil {
		panic(err)
	}

	return schema.AccessToken
}