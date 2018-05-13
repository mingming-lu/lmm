package user

import (
	"encoding/json"
	"fmt"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/usecase"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func SignUp(c *elesion.Context) {
	id, err := usecase.New(repository.New()).SignUp(c.Request.Body)
	switch err {
	case nil:
		c.Header("location", fmt.Sprintf("/users/%d", id)).Status(http.StatusCreated).String("Success")
	case usecase.ErrDuplicateUserName:
		c.Status(http.StatusBadRequest).String(usecase.ErrDuplicateUserName.Error())
	default:
		c.Status(http.StatusInternalServerError).String(http.StatusText(http.StatusInternalServerError))
	}
}

func SignIn(c *elesion.Context) {
	info := model.Minimal{}
	err := json.NewDecoder(c.Request.Body).Decode(&info)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body").Error(err.Error())
		return
	}

	user, err := usecase.SignIn(info.Name, info.Password)
	if err != nil {
		c.Status(http.StatusNotFound).String(err.Error())
		return
	}

	c.Status(http.StatusOK).JSON(user)
}

func Verify(c *elesion.Context) {
	_, err := usecase.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("Authorized")
}
