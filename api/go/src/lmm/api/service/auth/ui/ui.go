package ui

import (
	"lmm/api/http"
	"lmm/api/service/auth/application"
	"lmm/api/service/auth/application/command"
)

const (
	errGrantTypeRequired = "grant type required"
)

// UI is the auth UI
type UI struct {
	appService *application.Service
}

// NewUI returns a new ui pointer
func NewUI(service *application.Service) *UI {
	return &UI{
		appService: service,
	}
}

// Login handles POST /v1/auth/login
func (ui *UI) Login(c http.Context) {
	auth := c.Request().Header.Get("Authorization")
	var body loginRequestBody
	if err := c.Request().Bind(&body); err != nil {
		http.Warn(c, err.Error())
		c.String(http.StatusBadRequest, errGrantTypeRequired)
		return
	}

	token, err := ui.appService.Login(c, command.LoginCommand{
		AccessToken: auth,
		BasicAuth:   auth,
		GrantType:   body.GrantType,
	})
	switch err {
	case nil:
		c.JSON(http.StatusOK, loginResponse{
			AccessToken: token.Hashed(),
		})
	default:
		http.Panic(c, err.Error())
	}
}

// BearerAuth provides a bearer auth middleware
func (ui *UI) BearerAuth(next http.Handler) http.Handler {
	return func(c http.Context) {
		user, err := ui.appService.BearerAuth(c, c.Request().Header.Get("Authorization"))
		if err != nil {
			http.Warn(c, err.Error())
			http.Unauthorized(c)
			return
		}
		c.Request().Header.Set("X-LMM-ID", user.Name())
		next(c)
	}
}
