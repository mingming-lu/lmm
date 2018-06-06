package ui

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/http"
)

func Verify(c *http.Context) {
	user, ok := c.Values().Get("user").(*model.User)
	if !ok {
		http.Unauthorized(c)
		return
	}
	c.Status(http.StatusOK).JSON(VerifyResponse{
		ID:   user.ID,
		Name: user.Name,
	})
}
