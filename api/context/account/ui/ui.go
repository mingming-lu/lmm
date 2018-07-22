package ui

import (
	"fmt"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
	"lmm/api/http"
)

type UI struct {
	app *appservice.AppService
}

func New(userRepo repository.UserRepository) *UI {
	app := appservice.New(userRepo)
	return &UI{app: app}
}

func (ui *UI) SignIn(c *http.Context) {
	user, err := ui.app.SignIn(c.Request.Body)
	switch err {
	case nil:
		c.JSON(http.StatusOK, SignInResponse{
			ID:    user.ID(),
			Name:  user.Name(),
			Token: user.Token(),
		})
	case service.ErrInvalidBody:
		http.BadRequest(c)
	case service.ErrInvalidUserNameOrPassword:
		c.String(http.StatusNotFound, service.ErrInvalidUserNameOrPassword.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) SignUp(c *http.Context) {
	id, err := ui.app.SignUp(c.Request.Body)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/users/%d", id)).String(http.StatusCreated, "Success")
	case service.ErrDuplicateUserName:
		c.String(http.StatusBadRequest, service.ErrDuplicateUserName.Error())
	case service.ErrInvalidBody:
		http.BadRequest(c)
	case service.ErrInvalidUserNameOrPassword:
		c.String(http.StatusBadRequest, service.ErrInvalidUserNameOrPassword.Error())
	default:
		c.Logger().Error(err.Error())
		http.InternalServerError(c)
	}
}

func (ui *UI) Verify(c *http.Context) {
	user, ok := c.Values().Get("user").(*model.User)
	if !ok {
		http.Unauthorized(c)
		return
	}
	c.JSON(http.StatusOK, VerifyResponse{
		ID:   user.ID(),
		Name: user.Name(),
	})
}

func (ui *UI) BearerAuth(handler http.Handler) http.Handler {
	return func(c *http.Context) {
		auth := c.Request.Header.Get("Authorization")
		user, err := ui.app.BearerAuth(auth)
		if err != nil {
			c.Logger().Error("%s: '%s'", err.Error(), auth)
			http.Unauthorized(c)
			return
		}
		c.Values().Set("user", user)
		handler(c)
	}
}
