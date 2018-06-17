package ui

import (
	"fmt"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/http"
	"log"
)

func SignIn(c *http.Context) {
	auth := Auth{}
	err := c.Request.ScanBody(&auth)
	if err != nil {
		http.BadRequest(c)
		log.Println(err)
		return
	}

	user, err := appservice.New(repository.New()).SignIn(auth.Name, auth.Password)
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

func SignUp(c *http.Context) {
	uc := appservice.New(repository.New())
	auth := &Auth{}
	err := c.Request.ScanBody(&auth)
	if err != nil {
		http.BadRequest(c)
		return
	}
	id, err := uc.SignUp(auth.Name, auth.Password)
	switch err {
	case nil:
		c.Header("Location", fmt.Sprintf("/users/%d", id)).String(http.StatusCreated, "Success")
	case appservice.ErrDuplicateUserName:
		c.String(http.StatusBadRequest, appservice.ErrDuplicateUserName.Error())
	default:
		http.InternalServerError(c)
	}
}

func Verify(c *http.Context) {
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
