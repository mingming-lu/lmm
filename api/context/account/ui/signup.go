package ui

import (
	"encoding/json"
	"fmt"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/usecase"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func SignUp(c *elesion.Context) {
	uc := usecase.New(repository.New())
	auth := &usecase.Auth{}
	err := json.NewDecoder(c.Request.Body).Decode(auth)
	if err != nil {
		c.Status(http.StatusBadRequest).String(http.StatusText(http.StatusBadRequest))
	}
	id, err := uc.SignUp(auth.Name, auth.Password)
	switch err {
	case nil:
		c.Header("location", fmt.Sprintf("/users/%d", id)).Status(http.StatusCreated).String("Success")
	case usecase.ErrDuplicateUserName:
		c.Status(http.StatusBadRequest).String(usecase.ErrDuplicateUserName.Error())
	default:
		c.Status(http.StatusInternalServerError).String(http.StatusText(http.StatusInternalServerError))
	}
}
