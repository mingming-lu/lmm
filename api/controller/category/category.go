package category

import (
	"encoding/json"
	"fmt"
	"lmm/api/db"
	model "lmm/api/domain/model/category"
	usecase "lmm/api/usecase/category"
	"lmm/api/usecase/user"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

func Register(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	m := model.Category{}
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	// TODO handle already exists
	id, err := usecase.Register(usr.ID, m.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}

	c.Header("Location", fmt.Sprintf("categories/%d", id)).Status(http.StatusCreated).String("success")
}

func Update(c *elesion.Context) {
	categoryID, err := strconv.ParseUint(c.Params.ByName("category"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid category id").Error(err.Error())
		return
	}

	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	m := model.Category{}
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	err = usecase.Update(usr.ID, categoryID, m.Name)
	if err != nil && err == db.ErrNoChange {
		c.Status(http.StatusNoContent).String("No change").Error(err.Error())
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}

	c.Status(http.StatusOK).String("success")
}

func GetByUser(c *elesion.Context) {
	userID, err := strconv.ParseUint(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user id").Error(err.Error())
	}

	blog, err := usecase.FetchByUser(userID)
	if err != nil {
		c.Status(http.StatusNotFound).String("Not Found").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(blog)
}

func Delete(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	categoryID, err := strconv.ParseUint(c.Params.ByName("category"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid category id").Error(err.Error())
		return
	}

	err = usecase.Delete(usr.ID, categoryID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
