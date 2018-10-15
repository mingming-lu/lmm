package ui

import (
	"context"

	"lmm/api/http"
	"lmm/api/service/auth/application"
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
	token, err := ui.appService.BasicAuth(c, c.Request().Header.Get("Authorization"))
	switch err {
	case nil:
		c.JSON(http.StatusOK, loginResponse{
			AccessToken: token.Hashed(),
		})
	default:
		http.Warn(c, err.Error())
		http.Unauthorized(c)
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
		next(c.With(
			context.WithValue(context.Background(), http.StrCtxKey("user"), user),
		))
	}
}
