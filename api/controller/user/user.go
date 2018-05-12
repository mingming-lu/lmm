package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akinaru-lu/elesion"

	model "lmm/api/domain/model/user"
	repo "lmm/api/domain/repository/user"
	usecase "lmm/api/usecase/user"
)

func SignUp(c *elesion.Context) {
	id, err := usecase.New(repo.New()).SignUp(c.Request.Body)
	switch err {
	case nil:
		c.Header("location", fmt.Sprintf("/users/%d", id)).Status(http.StatusCreated).String("Success")
	case usecase.ErrDuplicateUserName:
		c.Status(http.StatusBadRequest).String(http.StatusText(http.StatusBadRequest))
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
