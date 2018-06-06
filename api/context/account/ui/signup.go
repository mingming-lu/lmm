package ui

import (
	"fmt"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/repository"
	"lmm/api/http"
)

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
		c.Header("Location", fmt.Sprintf("/users/%d", id)).Status(http.StatusCreated).String("Success")
	case appservice.ErrDuplicateUserName:
		c.Status(http.StatusBadRequest).String(appservice.ErrDuplicateUserName.Error())
	default:
		http.InternalServerError(c)
	}
}
