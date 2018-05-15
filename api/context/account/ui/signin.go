package ui

import (
	"encoding/json"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/usecase"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func SignIn(c *elesion.Context) {
	auth := Auth{}
	err := json.NewDecoder(c.Request.Body).Decode(&auth)
	if err != nil {
		c.Status(http.StatusBadRequest).String(http.StatusText(http.StatusBadRequest)).Error(err.Error())
		return
	}

	user, err := usecase.New(repository.New()).SignIn(auth.Name, auth.Password)
	switch err {
	case nil:
		c.Status(http.StatusOK).JSON(SignInResponse{
			ID:    user.ID,
			Name:  user.Name,
			Token: user.Token,
		})
	case usecase.ErrEmptyUserNameOrPassword:
		c.Status(http.StatusBadRequest).String(err.Error())
	case usecase.ErrInvalidUserNameOrPassword:
		c.Status(http.StatusNotFound).String(err.Error())
	default:
		c.Status(http.StatusInternalServerError).String(http.StatusText(http.StatusInternalServerError)).Error(err.Error())
	}
}
