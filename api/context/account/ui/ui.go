package ui

import (
	"fmt"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/model"
	"lmm/api/http"
	"lmm/api/storage"
	"log"
)

type UI struct {
	app *appservice.AppService
}

func New(db *storage.DB) *UI {
	app := appservice.New(db)
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
	case appservice.ErrEmptyUserNameOrPassword:
		c.String(http.StatusBadRequest, appservice.ErrEmptyUserNameOrPassword.Error())
	case appservice.ErrInvalidUserNameOrPassword:
		c.String(http.StatusNotFound, appservice.ErrInvalidUserNameOrPassword.Error())
	default:
		http.InternalServerError(c)
	}
}

func (ui *UI) SignUp(c *http.Context) {
	auth := &Auth{}
	err := c.Request.ScanBody(&auth)
	if err != nil {
		http.BadRequest(c)
		return
	}
	id, err := ui.app.SignUp(c.Request.Body)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/users/%d", id)).String(http.StatusCreated, "Success")
	case appservice.ErrDuplicateUserName:
		c.String(http.StatusBadRequest, appservice.ErrDuplicateUserName.Error())
	default:
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
		user, err := ui.app.BearerAuth(c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Println(err)
			http.Unauthorized(c)
			return
		}
		c.Values().Set("user", user)
		handler(c)
	}
}
