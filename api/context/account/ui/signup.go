package ui

import (
	"encoding/json"
	"fmt"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/repository"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func SignUp(c *elesion.Context) {
	uc := appservice.New(repository.New())
	auth := &appservice.Auth{}
	err := json.NewDecoder(c.Request.Body).Decode(auth)
	if err != nil {
		c.Status(http.StatusBadRequest).String(http.StatusText(http.StatusBadRequest))
		return
	}
	id, err := uc.SignUp(auth.Name, auth.Password)
	switch err {
	case nil:
		c.Header("location", fmt.Sprintf("/users/%d", id)).Status(http.StatusCreated).String("Success")
	case appservice.ErrDuplicateUserName:
		c.Status(http.StatusBadRequest).String(appservice.ErrDuplicateUserName.Error())
	default:
		c.Status(http.StatusInternalServerError).String(http.StatusText(http.StatusInternalServerError))
	}
}
