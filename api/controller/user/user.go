package user

import (
	"lmm/api/context/account/appservice"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func Verify(c *elesion.Context) {
	_, err := appservice.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("Authorized")
}
