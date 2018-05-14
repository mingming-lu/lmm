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
	if err != nil {
		c.Status(http.StatusNotFound).String(err.Error())
		return
	}

	response := SignInResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: user.Token,
	}

	c.Status(http.StatusOK).JSON(response)
}
