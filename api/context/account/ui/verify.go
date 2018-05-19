package ui

import (
	"lmm/api/context/account/domain/model"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func Verify(c *elesion.Context) {
	user, ok := c.Fields().Get("user").(*model.User)
	if !ok {
		c.Status(http.StatusUnauthorized).String(http.StatusText(http.StatusUnauthorized))
		return
	}
	c.Status(http.StatusOK).JSON(VerifyResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}
