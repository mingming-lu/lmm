package ui

import (
	"lmm/api/context/account/appservice"
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
		c.Status(http.StatusOK).JSON(SignInResponse{
			ID:    user.ID(),
			Name:  user.Name(),
			Token: user.EncodedToken(),
		})
	case appservice.ErrEmptyUserNameOrPassword:
		c.Status(http.StatusBadRequest).String(err.Error())
	case appservice.ErrInvalidUserNameOrPassword:
		c.Status(http.StatusNotFound).String(err.Error())
	default:
		http.InternalServerError(c)
	}
}
